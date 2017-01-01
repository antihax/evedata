package models

import (
	"net/http"
	"net/http/httputil"
)

func AddHTTPError(req *http.Request, res *http.Response) {
	reqText, _ := httputil.DumpRequest(req, true)
	resText, _ := httputil.DumpResponse(res, true)

	database.Exec(`INSERT  INTO httpErrors 
                    (url, status, request, response, time) 
                    VALUES(?,?,?,?,NOW());`,
		req.URL.String(), res.StatusCode, reqText, resText)
}
