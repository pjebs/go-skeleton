package services

import (
	"net/http"
)

func ServiceProviders(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(rw, r)
}
