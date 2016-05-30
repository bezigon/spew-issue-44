package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/davecgh/go-spew/spew"
	"github.com/goincremental/negroni-oauth2"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"golang.org/x/oauth2/vk"
)

func main() {

	secureMux := http.NewServeMux()

	// Routes that require a logged in user
	// can be protected by using a separate route handler
	// If the user is not authenticated, they will be
	// redirected to the login path.
	secureMux.HandleFunc("/restrict", func(w http.ResponseWriter, req *http.Request) {
		token := oauth2.GetToken(req)
		fmt.Fprintf(w, "OK: %s", token.Access())
	})

	secure := negroni.New()
	secure.Use(oauth2.LoginRequired())
	secure.UseHandler(secureMux)

	n := negroni.New()
	n.Use(sessions.Sessions("my_session", cookiestore.New([]byte("secret123"))))

	config := &oauth2.Config{
		ClientID:     "5485578",
		ClientSecret: "JEISAVYEWh9zmhEB5b5V",
		RedirectURL:  "http://localhost:5000/oauth2callback",
		Scopes:       []string{"friends", "audio", "video", "offline"},
	}

	provider := oauth2.NewOAuth2Provider(config, vk.Endpoint.AuthURL, vk.Endpoint.TokenURL)
	n.Use(provider)

	router := http.NewServeMux()

	//routes added to mux do not require authentication
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		token := oauth2.GetToken(req)
		if token == nil || !token.Valid() {
			fmt.Fprintf(w, "not logged in, or the access token is expired")
			return
		}
		fmt.Fprintf(w, "logged in")

		spew.Dump(token)
	})

	//There is probably a nicer way to handle this than repeat the restricted routes again
	//of course, you could use something like gorilla/mux and define prefix / regex etc.
	router.Handle("/restrict", secure)

	n.UseHandler(router)

	n.Run(":5000")
}
