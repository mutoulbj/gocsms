package dto

import (
	"time"

	"github.com/mutoulbj/gocsms/internal/models"
)

type OrganizationResponse struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Slug           string   `json:"slug"`
	Description    string   `json:"description,omitempty"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
	ChargeStations []string `json:"charge_stations,omitempty"` // List of Charge Station IDs
}

type OrganizationCreateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Slug        string `json:"slug" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"omitempty,max=500"`
}

func ToOrganizationResponse(o *models.Organization) *OrganizationResponse {
	if o == nil {
		return nil
	}
	response := &OrganizationResponse{
		ID:          o.ID.String(),
		Name:        o.Name,
		Slug:        o.Slug,
		Description: o.Description,
		CreatedAt:   o.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   o.UpdatedAt.Format(time.RFC3339),
	}

	if len(o.ChargeStations) > 0 {
		response.ChargeStations = make([]string, len(o.ChargeStations))
		for i, cs := range o.ChargeStations {
			response.ChargeStations[i] = cs.ID.String()
		}
	}

	return response
}
