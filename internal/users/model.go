package users

import "time"

type UserModel struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
}
