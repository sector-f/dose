package main

type AuthService interface {
	AuthRequired() bool
	CheckAuthentication(username, password string) (bool, error)
}
