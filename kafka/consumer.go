package kafka

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

// StartConsumer starts a Sarama consumer group that calls handler for each message.
// Do NOT defer cancel() here; the goroutine will cancel on OS signal.
func StartConsumer(brokers []string, groupID, topic string, handler func(message string)) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("[kafka] Error creating consumer group client: %v", err)
	}

	// Context we can cancel on SIGINT/SIGTERM
	ctx, cancel := context.WithCancel(context.Background())

	// Listen for termination signals and cancel context when received
	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
		<-sigterm
		log.Println("[kafka] shutdown signal received, canceling consumer context")
		cancel()
	}()

	// Run consumer loop in background goroutine
	go func() {
		for {
			handlerWrapper := &consumerGroupHandler{handler: handler}
			if err := client.Consume(ctx, []string{topic}, handlerWrapper); err != nil {
				// log errors and continue; if context canceled, exit loop
				log.Printf("[kafka] consumer error: %v", err)
			}
			if ctx.Err() != nil {
				// context canceled -> exit goroutine
				log.Println("[kafka] consumer context canceled, exiting consumer loop")
				return
			}
		}
	}()
}

type consumerGroupHandler struct {
	handler func(message string)
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.handler(string(msg.Value))
		sess.MarkMessage(msg, "")
	}
	return nil
}
