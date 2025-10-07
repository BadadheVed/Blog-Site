package function

import (
	"log"

	"github.com/google/uuid"
	"github.com/yourname/blog-kafka/config"
	"github.com/yourname/blog-kafka/models"
)

func CreateNotificationForChannelMembers(channelID, authorID, blogID uuid.UUID, blogTitle string, notifType models.NotificationType) error {
	var members []models.ChannelMember

	if err := config.DB.Where("channel_id = ? AND user_id != ?", channelID, authorID).Find(&members).Error; err != nil {
		return err
	}

	for _, member := range members {
		notif := models.Notification{
			ID:      uuid.New(),
			Type:    notifType,
			Message: "Blog " + string(notifType) + ": " + blogTitle,
			UserID:  member.UserID,
			BlogID:  blogID,
		}

		if err := config.DB.Create(&notif).Error; err != nil {
			log.Println("Failed to create notification for user:", member.UserID, "err:", err)
		}
	}

	return nil
}
