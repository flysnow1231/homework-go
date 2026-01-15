package handler

import (
	"net/http"

	"blog/internal/pkg/resp"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	ReadyFn func() error
}

func (h *HealthHandler) Healthz(c *gin.Context) {
	resp.OK(c, gin.H{"status": "ok"})
}

func (h *HealthHandler) Readyz(c *gin.Context) {
	if h.ReadyFn != nil {
		if err := h.ReadyFn(); err != nil {
			resp.Fail(c, http.StatusServiceUnavailable, "not_ready")
			return
		}
	}
	resp.OK(c, gin.H{"status": "ready"})
}
