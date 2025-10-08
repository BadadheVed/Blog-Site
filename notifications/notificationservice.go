package notifications

import (
	"log"

	"github.com/google/uuid"
	"github.com/yourname/blog-kafka/config"
	"github.com/yourname/blog-kafka/models"
)

type NotificationService struct {
	WorkerPool *WorkerPool
}

func NewNotificationService(wp *WorkerPool) *NotificationService {
	return &NotificationService{WorkerPool: wp}
}

func (n *NotificationService) CreateNotificationForChannelMembers(
	channelID, authorID, blogID uuid.UUID,
	blogTitle string,
	notifType models.NotificationType,
) error {
	rows, err := config.DB.Model(&models.ChannelMember{}).
		Select("user_id").
		Where("channel_id = ? AND user_id != ?", channelID, authorID).
		Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var member models.ChannelMember
		if err := config.DB.ScanRows(rows, &member); err != nil {
			log.Printf("Error scanning member: %v", err)
			continue
		}

		job := NotificationJob{
			UserID:    member.UserID,
			BlogID:    blogID,
			BlogTitle: blogTitle,
			Type:      notifType,
		}
		n.WorkerPool.Submit(job)
		count++
	}

	log.Printf("Streamed and submitted %d notification jobs for '%s'\n", count, blogTitle)
	return nil
}

func (n *NotificationService) CreateSingleNotification(job NotificationJob) error {
	notif := models.Notification{
		ID:      uuid.New(),
		Type:    job.Type,
		Message: "Blog " + string(job.Type) + ": " + job.BlogTitle,
		UserID:  job.UserID,
		BlogID:  job.BlogID,
	}
	return config.DB.Create(&notif).Error
}
