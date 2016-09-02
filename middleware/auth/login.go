package auth

type Login struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Email    string `json:"Email"`
}
