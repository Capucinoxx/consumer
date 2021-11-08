package middleware

import "net/http"

// Middleware fonction faisant la passerelle entre la requête client et la logique
// métier de la requête
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Middlewares représente une collection de middleware
type Middlewares []Middleware
