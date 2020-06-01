package redigo

import (
	"sync"

	"github.com/boxofimagination/bxdk/go/redis/engine"
)

const (
	// default value of pipeline's numCmdHint
	defaultPipelineNumCmdHint = 2
)

// pipelineCmdErr is pipeline command and error
type pipelineCmdErr struct {
	cmd  string
	args []interface{}
	err  error
}

func (pce pipelineCmdErr) Name() string {
	return pce.cmd
}

func (pce pipelineCmdErr) Args() []interface{} {
	return pce.args
}

func (pce pipelineCmdErr) Err() error {
	return pce.err
}

// Pipeline creates new redigo pipeline
func (r *Redigo) Pipeline(retry, cmdNumHint int) engine.Pipeliner {
	if retry <= 0 {
		retry = 1
	}
	if cmdNumHint == 0 {
		cmdNumHint = defaultPipelineNumCmdHint
	}

	p := &pipeline{
		cli:        r,
		retry:      retry,
		cmdNumHint: cmdNumHint,
	}
	p.resetCmdBuf()
	return p
}

// pipeline is helper for redigo pipelined command.
// This pipeline could be used multiple times and from concurrent goroutines
type pipeline struct {
	cli *Redigo

	mux sync.Mutex

	// cmdErrs is buffer of command sent to this pipeline
	cmdErrs []engine.CmdErr

	// number of retry we do in case of pipeline execution failed.
	retry int

	// hints about number of commands on each pipeline
	cmdNumHint int
}

// AddRawCmd adds raw redis command to the pipeline
func (p *pipeline) AddRawCmd(cmd string, args ...interface{}) {
	p.mux.Lock()
	p.cmdErrs = append(p.cmdErrs, pipelineCmdErr{
		cmd:  cmd,
		args: args,
	})
	p.mux.Unlock()
}

// Exec executes the pipeline
func (p *pipeline) Exec() ([]engine.CmdErr, int, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	defer p.resetCmdBuf()

	var (
		ret      []engine.CmdErr
		err      error
		firstErr int
	)

	// execute the pipeline with retry, if needed
	for i := 0; i < p.retry; i++ {
		ret, firstErr, err = p.exec()
		if err == nil {
			return ret, firstErr, nil
		}
	}
	return ret, firstErr, err
}

func (p *pipeline) exec() ([]engine.CmdErr, int, error) {
	firstErr := -1

	conn, err := p.cli.getConn()
	if err != nil {
		return nil, firstErr, err
	}
	defer conn.Close()

	// buffer the command to redigo connection.
	// we don't do it earlier because if we do it earlier, we get more risk
	// that the connection become invalid/closed when we finally execute the pipeline.
	for _, cmd := range p.cmdErrs {
		err = conn.Send(cmd.Name(), cmd.Args()...)
		if err != nil {
			return nil, firstErr, err
		}
	}

	// flush the command
	err = conn.Flush()
	if err != nil {
		return nil, firstErr, err
	}

	// receive the errors
	for i := 0; i < len(p.cmdErrs); i++ {
		_, err = conn.Receive()
		if err != nil {
			pce := p.cmdErrs[i].(pipelineCmdErr)
			pce.err = err
			p.cmdErrs[i] = pce
			if firstErr < 0 {
				firstErr = i
			}
		}
	}
	return p.cmdErrs, firstErr, nil
}

// Discard resets the pipeline and discards queued commands
func (p *pipeline) Discard() error {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.resetCmdBuf()
	return nil
}

// Close closes the pipeline, releasing any open resources.
func (p *pipeline) Close() error {
	return nil
}

func (p *pipeline) resetCmdBuf() {
	p.cmdErrs = make([]engine.CmdErr, 0, p.cmdNumHint)
}
