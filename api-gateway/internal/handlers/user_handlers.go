package handlers

import (
	"context"
	"net/http"

	userProto "github.com/abin-saji-2003/MicroService-Clean-Architecture-Go-/user-service/api/proto"
	"github.com/gin-gonic/gin"
)

// ✅ Register User Routes
func RegisterUserRoutes(r *gin.Engine, userClient userProto.UserServiceClient) {
	userRoutes := r.Group("/api/users")
	{
		userRoutes.POST("/register", func(c *gin.Context) {
			var req userProto.RegisterUserRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			resp, err := userClient.RegisterUser(context.Background(), &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": resp.Message})
		})

		userRoutes.POST("/login", func(c *gin.Context) {
			var req userProto.LoginUserRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			resp, err := userClient.LoginUser(context.Background(), &req)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"user_id": resp.UserId, "name": resp.Name, "email": resp.Email})
		})
	}
}
