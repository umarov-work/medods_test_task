package dto

type UpdateTokensRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
