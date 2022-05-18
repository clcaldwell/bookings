package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/clcaldwell/bookings/pkg/config"
	"github.com/clcaldwell/bookings/pkg/handlers"
	"github.com/clcaldwell/bookings/pkg/render"

	"github.com/alexedwards/scs/v2"
)

var portNumber string = ":8080"
var app config.AppConfig
var session *scs.SessionManager

func main() {

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	print := fmt.Sprintf("Starting application on port %s", portNumber)
	fmt.Println(print)

	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
