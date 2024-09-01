package compute

type chanByteWriterCloser struct {
	ch chan []byte
}

func (c *chanByteWriterCloser) Write(p []byte) (n int, err error) {
	// p may be reused by the caller, so we need to copy it
	pCopy := make([]byte, len(p))
	copy(pCopy, p)
	c.ch <- pCopy
	return len(p), nil
}

func (c *chanByteWriterCloser) Close() error {
	close(c.ch)
	return nil
}
