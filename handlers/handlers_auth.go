package handlers

import (
	"encoding/json"
	"github.com/AdiInfiniteLoop/Authora/models"
	"github.com/AdiInfiniteLoop/Authora/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"time"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type JWTOutput struct {
	Token  string    `json:"token"`
	Expire time.Time `json:"expires"`
}

type SessionData struct {
	Token  string    `json:"token"`
	UserID uuid.UUID `json:"user_id"`
}

func (lac *LocalApiConfig) SignInHandler(c *gin.Context) {
	var userToAuth models.User //data provided by the user

	if err := c.ShouldBindJSON(&userToAuth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	var val = utils.ValidationOfUser(userToAuth)
	if len(val) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  val,
		})
		return
	}
	//Fetch the user from database and check if user exists or not
	foundUser, err := lac.DB.FindUserByEmail(c, userToAuth.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "No such user found",
		})
		return
	}

	if foundUser.Password != userToAuth.Password {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Unauthorized",
		})
		return
	}

	expirationTime := time.Now().Add(10 * time.Minute)
	//Creating a JWT Assignment
	claims := &Claims{
		Email: foundUser.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	// Whenever you're working with JWT signing in Go (using the standard jwt-go library),
	// you'll need to convert string secrets to byte slices.
	if err != nil {
		log.Println("Error while Signing Up the String", err)
		return
	}

	//Create a sessionId
	sessionId := uuid.New().String()
	//Create a session interface
	sessionData := map[string]interface{}{
		"token":  tokenString,
		"userId": foundUser.ID,
	}

	//Marshal the interface into JSON
	sessionDataJSON, err := json.Marshal(sessionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to encode the Session Data into JSON",
		})
		return
	}

	//Create a Redis Mapping for sessionId -> sessionData
	err = lac.RedisClient.Set(c, sessionId, sessionDataJSON, time.Until(expirationTime)).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed Redis Setting",
		})
		return
	}
	c.SetCookie("session_id", sessionId, int(time.Until(expirationTime).Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status":    "Success",
		"message":   "Found User By Email",
		"expiresAt": expirationTime,
	})
}

func (lac *LocalApiConfig) LogoutHandler(c *gin.Context) {
	//Find the id and check if it exists
	//remove from the session Now
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Unauthorized",
		})
		return
	}

	err = lac.RedisClient.Del(c, sessionID).Err()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Unauthorized",
		})
		return
	}
	//Set the  cookie to empty
	c.SetCookie("session_id", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Log Out Successful",
	})
	return
}

// middleware must always return a gin.handler{}

func (lac *LocalApiConfig) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get the sessionId from cookie
		//Get the Session Data as JSON from Redis
		//Convert the Session Data into golang struct
		//Get the Token and User ID from the unmarshalled session data
		//Check the availability then store in request Context(e.g., req.user)

		sessionId, err := c.Cookie("session_id")

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "unauthorized - no session",
			})
			return
		}

		sessionDataJSON, err := lac.RedisClient.Get(c, sessionId).Result()

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "unauthorized - no redis session data",
			})
			return
		}
		var sessionData SessionData
		err = json.Unmarshal([]byte(sessionDataJSON), &sessionData)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  "cannot unmarshal data",
			})
			return
		}
		parsedId, err := uuid.Parse(sessionData.UserID.String())
		sessionData.UserID = parsedId
		token, err := jwt.ParseWithClaims(sessionData.Token, &Claims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "invalid token",
			})
			return
		}
		c.Set("userId", sessionData.UserID)
		c.Next()
	}
}

func (lac *LocalApiConfig) HandlerAuthRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Authenticated Route working well",
	})
}
