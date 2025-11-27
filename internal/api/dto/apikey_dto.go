package dto

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type CreateAPIKeyRequest struct {
	Name string `json:"name"`
}

type APIKeyResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Prefix    string `json:"prefix"`
	CreatedAt string `json:"created_at"`
	RawKey    string `json:"raw_key,omitempty"` // Solo al crear
}
