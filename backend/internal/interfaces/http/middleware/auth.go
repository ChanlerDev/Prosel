package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/chanler/prosel/backend/internal/interfaces/http/response"
)

type AccessTokenParser interface {
	ParseAccessToken(token string) (userID string, err error)
}

const UserIDKey = "userID"

func Auth(tokens AccessTokenParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
			c.Abort()
			return
		}
		userID, err := tokens.ParseAccessToken(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token", nil)
			c.Abort()
			return
		}
		c.Set(UserIDKey, userID)
		c.Next()
	}
}

func CurrentUserID(c *gin.Context) string {
	value, _ := c.Get(UserIDKey)
	userID, _ := value.(string)
	return userID
}
