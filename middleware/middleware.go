package middleware

import (
	"log"
	"net/http"
	"time"
)

// Middleware fonction faisant la passerelle entre la requête client et la logique
// métier de la requête
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Middlewares représente une collection de middleware
type Middlewares []Middleware

// Logger fonction tampon calculant le temps d'exécution de la transaction.
func Logger(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(begin time.Time) {
			log.Printf("[%-7s] %q %s", r.Method, r.URL.String(), time.Since(begin))
		}(time.Now())

		h.ServeHTTP(w, r)
	})
}
