package repository

type UserService interface {
	Get(login string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByPhoneNumber(phone_number string) (*User, error)
}
