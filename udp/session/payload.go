package session

type payload struct {
	buffer []byte
	length int
}

func (p payload) buf() []byte {
	if p.buffer == nil {
		return nil
	}

	return p.buffer[0:p.length]
}

func (p payload) release() {
	if p.buffer == nil {
		return
	}
	bufferPool.Put(p.buffer)
}
