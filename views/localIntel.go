package views

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/antihax/evedata/eveConsumer"
	"github.com/antihax/goesi"

	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/server"
	"github.com/antihax/evedata/templates"

	"github.com/garyburd/redigo/redis"
)

func init() {
	evedata.AddRoute("localIntel", "GET", "/localIntel", localIntelPage)
	evedata.AddRoute("localIntel", "POST", "/J/localIntel", localIntel)
	evedata.AddRoute("localIntel", "GET", "/J/localIntel", localIntel)
}

func localIntelPage(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*60)
	p := newPage(r, "Local Intel Summary")
	hash := r.FormValue("hash")
	if hash != "" {
		p["HashURL"] = "/J/localIntel?hash=" + hash
	}

	templates.Templates = template.Must(template.ParseFiles("templates/localIntel.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)
	if err != nil {
		httpErr(w, err)
		return
	}
}

func localIntel(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*30)
	c := evedata.GlobalsFromContext(r.Context())

	hash := r.FormValue("hash")
	red := c.Cache.Get()
	defer red.Close()

	cache, err := redis.String(red.Do("GET", "EVEDATA_localIntel:"+hash))
	if err == nil {
		fmt.Fprintf(w, cache)
		return
	}

	type localdata struct {
		Local string `json:"local"`
	}
	var locl localdata
	if r.Body == nil {
		httpErrCode(w, http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&locl)
	if err != nil || len(locl.Local) == 0 {
		httpErrCode(w, http.StatusNotFound)
		return
	}

	names := strings.Split(locl.Local, "\n")
	newNames := removeDuplicatesAndValidate(names)

	// Get any one we don't know
	eveConsumer.CharSearchAddToQueue(newNames, &red)

	v, err := models.GetLocalIntel(newNames)
	if err != nil {
		httpErr(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(v)
	if err != nil {
		httpErr(w, err)
		return
	}

	w.Write(buf.Bytes())

	red.Do("SETEX", "EVEDATA_localIntel:"+hash, 86400, buf.Bytes())
}

// Remove any duplicate characters and delete any invalid entries
func removeDuplicatesAndValidate(xs []string) []interface{} {
	var n []interface{}
	found := make(map[string]bool)

	for _, x := range xs {
		x = strings.TrimSpace(x)
		if goesi.ValidCharacterName(x) {
			if !found[x] {
				found[x] = true
				n = append(n, x)
			}
		}
	}

	return n
}
