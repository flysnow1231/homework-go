package handler

import (
	"blog/internal/model"
	"blog/internal/pkg/resp"
	"blog/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserHandler struct {
	svc *service.UserService
	log *zap.Logger
}

func NewUserHandler(svc *service.UserService, log *zap.Logger) *UserHandler {
	return &UserHandler{svc: svc, log: log}
}

//type createUserReq struct {
//	Name string `json:"name" binding:"required"`
//}

func (h *UserHandler) Create(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.log.Info("Create handler called", zap.Any("user", user))

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	if _, err := h.svc.CreateUser(c, &user); err != nil {
		h.log.Error("Create handler failed", zap.Any("user", user), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
	username := c.Param("username")

	pwd := c.Param("password")
	if len(username) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return

	}

	if len(pwd) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return

	}

	u, err := h.svc.GetUser(c.Request.Context(), username, pwd)
	if err != nil {
		resp.Fail(c, http.StatusNotFound, "not_found")
		return

	}

	resp.OK(c, u)
}
