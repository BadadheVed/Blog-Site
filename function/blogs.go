package function

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourname/blog-kafka/config"
	middleware "github.com/yourname/blog-kafka/middlewares"
	"github.com/yourname/blog-kafka/models"
	"gorm.io/gorm"
)

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

	var input struct {
		Title string  `json:"title" binding:"required"`
		Body  *string `json:"body"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blog := models.Blog{
		ID:       uuid.New(),
		Title:    input.Title,
		Body:     input.Body,
		AuthorId: userID,
	}

	if err := config.DB.Create(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create blog"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "blog created successfully", "blog": blog})
}

func EditBlog(c *gin.Context) {
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

func getBlogs(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")
}
