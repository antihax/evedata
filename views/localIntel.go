package views

import (
	"bytes"
	"encoding/json"
	"errors"
	"evedata/appContext"
	"evedata/eveapi"
	"evedata/models"
	"evedata/server"
	"evedata/templates"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddRoute("localIntel", "GET", "/localIntel", localIntelPage)
	evedata.AddRoute("localIntel", "POST", "/J/localIntel", localIntel)
	evedata.AddRoute("localIntel", "GET", "/J/localIntel", localIntel)
}

func localIntelPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 60*60)
	p := newPage(s, r, "Local Intel Summary")
	hash := r.FormValue("hash")
	if hash != "" {
		p["HashURL"] = "/J/localIntel?hash=" + hash
	}

	templates.Templates = template.Must(template.ParseFiles("templates/localIntel.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func localIntel(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 60*30)
	hash := r.FormValue("hash")
	red := c.Cache.Get()
	defer red.Close()

	cache, err := redis.String(red.Do("GET", "EVEDATA_localIntel:"+hash))
	if err == nil {
		fmt.Fprintf(w, cache)
		return http.StatusOK, nil
	}

	type localdata struct {
		Local string `json:"local"`
	}
	var locl localdata
	if r.Body == nil {
		return http.StatusNotFound, errors.New("No Data Received")
	}
	err = json.NewDecoder(r.Body).Decode(&locl)
	if err != nil || len(locl.Local) == 0 {
		return http.StatusNotFound, err
	}

	names := strings.Split(locl.Local, "\n")
	newNames := removeDuplicatesAndValidate(names)

	v, err := models.GetLocalIntel(newNames)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(v)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Write(buf.Bytes())

	red.Do("SETEX", "EVEDATA_localIntel:"+hash, 86400, buf.Bytes())
	return http.StatusOK, nil
}

// Remove any duplicate characters and delete any invalid entries
func removeDuplicatesAndValidate(xs []string) []interface{} {
	found := make(map[string]bool)
	n := make([]interface{}, 0)

	for _, x := range xs {
		if eveapi.ValidCharacterName(x) {
			if !found[x] {
				found[x] = true
				n = append(n, x)
			}
		}
	}
	return n
}
