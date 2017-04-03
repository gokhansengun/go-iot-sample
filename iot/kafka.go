package iot

import (
	"github.com/Shopify/sarama"
	"github.com/go-martini/martini"
)

// KafkaSession is the struct to keep Sarama
type KafkaSession struct {
	*sarama.Config
	kafkaBrokerConnStrList []string
	syncProducer           sarama.SyncProducer
}

// NewKafkaSession connects to the Kafka brokers.
func NewKafkaSession(connStrList []string) *KafkaSession {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	brokers := connStrList
	producer, err := sarama.NewSyncProducer(brokers, config)

	if err != nil {
		panic(err)
	}

	return &KafkaSession{config, connStrList, producer}
}

// ProduceMessage produces a message in the topic specified
func (kafka *KafkaSession) ProduceMessage(content string, topic string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(content),
	}

	// partition, offset, err := kafka.syncProducer.SendMessage(msg)
	_, _, err := kafka.syncProducer.SendMessage(msg)
	if err != nil {
		return err
	}

	// fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)

	return nil
}

// NewKafkaHandler adds Kafka to the Martini pipeline
func (kafka *KafkaSession) NewKafkaHandler() martini.Handler {
	return func(context martini.Context) {
		context.Map(kafka)
		context.Next()
	}
}
