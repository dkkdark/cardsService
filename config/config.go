package config

const (
	DBConnection     = "postgres://postgres:135274@localhost/tasks_db?sslmode=disable"
	DBConnectionType = "postgres"
	HTTPPort         = "80"
	PrivateKey       = `-----BEGIN RSA PRIVATE KEY-----
	## KEY ##
-----END RSA PRIVATE KEY-----`

	PublicKey = `-----BEGIN PUBLIC KEY-----
MFswDQYJKoZIhvcNAQEBBQADSgAwRwJAZSOF9iuJYCxwbUSwFoNteC+Z0rifXvvh
JK5NghtimuJmD5xfwySLCwXhraKfXEUtz+T6XXA2Rp1tY+pVq+FHwQIDAQAB
-----END PUBLIC KEY-----`

	MasterPassword = "135274"
)
