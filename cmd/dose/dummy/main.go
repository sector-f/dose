package dummy

type NoAuth struct{}

func (a NoAuth) AuthRequired() bool {
	return false
}

func (a NoAuth) CheckAuthentication(username, password string) (bool, error) {
	return true, nil
}

type EmptyAuth struct{}

func (a EmptyAuth) AuthRequired() bool {
	return true
}

func (a EmptyAuth) CheckAuthentication(username, password string) (bool, error) {
	return true, nil
}
