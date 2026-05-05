package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"reader-club/internal/domain/entity"
	"reader-club/internal/domain/repository"
	"reader-club/internal/interfaces/http/httperr"
)

// RequireClubRole aborts with 403 if the actor does not hold one of the allowed roles in the club.
// Expects :club_id path param and actor_id context value (set by Authenticate).
func RequireClubRole(members repository.MembershipRepository, allowed ...entity.Role) gin.HandlerFunc {
	roleSet := make(map[entity.Role]struct{}, len(allowed))
	for _, r := range allowed {
		roleSet[r] = struct{}{}
	}

	return func(c *gin.Context) {
		actorID := ActorID(c)
		clubID, err := uuid.Parse(c.Param("club_id"))
		if err != nil {
			httperr.Abort(c, http.StatusBadRequest, "validation_error", "invalid club_id")
			return
		}

		membership, err := members.FindByUserAndClub(context.Background(), actorID, clubID)
		if err != nil {
			httperr.Abort(c, http.StatusForbidden, "forbidden", "not a member of this club")
			return
		}

		if _, ok := roleSet[membership.Role]; !ok {
			httperr.Abort(c, http.StatusForbidden, "forbidden", "insufficient club role")
			return
		}

		c.Next()
	}
}

// RequireGlobalRoles aborts with 403 if the actor's JWT roles do not include at least one allowed role.
func RequireGlobalRoles(allowed ...string) gin.HandlerFunc {
	roleSet := make(map[string]struct{}, len(allowed))
	for _, r := range allowed {
		roleSet[r] = struct{}{}
	}

	return func(c *gin.Context) {
		for _, role := range ActorRoles(c) {
			if _, ok := roleSet[role]; ok {
				c.Next()
				return
			}
		}
		httperr.Abort(c, http.StatusForbidden, "forbidden", "insufficient global role")
	}
}
