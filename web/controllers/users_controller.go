package controllers

import (
	"bodyMaxIndex/models"
	"bodyMaxIndex/services"

	"github.com/kataras/iris"
)

// UsersController is our /users API controller
// Requires basic authentication
type UsersController struct {
	Ctx iris.Context

	// Our UserService, it's an interface which
	Service services.UserService
}

// Get returns list of the users
//func (c *UsersController) Get() (results []models.User) {
//	return c.Service.GetAll()
//}

// GetBy returns an user ---[GET] http://localhost:8080/users/{{id}}
func (c *UsersController) GetBy(id int64) (user models.User, found bool) {
	u, found := c.Service.GetByID(id)
	if !found {
		c.Ctx.Values().Set("message", "User nor be found!")
	}
	return u, found
}

// PutBy updates a user ---[PUT] http://localhost:8080/users/1
func (c *UsersController) PutBy(id int64) (models.User, error) {
	
	// Read form paaameteres in iris context
	u := models.User{}
	if err := c.Ctx.ReadForm(&u); err != nil {
		return u, err
	}

	return c.Service.Update(id, u)
}

