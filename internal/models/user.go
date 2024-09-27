package models

type User struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Role     string `db:"role"`
}

func NewUser(username, password, role string) *User {
	return &User{
		Username: username,
		Password: password,
		Role:     role,
	}
}
