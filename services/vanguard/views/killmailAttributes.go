package views

import (
	"net/http"
	"strconv"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddRoute("GET", "/killmailAttributes",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w, "killmailAttributes.html", time.Hour*24*31, newPage(r, "Killmail Attribute Browser"))
		})

	vanguard.AddRoute("GET", "/J/killmailAttributes", killmailAttributesAPI)
	vanguard.AddRoute("GET", "/J/offensiveGroups", offensiveGroupsAPI)
}

func offensiveGroupsAPI(w http.ResponseWriter, r *http.Request) {
	v, err := models.GetOffensiveShipGroupID()
	if err != nil {
		httpErr(w, err)
		return
	}
	renderJSON(w, v, time.Hour*24*31)
}

func killmailAttributesAPI(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.ParseInt(r.FormValue("groupID"), 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	v, err := models.GetKillmailAttributes(groupID)
	if err != nil {
		httpErr(w, err)
		return
	}
	renderJSON(w, v, time.Hour*24)
}
