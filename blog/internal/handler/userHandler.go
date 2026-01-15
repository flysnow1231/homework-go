package handler

import (
	"blog/internal/model"
	"blog/internal/pkg/resp"
	"blog/internal/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler { return &UserHandler{svc: svc} }

type createUserReq struct {
	Name string `json:"name" binding:"required"`
}

func (h *UserHandler) Create(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	if _, err := h.svc.CreateUser(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *UserHandler) Get(c *gin.Context) {
	username := c.Param("username")
	if len(username) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})

	}

	u, err := h.svc.GetUser(c.Request.Context(), username)
	if err != nil {
		resp.Fail(c, http.StatusNotFound, "not_found")

	}
	resp.OK(c, u)
}
