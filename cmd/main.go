package main

import (
	"database/sql"
	"github.com/AdiInfiniteLoop/Authora/handlers"
	"github.com/AdiInfiniteLoop/Authora/internal/config"
	"github.com/AdiInfiniteLoop/Authora/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

func main() {
	//Initialize Redis Here
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//pong, err := client.Ping(ctx).Result()
	//log.Println(pong, err)

	//Initialize the database here
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Cannot find the env ")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Println("Cannot find the database url")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Println("Cannot connect to the database")
	}

	var testQuery int
	err = conn.QueryRow("SELECT 1").Scan(&testQuery)
	if err != nil {
		log.Println("Error While Querying the database")
	} else {
		log.Println(testQuery)
		log.Println("Connection Successful !!! Test Query ran successfully")
	}

	//Setting Api Configuration
	//apiConfig is a struct type that stores the api configuration like db_url, authentication, apiKey, etc
	apiConfig := &config.ApiConfig{
		DB:          database.New(conn),
		RedisClient: client,
	}

	LocalApiConfig := handlers.LocalApiConfig{
		ApiConfig: apiConfig,
	}

	//Initialize the routers
	router := gin.Default() //Sets up the router

	//Cors

	authorized := router.Group("/")
	authorized.Use(LocalApiConfig.AuthMiddleware())
	{
		authorized.GET("/health-check", LocalApiConfig.HandlerCheckReadiness)
		authorized.GET("/auth-route", LocalApiConfig.HandlerAuthRoute)
	}

	router.POST("/sign-in", LocalApiConfig.SignInHandler)
	router.POST("/logout", LocalApiConfig.LogoutHandler)
	router.POST("/sign-up", LocalApiConfig.CreateUserHandler)
	log.Fatal(router.Run(":8080"))
}
