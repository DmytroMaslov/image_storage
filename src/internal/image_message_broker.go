package internal

type ImageProducer interface {
	Publish(body []byte, contentType string) error
}

type ImageConsumer interface {
	RunConsumer(workerPoolSize int) (err error)
}
