package grace

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/boxofimagination/bxdk/go/log"
	"github.com/boxofimagination/bxdk/go/socketmaster"
)

func ServerHTTP(srv *http.Server, address string) error {
	lis, err := Listen(address)
	if err != nil {
		return err
	}

	stoppedCh := WaitTermSig(srv.Shutdown)

	log.Printf("http server running on address: %v", address)

	go socketmaster.NotifyMaster()

	// start serving
	if err := srv.Serve(lis); err != http.ErrServerClosed {
		return err
	}

	<-stoppedCh
	log.Println("HTTP server stopped")
	return nil
}

// It returns channel which will be closed after the signal received and the handler executed.
// We can use the signal to wait for the shutdown to be finished.
func WaitTermSig(handler func(ctx context.Context) error) <-chan struct{} {
	stoppedCh := make(chan struct{})
	go func() {
		signals := make(chan os.Signal, 1)

		// wait for the sigterm
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-signals

		// We received an os signal, shut down.
		if err := handler(context.Background()); err != nil {
			log.Printf("graceful shutdown  failed: %v", err)
		} else {
			log.Println("gracefull shutdown succeed")
		}
		close(stoppedCh)
	}()
	return stoppedCh
}

func Listen(port string) (net.Listener, error) {
	var l net.Listener

	// see if we run under socketmaster
	fd := os.Getenv("EINHORN_FDS")
	if fd != "" {
		sock, err := strconv.Atoi(fd)
		if err != nil {
			return nil, err
		}
		log.Println("detected socketmaster, listening on", fd)
		file := os.NewFile(uintptr(sock), "listener")
		fl, err := net.FileListener(file)
		if err != nil {
			return nil,err
		}
		l=fl
	}

	if l != nil { // we already have the listener, which listen on EINHORN_FDS
		notifSocketMaster()
		return l, nil
	}

	// we are not using socketmaster, no need to notify
	return net.Listen("tcp4", port)
}



func notifSocketMaster() {
	go func() {
		err := socketmaster.NotifyMaster()
		if err != nil {
			log.Printf("failed to notify socketmaster: %v, maybe we don't use it?", err)
		} else {
			log.Println("successfully notify socketmaster")
		}
	}()
}
