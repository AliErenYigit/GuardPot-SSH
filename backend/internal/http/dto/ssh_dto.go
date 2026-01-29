package dto

type CreateSSHConnectionRequest struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`

	AuthType string `json:"authType"` // "password" | "private_key"
	Secret   string `json:"secret"`   // password OR private key pem (plain, will be encrypted at rest)
}
