package handler

import (
	"blog/internal/model"
	"blog/internal/pkg/resp"
	"blog/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type PostHandler struct {
	svc *service.PostService
	log *zap.Logger
}

func NewPostHandler(svc *service.PostService, log *zap.Logger) *PostHandler {
	return &PostHandler{svc: svc, log: log}
}

//type createUserReq struct {
//	Name string `json:"name" binding:"required"`
//}

func (h *PostHandler) AddPost(c *gin.Context) {
	var post model.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.log.Info("Create handler called", zap.Any("post", post))

	if _, err := h.svc.CreatePost(c, &post); err != nil {
		h.log.Error("Create handler failed", zap.Any("user", post), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "post successfully"})
}

func (h *PostHandler) QueryPostByUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("userid"))
	page, _ := strconv.Atoi(c.Param("page"))
	size, _ := strconv.Atoi(c.Param("size"))

	posts, err := h.svc.QueryPosts(c, userId, page, size)
	if err != nil {
		h.log.Error("queryPost error", zap.Any("userid", userId), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queryPostByUser"})
		return
	}
	resp.OK(c, posts)
}
