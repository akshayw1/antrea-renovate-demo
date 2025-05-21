package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/akshayw1/antrea-renovate-demo/internal/auth"
	"github.com/akshayw1/antrea-renovate-demo/internal/storage"
)

// APIServer encapsulates the API server components
type APIServer struct {
	storage *storage.MemoryStorage
}

// NewAPIServer creates a new API server instance
func NewAPIServer() *APIServer {
	return &APIServer{
		storage: storage.NewMemoryStorage(),
	}
}

// Setup configures the routes on the given router
func (s *APIServer) Setup(router *gin.Engine) {
	// Public routes
	router.POST("/login", s.handleLogin)

	// Protected routes
	authorized := router.Group("/api")
	authorized.Use(s.authMiddleware())
	{
		authorized.GET("/items", s.listItems)
		authorized.GET("/items/:id", s.getItem)
		authorized.POST("/items", s.createItem)
		authorized.DELETE("/items/:id", s.deleteItem)
	}
}

// authMiddleware validates JWT tokens
func (s *APIServer) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		username, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store username in context
		c.Set("username", username)
		c.Next()
	}
}

// handleLogin generates a JWT token for a user
func (s *APIServer) handleLogin(c *gin.Context) {
	var credentials struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real app, validate credentials against a database
	// For this demo, accept any username/password
	
	tokenString, err := auth.GenerateToken(credentials.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

// listItems returns all stored items
func (s *APIServer) listItems(c *gin.Context) {
	items := s.storage.ListItems()
	c.JSON(http.StatusOK, items)
}

// getItem retrieves a specific item by ID
func (s *APIServer) getItem(c *gin.Context) {
	id := c.Param("id")
	
	item, exists := s.storage.GetItem(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	
	c.JSON(http.StatusOK, item)
}

// createItem adds a new item
func (s *APIServer) createItem(c *gin.Context) {
	var item storage.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// In a real app, validate and generate ID if not provided
	if item.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}
	
	s.storage.StoreItem(item)
	c.JSON(http.StatusCreated, item)
}

// deleteItem removes an item by ID
func (s *APIServer) deleteItem(c *gin.Context) {
	id := c.Param("id")
	
	if deleted := s.storage.DeleteItem(id); !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}