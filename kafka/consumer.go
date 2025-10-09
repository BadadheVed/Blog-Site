package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/yourname/blog-kafka/notifications"
)

var MyWorkerPool *notifications.WorkerPool
var MyNotifService *notifications.NotificationService

func SetWorkerPoolAndService(wp *notifications.WorkerPool, ns *notifications.NotificationService) {
	MyWorkerPool = wp
	MyNotifService = ns
}

func StartConsumer(brokers []string, groupID, topic string) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_8_0_0

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("[kafka] Error creating consumer group client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
		<-sigterm
		log.Println("[kafka] shutdown signal received, canceling consumer context")
		cancel()
	}()

	go func() {
		for {
			handlerWrapper := &consumerGroupHandler{}
			if err := client.Consume(ctx, []string{topic}, handlerWrapper); err != nil {
				log.Printf("[kafka] consumer error: %v", err)
			}
			if ctx.Err() != nil {
				log.Println("[kafka] consumer context canceled, exiting consumer loop")
				return
			}
		}
	}()
}

type consumerGroupHandler struct{}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if MyWorkerPool == nil || MyNotifService == nil {
			log.Println("[worker-pool] or NotificationService not initialized, skipping message")
			continue
		}

		var payload NotificationPayload
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			log.Printf("[kafka] failed to unmarshal message: %v", err)
			continue
		}

		err := MyNotifService.CreateNotificationForChannelMembers(
			payload.ChannelID,
			payload.AuthorID,
			payload.BlogID,
			payload.BlogTitle,
			payload.Type,
		)
		if err != nil {
			log.Printf("[kafka] failed to create notifications: %v", err)
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
