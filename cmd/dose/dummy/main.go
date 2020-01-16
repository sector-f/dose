package dummy

// Auth implements dose.AuthService
// It always returns true
type Auth struct{}

func (a Auth) CheckAuthentication(username, password string) (bool, error) {
	return true, nil
}
