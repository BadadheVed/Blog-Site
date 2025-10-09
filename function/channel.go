package function

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourname/blog-kafka/config"
	middleware "github.com/yourname/blog-kafka/middlewares"
	"github.com/yourname/blog-kafka/models"
	"gorm.io/gorm"
)

func CreateChannel(c *gin.Context) {

	userIDStr, _, ok := middleware.ExtractUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id in token"})
		return
	}

	var input struct {
		Name        string  `json:"name" binding:"required"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	channel := models.Channel{
		ID:          uuid.New(),
		Name:        input.Name,
		Description: input.Description,
	}

	if err := config.DB.Create(&channel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Channel already exists"})
		return
	}

	sub := models.ChannelMember{
		ID:        uuid.New(),
		UserID:    userID,
		ChannelID: channel.ID,
	}

	if err := config.DB.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to subscribe creator"})
		return
	}

	if err := config.DB.Preload("User").Preload("Channel").First(&sub, "id = ?", sub.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load associations"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Channel created & subscribed",
		"subscription": sub,
	})
}

func SubscribeChannel(c *gin.Context) {

	userIDStr, _, ok := middleware.ExtractUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id in token"})
		return
	}

	channelIDStr := c.Param("id")
	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid channel id"})
		return
	}

	var channel models.Channel
	if err := config.DB.First(&channel, "id = ?", channelID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	var existing models.ChannelMember
	err = config.DB.Where("user_id = ? AND channel_id = ?", userID, channelID).First(&existing).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already subscribed"})
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	sub := models.ChannelMember{
		ID:        uuid.New(),
		UserID:    userID,
		ChannelID: channelID,
	}
	if err := config.DB.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to subscribe"})
		return
	}

	var size int64

	db := config.DB.Where("channel_id = ?", channel.ID).Count(&size)
	if db.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count channel members"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscribed successfully", "Channel Is": sub, "member_count": size})
}

func GetChannels(c *gin.Context) {
	userIDStr, _, ok := middleware.ExtractUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id in token"})
		return
	}

	var channels []models.Channel

	err = config.DB.
		Model(&models.Channel{}).
		Joins("JOIN channel_members ON channel_members.channel_id = channels.id").
		Where("channel_members.user_id = ?", userID).
		Find(&channels).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch channels"})
		return
	}

	if len(channels) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "you are not subscribed to any channels", "channels": []models.Channel{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "fetched subscribed channels successfully",
		"channels": channels,
	})
}

func JoinChannel(c *gin.Context) {
	userIDStr, _, ok := middleware.ExtractUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id in token"})
		return
	}

	channelIDStr := c.Param("channelId")
	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid channel id"})
		return
	}

	var channel models.Channel
	if err := config.DB.First(&channel, "id = ?", channelID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	var existing models.ChannelMember
	err = config.DB.Where("user_id = ? AND channel_id = ?", userID, channelID).First(&existing).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already joined this channel"})
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	sub := models.ChannelMember{
		ID:        uuid.New(),
		UserID:    userID,
		ChannelID: channelID,
	}
	if err := config.DB.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to join channel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "joined channel successfully",
		"channel": sub,
	})
}
