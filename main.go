package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// User represents a user in our system
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Config struct {
	DBHost     string `json:"db_host"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
	DBPort     string `json:"db_port"`
	AppPort    string `json:"app_port"`
}

var db *sql.DB
var config Config

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Failed to open config.json: %v", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Failed to decode config.json: %v", err)
	}

	// Override with environment variables if set
	if v := os.Getenv("DB_HOST"); v != "" {
		config.DBHost = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		config.DBUser = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		config.DBPassword = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		config.DBName = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		config.DBPort = v
	}
	if v := os.Getenv("APP_PORT"); v != "" {
		config.AppPort = v
	}
}

func initDB() {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Read and execute SQL from init.sql
	initSQL, err := os.ReadFile("init.sql")
	if err != nil {
		log.Fatalf("Failed to read init.sql: %v", err)
	}
	_, err = db.Exec(string(initSQL))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	loadConfig()
	initDB()
	defer db.Close()

	// Initialize Gin router
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

	// Serve the landing page
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// Simple health check endpoint
	r.GET("/health", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  "Database not reachable",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// API routes group
	api := r.Group("/api")
	{
		// GET /api/users
		api.GET("/users", getUsers)
		// POST /api/users
		api.POST("/users", createUser)
	}

	// Run the server
	r.Run(":" + config.AppPort)
}

// createUser handles POST requests to create a new user
func createUser(c *gin.Context) {
	var newUser User

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Input validation
	if newUser.ID == "" || newUser.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID and Name are required"})
		return
	}
	if len(newUser.ID) > 50 || len(newUser.Name) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID or Name too long"})
		return
	}

	// Use prepared statement to prevent SQL injection
	stmt, err := db.Prepare("INSERT INTO users (id, name) VALUES ($1, $2)")
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(newUser.ID, newUser.Name)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": newUser,
	})
}

// getUsers handles GET requests to retrieve users
func getUsers(c *gin.Context) {
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		log.Printf("Failed to query users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			log.Printf("Failed to scan user row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
