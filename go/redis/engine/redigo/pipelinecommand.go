package redigo

func (p *pipeline) Incr(key string) {
	p.AddRawCmd("INCR", key)
}

func (p *pipeline) IncrBy(key string, value int64) {
	p.AddRawCmd("INCRBY", key, value)
}

func (p *pipeline) Decr(key string) {
	p.AddRawCmd("DECR", key)
}

func (p *pipeline) DecrBy(key string, value int64) {
	p.AddRawCmd("DECRBY", key, value)
}

func (p *pipeline) Expire(key string, expiry int) {
	p.AddRawCmd("EXPIRE", key, expiry)
}

func (p *pipeline) Delete(keys ...string) {
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		args[i] = key
	}
	p.AddRawCmd("DEL", args...)
}

func (p *pipeline) HMSet(key string, kv map[string]interface{}) {
	var (
		args = make([]interface{}, 1+(len(kv)*2))
		idx  = 1
	)
	args[0] = key
	for k, v := range kv {
		args[idx] = k
		args[idx+1] = v
		idx += 2
	}
	p.AddRawCmd("HMSET", args...)
}

func (p *pipeline) HDel(key string, fields ...string) {
	args := make([]interface{}, len(fields)+1)
	args[0] = key
	for i, field := range fields {
		args[i+1] = field
	}
	p.AddRawCmd("HDEL", args...)
}
