package middleware

import "net/http"

func APIKeyMiddleware(validKey string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 🔥 IMPORTANT: Allow preflight
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}
		
		// Skip health endpoint
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		apiKey := r.Header.Get("X-API-Key")

		if apiKey == "" || apiKey != validKey {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"Unauthorized"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}