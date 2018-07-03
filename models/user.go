package models

import(
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User ..
type User struct {
	ID				int64     	`json:"id" form:"id"`
	Firstname      	string    	`json:"firstname" form:"firstname"`
	Username       	string    	`json:"username" form:"username"`
	BMI				float64		`json:"bmi" form:"bmi"`
	Gender			string		`json:"gender" form:"gender"`
	Country			string		`json:"country" form:"country"`
	HashedPassword 	[]byte    	`json:"-" form:"-"`
	CreatedAt      	time.Time 	`json:"created_at" form:"created_at"`
}

// IsValid ...
func (u User) IsValid() bool {
	return u.ID > 0
}

// GeneratePassword generates a hashed password
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// ValidatePassword validates if passwords match
func ValidatePassword(userPassword string, hashed []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashed, []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}