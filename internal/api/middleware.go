package api

import "net/http"

type Middleware func(http.Handler) http.Handler

func chain(h http.Handler, m ...Middleware) http.Handler {
	// loop in reverse to preserve middleware order
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}

func noSniffMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(res, req)
	})
}
