package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type RabbitMQConfig struct {
	channel *amqp.Channel
}

func Connect(rabbitmqURL string) (*RabbitMQConfig, error) {
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, err
	}

	ch, err := DeclareChannel(conn)
	if err != nil {
		return nil, err
	}

	return &RabbitMQConfig{channel: ch}, nil
}

func DeclareChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	ch.Qos(1, 0, false)

	return ch, nil
}

func (c *RabbitMQConfig) DeclareQueue(queueName string) (amqp.Queue, error) {
	q, err := c.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return amqp.Queue{}, err
	}
	return q, nil
}

func (c *RabbitMQConfig) Publish(q amqp.Queue, body []byte) error {
	return c.channel.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}

func (c *RabbitMQConfig) Consume(q amqp.Queue) (<-chan amqp.Delivery, error) {
	msgs, err := c.channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
