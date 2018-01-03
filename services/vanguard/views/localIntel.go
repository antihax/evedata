package views

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
	"github.com/antihax/evedata/services/vanguard/templates"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
)

func init() {
	vanguard.AddRoute("localIntel", "GET", "/localIntel", localIntelPage)
	vanguard.AddRoute("localIntel", "POST", "/J/localIntel", localIntel)
	vanguard.AddRoute("localIntel", "GET", "/J/localIntel", localIntel)
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
	c := vanguard.GlobalsFromContext(r.Context())

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
		httpErrCode(w, err, http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&locl)
	if err != nil || len(locl.Local) == 0 {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	names := strings.Split(locl.Local, "\n")

	work := []redisqueue.Work{}
	// Add any characters we do not know to the list
	newNames := removeDuplicatesAndValidate(names)
	for _, name := range newNames {
		work = append(work, redisqueue.Work{Operation: "charSearch", Parameter: name})
	}
	c.OutQueue.QueueWork(work, redisqueue.Priority_Urgent)

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
