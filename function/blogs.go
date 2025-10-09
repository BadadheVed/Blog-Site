package function

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourname/blog-kafka/config"
	"github.com/yourname/blog-kafka/kafka"

	middleware "github.com/yourname/blog-kafka/middlewares"
	"github.com/yourname/blog-kafka/models"
	"github.com/yourname/blog-kafka/notifications"

	"gorm.io/gorm"
)

var MyWorkerPool *notifications.WorkerPool

func SetWorkerPool(wp *notifications.WorkerPool) {
	MyWorkerPool = wp
}

func CreateBlog(c *gin.Context) {
	userIDStr, _, ok := middleware.ExtractUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	// Get Channel ID from route parameter
	channelIDStr := c.Param("channelId")
	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid channel id"})
		return
	}

	// Verify that channel exists
	var channel models.Channel
	if err := config.DB.First(&channel, "id = ?", channelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
		return
	}

	// Parse input
	var input struct {
		Title string  `json:"title" binding:"required"`
		Body  *string `json:"body"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title cannot be empty"})
		return
	}

	// Create Blog
	blog := models.Blog{
		ID:        uuid.New(),
		Title:     input.Title,
		Body:      input.Body,
		AuthorId:  userID,
		ChannelID: channelID,
	}

	if err := config.DB.Create(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload the blog with related Author and Channel data
	if err := config.DB.Preload("Author").Preload("Channel").First(&blog, "id = ?", blog.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load relations"})
		return
	}
	payload := kafka.NotificationPayload{
		ChannelID: channelID,
		AuthorID:  userID,
		BlogID:    blog.ID,
		BlogTitle: blog.Title,
		Type:      models.NotificationTypeNew,
	}

	if err := kafka.PublishNotification(payload); err != nil {
		log.Printf("Failed to publish notification to Kafka: %v", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "blog created successfully under channel",
		"blog":    blog,
	})

}

func EditBlog(c *gin.Context) {
	userIDStr, _, ok := middleware.ExtractUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User Id"})
		return
	}

	blogIDStr := c.Param("id")
	blogID, err := uuid.Parse(blogIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog id"})
		return
	}

	var blog models.Blog
	if err := config.DB.First(&blog, "id = ?", blogID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if blog.AuthorId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not authorized to edit this blog"})
		return
	}

	var input struct {
		Body *string `json:"body"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&blog).Update("body", input.Body).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update blog body"})
		return
	}
	payload := kafka.NotificationPayload{
		ChannelID: blog.ChannelID,
		AuthorID:  userID,
		BlogID:    blog.ID,
		BlogTitle: blog.Title,
		Type:      models.NotificationTypeEdited,
	}

	if err := kafka.PublishNotification(payload); err != nil {
		log.Printf("Failed to publish notification to Kafka: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "blog body updated successfully"})
}

func DeleteBlog(c *gin.Context) {
	userIDStr, _, ok := middleware.ExtractUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	blogIDStr := c.Param("id")
	blogID, err := uuid.Parse(blogIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog id"})
		return
	}

	var blog models.Blog
	if err := config.DB.First(&blog, "id = ?", blogID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if blog.AuthorId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not authorized to delete this blog"})
		return
	}

	if err := config.DB.Delete(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete blog"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "blog deleted successfully"})

}

func GetBlogs(c *gin.Context) {
	userIDStr, _, ok := middleware.ExtractUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var channelIDs []uuid.UUID
	if err := config.DB.Model(&models.ChannelMember{}).
		Where("user_id = ?", userID).
		Pluck("channel_id", &channelIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user channels"})
		return
	}

	if len(channelIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"page":  page,
			"limit": limit,
			"total": 0,
			"blogs": []models.Blog{},
		})
		return
	}

	var total int64
	if err := config.DB.Model(&models.Blog{}).
		Where("channel_id IN ?", channelIDs).
		Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count blogs"})
		return
	}

	var blogs []models.Blog
	if err := config.DB.Preload("Author").Preload("Channel").
		Where("channel_id IN ?", channelIDs).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&blogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch blogs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"total": total,
		"blogs": blogs,
	})
}

func GetBlogById(c *gin.Context) {
	userIDStr, _, ok := middleware.ExtractUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	_, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	blogIDStr := c.Param("id")
	blogID, err := uuid.Parse(blogIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog id"})
		return
	}

	var blog models.Blog
	if err := config.DB.Preload("Author").Preload("Channel").First(&blog, "id = ?", blogID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch blog"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blog": blog})
}
