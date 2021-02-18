package health

import "net/http"

// Handler executes the health checks and returns http.StatusOK if they are successful.
func Handler(checker Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := checker.Check(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	}
}
