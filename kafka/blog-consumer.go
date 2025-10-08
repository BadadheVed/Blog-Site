package kafka

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/yourname/blog-kafka/models"
	"github.com/yourname/blog-kafka/notifications"
)

// BlogEvent is the message shape produced by your API when a blog is created/edited
type BlogEvent struct {
	ChannelID uuid.UUID               `json:"channel_id"`
	AuthorID  uuid.UUID               `json:"author_id"`
	BlogID    uuid.UUID               `json:"blog_id"`
	BlogTitle string                  `json:"blog_title"`
	Type      models.NotificationType `json:"type"`
}

// StartBlogEventConsumer wires the generic consumer to your NotificationService.
// It unmarshals each message into BlogEvent and calls CreateNotificationForChannelMembers.
func StartBlogEventConsumer(notifSvc *notifications.NotificationService, brokers []string, groupID, topic string) {
	StartConsumer(brokers, groupID, topic, func(message string) {
		var ev BlogEvent
		if err := json.Unmarshal([]byte(message), &ev); err != nil {
			log.Printf("[kafka] invalid blog event payload: %v; raw=%s", err, message)
			return
		}

		log.Printf("[kafka] blog event received: blog_id=%s channel_id=%s", ev.BlogID.String(), ev.ChannelID.String())

		if notifSvc == nil {
			log.Println("[kafka] notification service is nil; cannot handle blog event")
			return
		}

		if err := notifSvc.CreateNotificationForChannelMembers(ev.ChannelID, ev.AuthorID, ev.BlogID, ev.BlogTitle, ev.Type); err != nil {
			log.Printf("[kafka] failed to create notifications for blog %s: %v", ev.BlogID.String(), err)
		}
	})
}
