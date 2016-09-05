package auth

type AuthRequest struct {
	Username    string `json:"Username"`
	Email       string `json:"Email"`
	AppToken    string `json:"AppToken"`
	AccessLevel int    `json:"AccessLevel"`
	Path        string `json:"Path"`
}
