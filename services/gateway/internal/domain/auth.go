package domain

const (
	OwnerKindGuest = "guest"
	OwnerKindUser  = "user"
)

type RegisterRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	CreatedAt   string `json:"created_at,omitempty"`
}

type UserSession struct {
	SessionToken string `json:"session_token"`
	User         User   `json:"user"`
	CreatedAt    string `json:"created_at,omitempty"`
	LastSeenAt   string `json:"last_seen_at,omitempty"`
}
