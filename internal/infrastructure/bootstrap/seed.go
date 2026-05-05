package bootstrap

import (
	"context"
	"log"
	"os"

	"reader-club/internal/domain/entity"
	"reader-club/internal/domain/repository"
)

type seedCandidate struct {
	email    string
	password string
	name     string
	role     entity.Role
}

// SeedDefaultUsers creates seed users from environment variables.
// Missing env vars or pre-existing accounts are silently skipped.
func SeedDefaultUsers(ctx context.Context, users repository.UserRepository) {
	candidates := []seedCandidate{
		{
			email:    os.Getenv("ADMIN_EMAIL"),
			password: os.Getenv("ADMIN_PASSWORD"),
			name:     envOr("ADMIN_NAME", "System Admin"),
			role:     entity.RoleAdmin,
		},
		{
			email:    os.Getenv("VISITOR_EMAIL"),
			password: os.Getenv("VISITOR_PASSWORD"),
			name:     envOr("VISITOR_NAME", "Read Only User"),
			role:     entity.RoleVisitor,
		},
	}

	for _, c := range candidates {
		if c.email == "" || c.password == "" {
			continue
		}
		seedOne(ctx, users, c)
	}
}

func seedOne(ctx context.Context, users repository.UserRepository, c seedCandidate) {
	exists, err := users.ExistsByEmail(ctx, c.email)
	if err != nil {
		log.Printf("[WARN] seed: could not check existence of %s: %v", c.email, err)
		return
	}
	if exists {
		log.Printf("[INFO] seed: %s already exists, skipping", c.email)
		return
	}

	user, err := entity.NewUser(c.name, c.email, c.password)
	if err != nil {
		log.Printf("[WARN] seed: could not build user %s: %v", c.email, err)
		return
	}
	user.GlobalRole = c.role

	if err := users.Save(ctx, user); err != nil {
		log.Printf("[WARN] seed: could not save user %s: %v", c.email, err)
		return
	}

	log.Printf("[INFO] seed: created %s with role %s", c.email, c.role)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
