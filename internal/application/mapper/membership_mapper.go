package mapper

import (
	"reader-club/internal/application/dto"
	"reader-club/internal/domain/entity"
)

func ToMembershipResponse(m *entity.Membership) dto.MembershipResponse {
	return dto.MembershipResponse{
		ID:       m.ID,
		User:     ToUserResponse(&m.User),
		Role:     m.Role,
		JoinedAt: m.JoinedAt,
	}
}
