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

// MeassurementController ...
type MeassurementController struct {
	Ctx iris.Context
	Service services.UserService
	Session *sessions.Session
}

// Category ...
type Category struct{
	From	float64
	To		float64
	Name	string
}

// getCategories ...
func getCategories(countryID int)[]Category{
	return []Category{
		Category{Name: "Very severely underweight", From: 0, To: 15},
		Category{Name: "Severely underweight", From: 15, To: 16},
		Category{Name: "Underweight", From: 16, To: 18.5},
		Category{Name: "Normal (healthy weight)", From: 18.5, To: 25},
		Category{Name: "Overweight", From: 25, To: 30},
		Category{Name: "Obese Class I (Moderately obese)", From: 30, To: 35},
		Category{Name: "Obese Class II (Severely obese)", From: 35, To: 40},
		Category{Name: "Obese Class III (Very severely obese)", From: 40, To: 45},
		Category{Name: "Obese Class IV (Morbidly Obese)", From: 45, To: 50},
		Category{Name: "Obese Class V (Super Obese)", From: 50, To: 55},
		Category{Name: "Obese Class VI (Hyper Obese)", From: 55, To: 150},
	}

}

// getUserCategory ...
func getUserCategory(user models.User)Category{
	var catResult Category;
   	for _,c:= range getCategories(1){
		if(c.From<user.BMI && c.To>user.BMI){
			catResult = c;
			break;
		}
   	}
   	return catResult;
}

func (c *MeassurementController) getCurrentUserID() int64 {
	userID := c.Session.GetInt64Default(userIDKey, 0)
	return userID
}

func (c *MeassurementController) isLoggedIn() bool {
	return c.getCurrentUserID() > 0
}

func (c *MeassurementController) logout() {
	c.Session.Destroy()
}

// BeforeActivation ...
func (c *MeassurementController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/bmi", "PostCalculate")
}

// AfterActivation ...
func (c *MeassurementController) AfterActivation(a mvc.AfterActivation) {
	if a.Singleton() {
		panic("MeassurementController should be stateless, a request-scoped, we have a 'Session' which depends on the context.")
	}
}


// PostCalculate handles POST: http://localhost:8080/user/calculate
func (c *MeassurementController) PostCalculate()mvc.Result{
	if !c.isLoggedIn() {
		// if it's not logged in then redirect user to the login page.
		return mvc.Response{Path: "/user/login"}
	}
	u, found := c.Service.GetByID(c.getCurrentUserID())
	if !found {
		c.logout()
		return mvc.Response{
			Path: "/user/register",
		}
	}

	var (
		height,_ = strconv.ParseFloat(c.Ctx.FormValue("height"),32)
		mass,_ = strconv.ParseFloat(c.Ctx.FormValue("mass"),32)
		calcType = c.Ctx.FormValue("type")
	)

	// if type is for m-kg
	if(calcType == "1"){
		u.BMI = mass/(height*height)
	}else{
		u.BMI = (mass/(height*height))*703
	}

	uu, err := c.Service.Update(c.getCurrentUserID(),u)
	if(err != nil){
	}

	cat := getUserCategory(uu);

	return mvc.View{
		Name: "user/me.html",
		Data: iris.Map{
			"Title": "Profile of " + uu.Username,
			"User":  uu,
			"BMI": fmt.Sprintf("%f", u.BMI),
			"Category": cat,
		},
	}

}
