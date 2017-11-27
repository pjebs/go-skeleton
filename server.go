package main

import (
	"log"
	"net/http"

	"github.com/pjebs/restgate"
	"github.com/urfave/negroni"

	//Framework
	route "./app/http"
	s "./app/providers"
	c "./config"
)

func main() {

	app := negroni.New()

	app.Use(negroni.NewRecovery())
	app.Use(negroni.NewLogger())
	//Add Security Credential Requirements to API routes.
	//For production, keep HTTPSProtection = true
	HTTPSProtection := false
	if HTTPSProtection {
		app.Use(restgate.New("X-Auth-Key", "X-Auth-Secret", restgate.Static, restgate.Config{HTTPSProtectionOff: false, Key: []string{c.API_ENDPOINT_KEY}, Secret: []string{c.API_ENDPOINT_SECRET}}))
	}
	app.Use(negroni.HandlerFunc(s.ServiceProviders)) //Initializes Service Providers
	app.UseHandler(route.New())
	http.Handle("/", app)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
