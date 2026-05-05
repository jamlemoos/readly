package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"reader-club/internal/infrastructure/auth"
	"reader-club/internal/interfaces/http/httperr"
)

const (
	ContextUserID = "actor_id"
	ContextEmail  = "actor_email"
	ContextRoles  = "actor_roles"
)

func Authenticate(tokens auth.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			httperr.Abort(c, http.StatusUnauthorized, "unauthorized", "missing or malformed token")
			return
		}

		claims, err := tokens.Validate(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			httperr.Abort(c, http.StatusUnauthorized, "unauthorized", "invalid or expired token")
			return
		}

		actorID, err := uuid.Parse(claims.UserID)
		if err != nil {
			httperr.Abort(c, http.StatusUnauthorized, "unauthorized", "invalid token subject")
			return
		}

		c.Set(ContextUserID, actorID)
		c.Set(ContextEmail, claims.Email)
		c.Set(ContextRoles, claims.Roles)
		c.Next()
	}
}

func ActorID(c *gin.Context) uuid.UUID {
	return c.MustGet(ContextUserID).(uuid.UUID)
}

func ActorRoles(c *gin.Context) []string {
	roles, _ := c.Get(ContextRoles)
	if r, ok := roles.([]string); ok {
		return r
	}
	return nil
}
