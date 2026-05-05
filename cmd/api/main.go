package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	applogger "reader-club/internal/application/logger"
	authuc "reader-club/internal/application/usecase/auth"
	bookuc "reader-club/internal/application/usecase/book"
	bookclubuc "reader-club/internal/application/usecase/bookclub"
	meetinguc "reader-club/internal/application/usecase/meeting"
	membershipuc "reader-club/internal/application/usecase/membership"
	themeuc "reader-club/internal/application/usecase/theme"
	"reader-club/internal/infrastructure/auth"
	"reader-club/internal/infrastructure/bootstrap"
	"reader-club/internal/infrastructure/persistence"
	httpserver "reader-club/internal/interfaces/http"
	"reader-club/internal/interfaces/http/controller"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	logger := applogger.NewStdLogger()

	dbHost := getenv("DB_HOST", "localhost")
	dbPort := getenv("DB_PORT", "5432")
	dbName := getenv("DB_NAME", "reader_club")
	log.Printf("connecting to database host=%s port=%s db=%s", dbHost, dbPort, dbName)

	port, _ := strconv.Atoi(dbPort)
	db, err := persistence.Connect(persistence.Config{
		Host:     dbHost,
		Port:     port,
		User:     getenv("DB_USER", "postgres"),
		Password: getenv("DB_PASSWORD", "postgres"),
		DBName:   dbName,
		SSLMode:  getenv("DB_SSLMODE", "disable"),
	})
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	if err := persistence.Migrate(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	userRepo := persistence.NewUserRepository(db)

	bootstrap.SeedDefaultUsers(context.Background(), userRepo)
	clubRepo := persistence.NewBookClubRepository(db)
	memberRepo := persistence.NewMembershipRepository(db)
	suggestionRepo := persistence.NewBookSuggestionRepository(db)
	themeRepo := persistence.NewMonthlyThemeRepository(db)
	meetingRepo := persistence.NewMeetingRepository(db)
	auditRepo := persistence.NewAuditLogRepository(db)

	tokenService := auth.NewJWTService(
		getenv("JWT_SECRET", "change-me-in-production"),
		24*time.Hour,
	)

	registerUC := authuc.NewRegister(userRepo)
	loginUC := authuc.NewLogin(userRepo, tokenService)
	createClubUC := bookclubuc.NewCreateBookClub(clubRepo, memberRepo, auditRepo, logger)
	getClubUC := bookclubuc.NewGetBookClub(clubRepo)
	listClubsUC := bookclubuc.NewListBookClubs(clubRepo)
	deleteClubUC := bookclubuc.NewDeleteBookClub(clubRepo, meetingRepo, themeRepo, suggestionRepo, memberRepo, auditRepo, logger)
	joinClubUC := membershipuc.NewJoinClub(clubRepo, memberRepo, userRepo)
	suggestBookUC := bookuc.NewSuggestBook(memberRepo, suggestionRepo, userRepo, auditRepo, logger)
	listSuggestUC := bookuc.NewListBookSuggestions(suggestionRepo)
	drawThemeUC := themeuc.NewDrawMonthlyTheme(memberRepo, suggestionRepo, themeRepo, auditRepo, logger)
	createMeetingUC := meetinguc.NewCreateMeeting(memberRepo, themeRepo, meetingRepo, auditRepo, logger)

	ctrls := httpserver.Controllers{
		Auth:       controller.NewAuthController(registerUC, loginUC),
		BookClub:   controller.NewBookClubController(createClubUC, getClubUC, listClubsUC, deleteClubUC),
		Membership: controller.NewMembershipController(joinClubUC),
		Book:       controller.NewBookSuggestionController(suggestBookUC, listSuggestUC),
		Theme:      controller.NewThemeController(drawThemeUC),
		Meeting:    controller.NewMeetingController(createMeetingUC),
	}

	router := httpserver.NewRouter(ctrls, tokenService, memberRepo)

	addr := fmt.Sprintf(":%s", getenv("PORT", "8080"))
	log.Printf("server listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
