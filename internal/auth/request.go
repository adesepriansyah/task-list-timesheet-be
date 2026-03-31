package auth

// loginRequest represents the request body for login.
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
