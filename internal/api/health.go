package api

import (
	"fmt"
	"net/http"
)

func (api *API) HealthHandler(res http.ResponseWriter, req *http.Request) {
	if !ensureMethod(res, req, http.MethodGet) {
		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := fmt.Fprint(res, "ok")

	if err != nil {
		api.log.Error(err.Error())
	}
}
