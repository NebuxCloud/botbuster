package api

import "net/http"

func ensureMethod(res http.ResponseWriter, req *http.Request, method string) bool {
	if req.Method != method {
		res.Header().Set("Allow", method)
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	return true
}

func writeError(res http.ResponseWriter, code int) {
	http.Error(
		res,
		http.StatusText(code),
		code,
	)
}
