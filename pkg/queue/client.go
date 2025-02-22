package queue

import (
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

// A Struct that represents the RabbitMQ Client
type RabbitMQClient struct {
	Name        string
	Conn        *amqp.Connection
	Channel     *amqp.Channel
	ConnectOnce sync.Once
	isConnected bool
}

// NewRabbitMQClient creates a new instance of RabbitMQClient
func NewRabbitMQClient(clientName string) *RabbitMQClient {
	return &RabbitMQClient{
		Name:        clientName,
		isConnected: false,
	}
}

// Connect to RabbitMQ
func (c *RabbitMQClient) Connect(amqpUrl string) error {
	var err error
	c.ConnectOnce.Do(func() {
		c.Conn, err = amqp.Dial(amqpUrl)
		if err != nil {
			return
		}
		c.Channel, err = c.Conn.Channel()
		if err != nil {
			return
		}
		c.isConnected = true
	})
	return nil
}

// Close the RabbitMQ connection
func (c *RabbitMQClient) Close() error {
	var err error
	// lets close the channel first
	if c.Channel != nil {
		err = c.Channel.Close()
	}
	if c.Conn != nil {
		err = c.Conn.Close()
	}
	// reset singleton
	c.Conn = nil
	c.Channel = nil
	c.isConnected = false
	c.ConnectOnce = sync.Once{}
	return err
}

// IsConnected returns true if the client is connected
func (c *RabbitMQClient) IsConnected() bool {
	return c.isConnected
}

// NewQueue Creates a New Remote RabbitMQ Queue
func (c *RabbitMQClient) NewQueue(queueName string) (*amqp.Queue, error) {
	q, err := c.Channel.QueueDeclare(
		queueName, // name
		true,      // durable - we want the queue to persist even if the server restarts
		false,     // delete when unused - we want the queue to always be available
		false,     // exclusive - only available on the current connection
		false,     // no-wait - we want the call to block if the queue is not available
		nil,       // arguments - we want the queue to be durable
	)
	return &q, err
}

// DeleteQueue Deletes a Remote RabbitMQ Queue
// Returns the Number of Messages deleted and error if the Queue does not exist
func (c *RabbitMQClient) DeleteQueue(queueName string) (int, error) {
	return c.Channel.QueueDelete(queueName, false, false, false)
}

// PurgeQueue Deletes a Remote RabbitMQ Queue
// Returns the Number of Messages deleted and error if the Queue does not exist
func (c *RabbitMQClient) PurgeQueue(queueName string) (int, error) {
	return c.Channel.QueuePurge(queueName, false)
}

// SetQoS sets the Quality of Service for a channel
func (c *RabbitMQClient) SetQoS() error {
	// QoS 0 - At Most Once
	// QoS 1 - At Least Once
	// QoS 2 - Exactly Once
	// prefetch count is 1
	// prefetch size is 0
	// global is false
	return c.Channel.Qos(1, 0, false)
}

// CheckQueueExists checks if a queue exists
func (c *RabbitMQClient) CheckQueueExists(queueName string) (bool, error) {
	_, err := c.Channel.QueueDeclarePassive(
		queueName, // name
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConsumeQueue Loops through a queue and calls the handleMessage function for each message
func (c *RabbitMQClient) ConsumeQueue(queueName string, handleMessage func([]byte) error) error {
	// declare a queue
	_, err := c.Channel.QueueDeclare(
		queueName, // name
		true,      // durable - we want the queue to persist even if the server restarts
		false,     // delete when unused - we want the queue to always be available
		false,     // exclusive - only available on the current connection
		false,     // no-wait - we want the call to block if the queue is not available
		nil,       // arguments - we want the queue to be durable
	)
	if err != nil {
		return err
	}
	// start consuming
	msgs, err := c.Channel.Consume(
		queueName, // queue
		c.Name,    // consumer
		false,     // auto-ack - we want the client to ack the message for us
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}
	// // loop through the messages
	// for d := range msgs {
	// 	callback(d)
	// }
	// return nil

	// blocking channel
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("[%s] Received message: %s", queueName, d.Body)
			err := handleMessage(d.Body)
			if err != nil {
				log.Printf("[%s] Error processing message: %v", queueName, err)
				if err := d.Nack(false, true); err != nil {
					log.Printf("[%s] Failed to nack message: %v", queueName, err)
				}
			} else {
				if err := d.Ack(false); err != nil {
					log.Printf("[%s] Failed to ack message: %v", queueName, err)
				}
			}
		}
	}()

	log.Printf(" [*] Waiting for messages on queue '%s'. To exit press CTRL+C", queueName)
	<-forever

	return nil
}
