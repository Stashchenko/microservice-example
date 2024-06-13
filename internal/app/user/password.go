package user

type Password string

func NewPassword(p string) (Password, error) {
	return Password(p), nil
}

func (p Password) String() string {
	return string(p)
}
