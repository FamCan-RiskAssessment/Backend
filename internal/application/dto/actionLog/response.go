package actionlogdto

import (
	"time"

	userdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/user"
)

type LogResponse struct {
	ID         uint                  `json:"id"`
	Actor      userdto.UserResponse  `json:"actor"`
	Target     *userdto.UserResponse `json:"target,omitempty"`
	Action     string                `json:"action"`
	Resource   *string               `json:"resource,omitempty"`
	ResourceID *uint                 `json:"resourceId,omitempty"`
	Details    string                `json:"details"`
	CreatedAt  time.Time             `json:"createdAt"`
	UpdatedAt  time.Time             `json:"updatedAt"`
}

type ActionType struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
