package config

const (
	DBConnection     = "postgres://postgres:135274@localhost/tasks_db?sslmode=disable"
	DBConnectionType = "postgres"
	HTTPPort         = "80"
	PrivateKey       = `-----BEGIN RSA PRIVATE KEY-----
MIIBOQIBAAJAZSOF9iuJYCxwbUSwFoNteC+Z0rifXvvhJK5NghtimuJmD5xfwySL
CwXhraKfXEUtz+T6XXA2Rp1tY+pVq+FHwQIDAQABAkAPwNi830smj8VzP5+t4grL
DZ8IE3m/cbw/2mZ4PYu+VBJ1YjzKecM/HSqq4mvQH8KgGQ02x/3f2MgJ+5Eadw+B
AiEAvuQGRduC8BeuIfiIHvcOw9rUJaN2DyHmHKYC+4q8zO8CIQCHoqQXLa+0zzg0
q/KcBJ8SfMxIsWuXTEX2yqdGURaWTwIgaZj+l1ptPp/65jP0KR0Gf/Xn8cJRJuHb
x/FWKQyAkOUCIQCDy9Rq+Wfc5+aTt+mdFRiFXGMc19nWQLVTZAQ63Zx3HQIgCQzr
TRtfs3Ax22LNLz2blixObO7FwoV1oC/5ovr6FJE=
-----END RSA PRIVATE KEY-----`

	PublicKey = `-----BEGIN PUBLIC KEY-----
MFswDQYJKoZIhvcNAQEBBQADSgAwRwJAZSOF9iuJYCxwbUSwFoNteC+Z0rifXvvh
JK5NghtimuJmD5xfwySLCwXhraKfXEUtz+T6XXA2Rp1tY+pVq+FHwQIDAQAB
-----END PUBLIC KEY-----`

	MasterPassword = "135274"
)
