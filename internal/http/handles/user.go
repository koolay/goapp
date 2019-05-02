package handles

import (
	"context"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetProfile(c *gin.Context) {

	ctx := context.Background()
	userProfile, err := h.userService.GetProfile(ctx, h.GetSession(c).UserID)
	if !h.CheckError(c, err) {
		return
	}
	h.SendSuccess(c, 200, userProfile)
}
