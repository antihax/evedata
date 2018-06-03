package views

import (
	"net/http"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddRoute("mutaplasmidEst", "GET", "/mutaplasmidEst",
		func(w http.ResponseWriter, r *http.Request) {
			mpt := r.FormValue("type")
			if mpt == "" {
				mpt = "Warp Disruptor"
			}

			v, err := models.GetMutaplasmidData(mpt)
			if err != nil {
				httpErr(w, err)
				return
			}
			p := newPage(r, "Mutaplasmid Estimator")

			// Get the mutaplasmid types
			types := make([]string, len(models.MutaplasmidTypes))
			i := 0
			for k := range models.MutaplasmidTypes {
				types[i] = k
				i++
			}

			p["Data"] = "[" + v.Data + "]"
			p["Graphs"] = "[" + v.MetaData + "]"
			p["Types"] = types

			renderTemplate(w,
				"mutaplasmidEst.html",
				time.Hour*24*31,
				p)
		})
}
