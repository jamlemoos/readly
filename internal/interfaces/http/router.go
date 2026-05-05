package http

import (
	"net/http"

	"reader-club/internal/domain/entity"
	"reader-club/internal/domain/repository"
	"reader-club/internal/infrastructure/auth"
	"reader-club/internal/interfaces/http/controller"
	"reader-club/internal/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	Auth       *controller.AuthController
	BookClub   *controller.BookClubController
	Membership *controller.MembershipController
	Book       *controller.BookSuggestionController
	Theme      *controller.ThemeController
	Meeting    *controller.MeetingController
}

func NewRouter(ctrls Controllers, tokens auth.TokenService, members repository.MembershipRepository) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	api := r.Group("/api/v1")

	api.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    "BookClub API",
			"version": "1.0.0",
			"status":  "ok",
		})
	})
	api.POST("/auth/register", ctrls.Auth.Register)
	api.POST("/auth/login", ctrls.Auth.Login)

	secured := api.Group("/")
	secured.Use(middleware.Authenticate(tokens))

	secured.POST("/clubs", middleware.RequireGlobalRoles(string(entity.RoleAdmin)), ctrls.BookClub.Create)
	secured.GET("/clubs", ctrls.BookClub.List)
	secured.GET("/clubs/:club_id", ctrls.BookClub.GetByID)

	adminOnly := secured.Group("/clubs/:club_id")
	adminOnly.Use(middleware.RequireClubRole(members, entity.RoleAdmin))
	adminOnly.DELETE("", ctrls.BookClub.Delete)
	adminOnly.POST("/themes/draw", ctrls.Theme.Draw)
	adminOnly.POST("/meetings", ctrls.Meeting.Create)

	secured.POST("/clubs/:club_id/join", ctrls.Membership.Join)

	memberOrAdmin := secured.Group("/clubs/:club_id")
	memberOrAdmin.Use(middleware.RequireClubRole(members, entity.RoleAdmin, entity.RoleMember))
	memberOrAdmin.POST("/suggestions", ctrls.Book.Suggest)

	clubMember := secured.Group("/clubs/:club_id")
	clubMember.Use(middleware.RequireClubRole(members, entity.RoleAdmin, entity.RoleMember, entity.RoleVisitor))
	clubMember.GET("/suggestions", ctrls.Book.List)

	return r
}
