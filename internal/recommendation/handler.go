package recommendation

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *RecommendationService
}

func NewHandler(service *RecommendationService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetRecommendations(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login required"})
		return
	}

	recommendations, err := h.service.GetRecommendations(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recommendations"})
		return
	}
	c.JSON(http.StatusOK, recommendations)
}
