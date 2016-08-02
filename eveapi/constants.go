package eveapi

import (
	"strings"
	"time"
)

type eveURI struct {
	AppManagement string
	CREST         string
	Images        string
	Login         string
	XML           string
}

var eveTQ = eveURI{
	AppManagement: "https://developers.eveonline.com/",
	CREST:         "https://crest-tq.eveonline.com/",
	Images:        "https://image.eveonline.com/",
	Login:         "https://login.eveonline.com/",
	XML:           "https://api.eveonline.com/",
}

var eveSisi = eveURI{
	AppManagement: "https://developers.testeveonline.com/",
	CREST:         "https://api-sisi.testeveonline.com/",
	Images:        "https://image.testeveonline.com/",
	Login:         "https://sisilogin.testeveonline.com/",
	XML:           "https://api.testeveonline.com/",
}

type EVETime struct {
	time.Time
}

// Cannot properly Unmarshal CCP's time stamps?
const eveTimeLayout = "2006-01-02T15:04:05"

func (c *EVETime) UnmarshalJSON(b []byte) (err error) {
	t := string(b)
	t = strings.Replace(t, `"`, "", -1)
	c.Time, err = time.Parse(eveTimeLayout, t)
	return
}
