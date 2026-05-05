package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin   Role = "ROLE_ADMIN"
	RoleMember  Role = "ROLE_MEMBER"
	RoleVisitor Role = "ROLE_VISITOR"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleMember, RoleVisitor:
		return true
	}
	return false
}

type Membership struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_membership"`
	ClubID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_membership"`
	Role     Role      `gorm:"not null"`
	JoinedAt time.Time

	User     User     `gorm:"foreignKey:UserID"`
	BookClub BookClub `gorm:"foreignKey:ClubID"`
}

func NewMembership(userID, clubID uuid.UUID, role Role) *Membership {
	return &Membership{
		ID:       uuid.New(),
		UserID:   userID,
		ClubID:   clubID,
		Role:     role,
		JoinedAt: time.Now(),
	}
}

func (m *Membership) CanManageClub() bool  { return m.Role == RoleAdmin }
func (m *Membership) CanSuggestBook() bool { return m.Role == RoleAdmin || m.Role == RoleMember }
