package auth

type contextKey string

const (
	ContextKeyUserID   contextKey = "user_id"
	ContextKeyUsername contextKey = "username"
	ContextKeyRoles    contextKey = "roles"
)
