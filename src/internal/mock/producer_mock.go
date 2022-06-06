package mock

type ProducerMock struct {
	FunctionPublish func([]byte, string) error
}

func (p *ProducerMock) Publish(body []byte, contentType string) error {
	return p.FunctionPublish(body, contentType)
}
