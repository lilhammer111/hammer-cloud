package middleware

import "net/http"

func TokenAuthMDW(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "wrong parameter", http.StatusBadRequest)
			return
		}
		username := r.Form.Get("username")
		if len(username) < 3 || !IsTokenValid(token) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		// h(w, r) is also right here
		h.ServeHTTP(w, r)
	}
}

func IsTokenValid(token string) bool {
	// TODO to judge if token expires, and validate token
	if len(token) != 40 {
		return false
	}
	return true
}
