package recommendation

import (
	"net/http"

	"app/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *RecommendationService
	log     *logger.Logger
}

func NewHandler(service *RecommendationService, log *logger.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
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
