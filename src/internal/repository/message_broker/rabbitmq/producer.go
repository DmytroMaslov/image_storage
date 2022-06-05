package rabbitmq

import (
	"image_storage/src/config"
	"image_storage/src/internal"
	"image_storage/src/internal/image_errors"
	"image_storage/src/pkg"

	"github.com/streadway/amqp"
)

type ImageProducer struct {
	amqpChan *amqp.Channel
	config   *config.Config
	log      pkg.Logger
}

func NewImageProducer(config *config.Config, log pkg.Logger) (internal.ImageProducer, error) {
	amqpCon, err := amqp.Dial(config.QueueConnectionString)
	if err != nil {
		log.Errorf("producer can't connect to Rabbit, reason: %s", err.Error())
		return nil, image_errors.ErrCantConnect
	}
	amqpChan, err := amqpCon.Channel()
	if err != nil {
		log.Errorf("producer can't create channel, reason: %s", err.Error())
		return nil, image_errors.ErrChannelCreate
	}
	var ImageProducer internal.ImageProducer = &ImageProducer{
		amqpChan: amqpChan,
		config:   config,
		log:      log,
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
		p.log.Errorf("producer can't declare queue, reason: %s", err.Error())
		err = image_errors.ErrQueueCreate
		return
	}
	return q, nil
}

func (p *ImageProducer) CloseChannel() error {
	return p.amqpChan.Close()
}

func (p *ImageProducer) Publish(body []byte, contentType string) (err error) {
	err = p.amqpChan.Publish(
		"",
		p.config.QueueName,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		})
	if err != nil {
		p.log.Errorf("producer can't publish message, reason: %s", err.Error())
		err = image_errors.ErrPublish
		return
	}
	return
}
