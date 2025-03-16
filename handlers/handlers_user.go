package handlers

import (
	"github.com/AdiInfiniteLoop/Authora/internal/config"
	"github.com/AdiInfiniteLoop/Authora/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type LocalApiConfig struct {
	//Composition relationship ("has a" relationship)
	*config.ApiConfig
}

func (lac *LocalApiConfig) CreateUserHandler(c *gin.Context) {
	type CreateUserParams struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	user := CreateUserParams{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Bad Request",
		})
		return
	}

	newUser, err := lac.DB.CreateUser(c, database.CreateUserParams{
		ID:       uuid.New(),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	})

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			c.JSON(http.StatusConflict, gin.H{
				"status": "error",
				"error":  "User already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Cannot Create this User",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Created A User successfully",
		"data":    newUser,
	})

	return
}
