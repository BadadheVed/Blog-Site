package workers

import (
	"log"

	"github.com/google/uuid"
	"github.com/yourname/blog-kafka/config"
	"github.com/yourname/blog-kafka/models"
)

func (n *NotificationJob) Process() error {
	notif := models.Notification{
		ID:      uuid.New(),
		Type:    n.Type,
		Message: "Blog " + string(n.Type) + ": " + n.BlogTitle,
		UserID:  n.UserID,
		BlogID:  n.BlogID,
	}
	if err := config.DB.Create(&notif).Error; err != nil {
		log.Println("❌ Failed to insert notification for user:", n.UserID, "error:", err)
		return err
	}
	log.Println("✅ Notification created for user:", n.UserID)
	return nil
}
