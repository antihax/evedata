package views

import (
	"encoding/json"
	"evedata/models"
	"evedata/server"
	"evedata/templates"
	"html/template"
	"net/http"
	"strconv"
)

func init() {
	evedata.AddRoute(evedata.Route{"login", "GET", "/login", loginPage})
	evedata.AddRoute(evedata.Route{"login", "POST", "/login", loginPostPage})
}

func loginPage(c *evedata.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	p := Page{
		Title: "Login to your account",
	}

	templates.Templates = template.Must(template.ParseFiles("templates/user/login.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func loginPostPage(c *evedata.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	u, err := models.AuthenticateUser(r.FormValue("username"), r.FormValue("password"), c.Db)

	encoder := json.NewEncoder(w)

	if err != nil {
		encoder.Encode(struct{ Status int }{0})
		return http.StatusOK, nil
	}

	r.AddCookie(&http.Cookie{
		Name:  "uid",
		Value: strconv.FormatInt(u.UID, 10),
		Path:  "/",
	})

	r.AddCookie(&http.Cookie{
		Name:  "pass",
		Value: u.Password,
		Path:  "/",
	})

	encoder.Encode(struct{ Status int }{1})

	return http.StatusOK, nil
}
