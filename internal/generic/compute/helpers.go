package compute

type chanByteWriterCloser struct {
	ch chan []byte
}

func (c *chanByteWriterCloser) Write(p []byte) (n int, err error) {
	c.ch <- p
	return len(p), nil
}

func (c *chanByteWriterCloser) Close() error {
	close(c.ch)
	return nil
}
