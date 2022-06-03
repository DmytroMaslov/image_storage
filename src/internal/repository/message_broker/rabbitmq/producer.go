package rabbitmq

import (
	"image_storage/src/config"
	"image_storage/src/internal"

	"github.com/streadway/amqp"
)

type ImageProducer struct {
	amqpChan *amqp.Channel
	config   *config.Config
}

func NewImageProducer(config *config.Config) (internal.ImageProducer, error) {
	amqpCon, err := amqp.Dial(config.QueueConnectionString)
	if err != nil {
		return nil, err
	}
	amqpChan, err := amqpCon.Channel()
	if err != nil {
		return nil, err
	}
	var ImageProducer internal.ImageProducer = &ImageProducer{
		amqpChan: amqpChan,
		config:   config,
	}
	return ImageProducer, nil
}

func (p *ImageProducer) SetUpQueue(queueName string) (q amqp.Queue, err error) {
	q, err = p.amqpChan.QueueDeclare(
		p.config.QueueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return
	}
	return q, nil
}

func (p *ImageProducer) CloseChannel() error {
	return p.amqpChan.Close()
}

func (p *ImageProducer) Publish(body []byte, contentType string) error {
	err := p.amqpChan.Publish(
		"",
		p.config.QueueName,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		})
	return err
}
