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
	vanguard.AddRoute("GET", "/killmailStatistics",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w, "killmailStatistics.html", time.Hour*24*31, newPage(r, "Killmail Statistics"))
		})

	vanguard.AddRoute("GET", "/J/killmailAttributes", killmailAttributesAPI)
	vanguard.AddRoute("GET", "/J/offensiveGroups", offensiveGroupsAPI)

	vanguard.AddRoute("GET", "/J/killmailStatistics", killmailStatisticsAPI)
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

	value, err := strconv.ParseInt(r.FormValue("value"), 10, 64)
	if err != nil {
		value = 0
	}

	points, err := strconv.ParseInt(r.FormValue("points"), 10, 64)
	if err != nil {
		points = 0
	}

	v, err := models.GetKillmailAttributes(groupID, value, points)
	if err != nil {
		httpErr(w, err)
		return
	}
	renderJSON(w, v, time.Hour*24)
}

func killmailStatisticsAPI(w http.ResponseWriter, r *http.Request) {
	v, err := models.GetKillmailStatistics()
	if err != nil {
		httpErr(w, err)
		return
	}
	renderJSON(w, v, time.Hour*24)
}
