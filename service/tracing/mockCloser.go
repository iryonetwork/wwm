package tracing

type MockCloser struct {
}

func (c MockCloser) Close() error {
	return nil
}
