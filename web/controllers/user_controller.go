package controllers

import(

	"fmt"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"

	"bodyMaxIndex/services"
	"bodyMaxIndex/models"
	
	"strconv"
)

// UserController ..
type UserController struct{
	Ctx iris.Context
	Service services.UserService
	Session *sessions.Session
}

const userIDKey = "UserID"

// getCurrentUserID ...
// TODO: create BaseCOntroller with these 3 methods by default and create methods to implement them
func (c *UserController) getCurrentUserID() int64 {
	userID := c.Session.GetInt64Default(userIDKey, 0)
	return userID
}

// isLoggedIn ... 
func (c *UserController) isLoggedIn() bool {
	return c.getCurrentUserID() > 0
}

func (c *UserController) logout() {
	c.Session.Destroy()
}

// TODO: move statics views to separeta folder
var registerStaticView = mvc.View{
	Name: "user/register.html",
	Data: iris.Map{"Title": "User Registration"},
}

// GetRegister handles GET: http://localhost:8080/user/register.
func (c *UserController) GetRegister() mvc.Result {
	if c.isLoggedIn() {
		c.logout()
	}

	return registerStaticView
}

// PostRegister handles POST: http://localhost:8080/user/register.
func (c *UserController) PostRegister() mvc.Result {
	var (
		firstname = c.Ctx.FormValue("firstname")
		username  = c.Ctx.FormValue("username")
		password  = c.Ctx.FormValue("password")
		height,_ = strconv.ParseFloat(c.Ctx.FormValue("height"),32)
		mass,_ = strconv.ParseFloat(c.Ctx.FormValue("mass"),32)
		calcType = c.Ctx.FormValue("type")
		BMI = 0.0
	)

	// if type is for m-kg
	if(calcType == "1"){
		BMI = mass/(height*height)
	}else{
		BMI = (mass/(height*height))*703
	}

	// create the new user, the password will be hashed by the service.
	u, err := c.Service.Create(password, models.User{
		Username:  username,
		Firstname: firstname,
		BMI: BMI,
	})

	// set the user's id to this session even if err != nil
	c.Session.Set(userIDKey, u.ID)

	return mvc.Response{
		Err: err,
		Path: "/user/me",
	}

}

var loginStaticView = mvc.View{
	Name: "user/login.html",
	Data: iris.Map{"Title": "User Login"},
}

// GetLogin handles GET: http://localhost:8080/user/login.
func (c *UserController) GetLogin() mvc.Result {
	// verify if users is already logged in
	if c.isLoggedIn() {
		c.logout()
	}

	return loginStaticView
}

// PostLogin handles POST: http://localhost:8080/user/register.
func (c *UserController) PostLogin() mvc.Result {
	var (
		username = c.Ctx.FormValue("username")
		password = c.Ctx.FormValue("password")
	)

	u, found := c.Service.GetByUsernameAndPassword(username, password)

	if !found {
		return mvc.Response{
			Path: "/user/register",
		}
	}

	c.Session.Set(userIDKey, u.ID)

	return mvc.Response{
		Path: "/user/me",
	}
}

// GetMe handles GET: http://localhost:8080/user/me.
func (c *UserController) GetMe() mvc.Result {
	if !c.isLoggedIn() {
		return mvc.Response{Path: "/user/login"}
	}

	u, found := c.Service.GetByID(c.getCurrentUserID())
	if !found {
		c.logout()
		return c.GetMe()
	}

	cat := getUserCategory(u);

	return mvc.View{
		Name: "user/me.html",
		Data: iris.Map{
			"Title": "Profile of " + u.Username,
			"User":  u,
			"BMI": fmt.Sprintf("%f", u.BMI),
			"Category": cat,
		},
	}
}

// AnyLogout handles All/Any HTTP Methods for: http://localhost:8080/user/logout.
func (c *UserController) AnyLogout() {
	if c.isLoggedIn() {
		c.logout()
	}
	c.Ctx.Redirect("/user/login")
}