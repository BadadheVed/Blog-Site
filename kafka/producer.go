package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type Producer struct {
	SyncProducer sarama.SyncProducer
}

func NewProducer(brokers []string) *Producer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	prod, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start producer: %v", err)
	}

	return &Producer{SyncProducer: prod}
}

func (p *Producer) Send(topic, message string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := p.SyncProducer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}

	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
}
