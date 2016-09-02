package auth

type AuthRequest struct {
	Username    string
	Email       string
	AppToken    string
	AccessLevel int8
	Path        string
}
