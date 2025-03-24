package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"monji/internal/database"
	"monji/internal/models"
	"monji/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware verifies JWT tokens and loads the user into context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}
		// Use strings.Split instead of jwt.SplitToken
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}
		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(utils.JWTSecret), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, ok := claims["user_id"].(float64)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				return
			}
			var user models.User
			row := database.DB.QueryRow(
				`SELECT id, first_name, last_name, email, company, password, role FROM users WHERE id = ?`,
				int(userID))
			if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Company, &user.Password, &user.Role); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				return
			}
			// Log the loaded user for debugging.
			fmt.Printf("AuthMiddleware: loaded user %+v\n", user)
			user.Password = ""
			c.Set("user", user)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}
	}
}

// AdminMiddleware ensures the user is an admin or superadmin.
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - user not in context"})
			return
		}
		usr, ok := user.(models.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - user type assertion failed"})
			return
		}
		// Log the role for debugging.
		fmt.Printf("AdminMiddleware: user role = %s\n", usr.Role)
		if usr.Role != "admin" && usr.Role != "superadmin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden - Admins Only"})
			return
		}
		c.Next()
	}
}
