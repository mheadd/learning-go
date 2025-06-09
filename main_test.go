package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
)

type testUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func setupTestDB() *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=usersdb sslmode=disable")
	if err != nil {
		panic(err)
	}
	// Run init.sql to ensure schema exists
	initSQL, err := os.ReadFile("init.sql")
	if err != nil {
		panic("Failed to read init.sql: " + err.Error())
	}
	_, err = db.Exec(string(initSQL))
	if err != nil {
		panic("Failed to execute init.sql: " + err.Error())
	}
	// Clean up users table before each test
	_, _ = db.Exec("DELETE FROM users")
	return db
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) { c.File("./static/index.html") })
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	api := r.Group("/api")
	{
		api.GET("/users", getUsers)
		api.POST("/users", createUser)
	}
	return r
}

func TestHealthEndpoint(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "healthy", resp["status"])
}

func TestCreateAndGetUser(t *testing.T) {
	db = setupTestDB()
	defer db.Close()
	r := setupRouter()

	// Create user
	user := testUser{ID: "test1", Name: "Test User"}
	body, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	// Get users
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/users", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var resp map[string][]testUser
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp["users"], 1)
	assert.Equal(t, user.ID, resp["users"][0].ID)
	assert.Equal(t, user.Name, resp["users"][0].Name)
}
