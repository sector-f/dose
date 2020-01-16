package main

type AuthService interface {
	CheckAuthentication(username, password string) (bool, error)
}
