package services

import (
	"errors"

	"bodyMaxIndex/repositories"
	"bodyMaxIndex/models"
)

// UserService handles CRUD operations of a user datamodel
type UserService interface {
	//GetAll()([]models.User)
	GetByID(id int64) (models.User, bool)
	GetByUsernameAndPassword(username, userPassword string) (models.User, bool)

	Update(id int64, user models.User) (models.User, error)
	UpdatePassword(id int64, newPassword string) (models.User, error)
	UpdateUsername(id int64, newUsername string) (models.User, error)

	Create(userPassword string, user models.User) (models.User, error)
}

// NewUserService returns the default user service.
func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

type userService struct {
	repo repositories.UserRepository
}

// GetByID returns a user based on its id.
func (s *userService) GetByID(id int64) (models.User, bool) {
	return s.repo.Select(func(m models.User) bool {
		return m.ID == id
	})
}

// GetByUsernameAndPassword returns a user based on its username and passowrd,
func (s *userService) GetByUsernameAndPassword(username, userPassword string) (models.User, bool) {
	if username == "" || userPassword == "" {
		return models.User{}, false
	}

	return s.repo.Select(func(m models.User) bool {
		if m.Username == username {
			hashed := m.HashedPassword
			if ok, _ := models.ValidatePassword(userPassword, hashed); ok {
				return true
			}
		}
		return false
	})
}

// Update updates every field from an existing User in 
// web/controllers/user_controller
func (s *userService) Update(id int64, user models.User) (models.User, error) {
	user.ID = id
	return s.repo.InsertOrUpdate(user)
}

// UpdatePassword updates a user's password.
func (s *userService) UpdatePassword(id int64, newPassword string) (models.User, error) {
	// update the user and return it.
	hashed, err := models.GeneratePassword(newPassword)
	if err != nil {
		return models.User{}, err
	}

	return s.Update(id, models.User{
		HashedPassword: hashed,
	})
}

// UpdateUsername updates a user's username.
func (s *userService) UpdateUsername(id int64, newUsername string) (models.User, error) {
	return s.Update(id, models.User{
		Username: newUsername,
	})
}

// Create inserts a new User
func (s *userService) Create(userPassword string, user models.User) (models.User, error) {
	if user.ID > 0 || userPassword == "" || user.Firstname == "" || user.Username == "" {
		return models.User{}, errors.New("unable to create this user")
	}

	hashed, err := models.GeneratePassword(userPassword)
	if err != nil {
		return models.User{}, err
	}
	user.HashedPassword = hashed

	return s.repo.InsertOrUpdate(user)
}

// GetAll ...
//func GetAll()(users []models.User){
//  return users;
//}