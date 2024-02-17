package handlers

import (
	"EM_test_task/internal/entities"
	"EM_test_task/pkg/storage"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type UserHandler struct {
	UserRepo storage.Repository
}

type userInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) Registration(c *gin.Context) {
	var input userInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("RegistrationUser BindJSON Error: ", err)
		return
	}

	_, err := h.UserRepo.GetUser(input.Username)
	if err == nil {
		err = errors.New("User already exists")
		c.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePrivate)
		slog.Error("Registration User already exists Error: ", err)
		return
	}

	userID, err := h.UserRepo.CreateUser(entities.User{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).SetType(gin.ErrorTypePrivate)
		slog.Error("RegistrationUser CreateUser Error: ", err)
		return
	}

	slog.Info("Person added successfully")
	c.JSON(http.StatusCreated, userID)
}

func (h *UserHandler) Login(c *gin.Context) {
	var input userInput
	jwtSecretKey := []byte(os.Getenv("JWT_SECRET"))

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		slog.Error("LoginUser BindJSON Error: ", err)
		return
	}

	user, err := h.UserRepo.GetUser(input.Username)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).SetType(gin.ErrorTypePrivate)
		slog.Error("Login GetUser Error: ", err)
		return
	}
	if user.Password != input.Password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or password is incorrect"})
		return
	}

	payload := jwt.MapClaims{
		"sub": user.Username,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).SetType(gin.ErrorTypePrivate)
		slog.Error("JWT SignedString Error: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": t})
}
