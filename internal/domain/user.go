package domain

type User struct {
	Id                     uint64
	Name                   string
	Email                  string
	Password               string
	EmailConfirmed         bool
	EmailConfirmationToken string
}
