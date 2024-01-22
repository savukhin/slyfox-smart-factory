package domain

import "strings"

type User struct {
	Username string `db:"username"`
	TotpKey  string `db:"totp_key"`
	AesKey   string `db:"aes_key"`
	Active   bool   `db:"active"`
}

func GenerateUser() User {
	return User{
		Username: "Some",
		TotpKey:  "Some secret",
		AesKey:   strings.Repeat("a", 32),
		Active:   true,
	}
}
