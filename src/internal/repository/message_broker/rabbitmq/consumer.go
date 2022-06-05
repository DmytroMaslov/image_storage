package rabbitmq

import (
	"encoding/json"
	"image_storage/src/config"
	"image_storage/src/internal"
	"image_storage/src/internal/domain"
	"image_storage/src/internal/image_errors"
	"image_storage/src/internal/services"
	"image_storage/src/pkg"

	"github.com/streadway/amqp"
)

type ImageConsumer struct {
	config  *config.Config
	amqpCon *amqp.Connection
	service services.ImageUseCases
	log     pkg.Logger
}

func NewImageConsumer(config *config.Config, service services.ImageUseCases, log pkg.Logger) (internal.ImageConsumer, error) {
	amqpCon, err := amqp.Dial(config.QueueConnectionString)
	if err != nil {
		log.Errorf("consumer can't connect to Rabbit, reason: %s", err.Error())
		return nil, image_errors.ErrCantConnect
	}
	var imageConsumer internal.ImageConsumer = &ImageConsumer{
		config:  config,
		amqpCon: amqpCon,
		service: service,
		log:     log,
	}
	return imageConsumer, nil
}

func (c *ImageConsumer) SetUpQueue(queueName string) (ch *amqp.Channel, err error) {
	ch, err = c.amqpCon.Channel()
	if err != nil {
		c.log.Errorf("consumer can't create channel, reason: %s", err.Error())
		return nil, image_errors.ErrChannelCreate
	}
	_, err = ch.QueueDeclare(
		c.config.QueueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		c.log.Errorf("consumer can't create channel, reason: %s", err.Error())
		return nil, image_errors.ErrChannelCreate
	}
	return
}

func (c *ImageConsumer) RunConsumer(workerPoolSize int, qualityArray []int) (err error) {
	ch, err := c.SetUpQueue(c.config.QueueName)
	if err != nil {
		return
	}
	msgs, err := ch.Consume(
		c.config.QueueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		c.log.Errorf("consumer can't consume message from channel, reason: %s", err.Error())
		return image_errors.ErrConsume
	}
	for i := 0; i < workerPoolSize; i++ {
		go c.worker(msgs, qualityArray)
	}
	return
}

func (c *ImageConsumer) worker(msgs <-chan amqp.Delivery, qualityArray []int) {
	for delivery := range msgs {
		myImage := new(domain.MyImage)
		err := json.Unmarshal(delivery.Body, &myImage)
		if err != nil {
			c.log.Errorf("worker can't unmarshal message, reason: %s", err.Error())
		}
		//add ack or reject logic
		c.log.Infof("message successful consumed with id: %s", myImage.Id)
		err = c.service.ReduceQuality(myImage, qualityArray)
		if err != nil {
			c.log.Errorf("worker can't reduce quality: %s", err.Error())
		}

	}
}
