package model

// User is an authentication entry for identifying the users of this service.
type User struct {
	ID       int
	Nickname string
	Password []byte
}
