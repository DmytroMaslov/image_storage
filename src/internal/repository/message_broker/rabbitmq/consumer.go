package rabbitmq

import (
	"encoding/json"
	"image_storage/src/config"
	"image_storage/src/internal"
	"image_storage/src/internal/domain"
	"image_storage/src/internal/services"
	"log"

	"github.com/streadway/amqp"
)

type ImageConsumer struct {
	config  *config.Config
	amqpCon *amqp.Connection
	service services.ImageUseCases
}

func NewImageConsumer(config *config.Config, service services.ImageUseCases) (internal.ImageConsumer, error) {
	amqpCon, err := amqp.Dial(config.QueueConnectionString)
	if err != nil {
		return nil, err
	}
	var imageConsumer internal.ImageConsumer = &ImageConsumer{
		config:  config,
		amqpCon: amqpCon,
		service: service,
	}
	return imageConsumer, nil
}

func (c *ImageConsumer) SetUpQueue(queueName string) (ch *amqp.Channel, err error) {
	ch, err = c.amqpCon.Channel()
	if err != nil {
		return
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
		return
	}
	return
}

func (c *ImageConsumer) RunConsumer(workerPoolSize int) (err error) {
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
		return
	}
	for i := 0; i < workerPoolSize; i++ {
		go c.worker(msgs)
	}
	return
}

func (c *ImageConsumer) worker(msgs <-chan amqp.Delivery) {
	for delivery := range msgs {
		myImage := new(domain.MyImage)
		err := json.Unmarshal(delivery.Body, &myImage)
		if err != nil {
			log.Println(err)
		}
		//add ack or reject logic
		//var newImage domain.MyImage
		newImage, err := c.service.OptimizeImages(myImage)
		if err != nil {
			log.Println(err)
		}
		err = c.service.Save(newImage)
		if err != nil {
			log.Println(err)
		}

	}
}
