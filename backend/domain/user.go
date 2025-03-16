package repository

type User struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"-"`
}
