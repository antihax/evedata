package views

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddRoute("GET", "/killmail", func(w http.ResponseWriter, r *http.Request) {

		idStr := r.FormValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpErr(w, err)
			return
		}

		km, err := models.GetKillmailDetails(id)
		if err != nil {
			httpErr(w, err)
			return
		}

		entity := "Someone"
		if km.AllianceID.Int64 > 0 {
			entity = km.AllianceName.String
		} else if km.CorporationID.Int64 > 0 {
			entity = km.CorporationName.String
		} else if km.CharacterID.Int64 > 0 {
			entity = km.CharacterName.String
		} else if km.FactionID.Int64 > 0 {
			entity = km.FactionName.String
		}

		fullEntity := ""
		if km.CharacterID.Int64 > 0 {
			fullEntity = km.CharacterName.String
		}
		if km.CorporationID.Int64 > 0 {
			if len(fullEntity) > 0 {
				fullEntity += " of "
			}
			fullEntity += km.CorporationName.String
		}
		if km.AllianceID.Int64 > 0 {
			fullEntity += " with " + km.AllianceName.String
		}
		if km.FactionID.Int64 > 0 {
			fullEntity += " fighting for the " + km.FactionName.String
		}

		title := fmt.Sprintf("%s lost their %s in %s (%.1f)",
			entity,
			km.TypeName,
			km.SolarSystemName,
			km.Security,
		)

		description := fmt.Sprintf("%s lost their %s in the %.1f security system of %s",
			fullEntity,
			km.TypeName,
			km.Security,
			km.SolarSystemName,
		)

		p := newPage(r, title)
		p["Killmail"] = km
		p["OG"] = OpenGraph{
			Image:       entityImage(int64(km.TypeID), "render", 128),
			Title:       title,
			Description: description,
		}
		renderTemplate(w, "killmail.html", time.Hour*24*31, p)
	})
}
