package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/altcha-org/altcha-lib-go"
)

func (api *API) CaptchaChallengeHandler(res http.ResponseWriter, req *http.Request) {
	if !ensureMethod(res, req, http.MethodGet) {
		return
	}

	challenge, err := api.captcha.CreateChallenge()

	if err != nil {
		api.log.Error(err.Error())
		writeError(res, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(res).Encode(challenge)

	if err != nil {
		api.log.Error(err.Error())
	}
}

func (api *API) CaptchaVerifyHandler(res http.ResponseWriter, req *http.Request) {
	if !ensureMethod(res, req, http.MethodPost) {
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		writeError(res, http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	decoded, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		writeError(res, http.StatusBadRequest)
		return
	}

	var payload altcha.Payload
	err = json.Unmarshal(decoded, &payload)
	if err != nil {
		writeError(res, http.StatusBadRequest)
		return
	}

	ok, err := api.captcha.VerifySolution(req.Context(), payload)

	if err != nil {
		api.log.Error(err.Error())
		writeError(res, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	if !ok {
		res.WriteHeader(http.StatusForbidden)
	}

	err = json.NewEncoder(res).Encode(map[string]bool{"success": ok})

	if err != nil {
		api.log.Error(err.Error())
	}
}
