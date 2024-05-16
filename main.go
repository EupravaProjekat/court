package main

import (
	"context"
	"errors"
	"github.com/EupravaProjekat/court/Repo"
	"github.com/EupravaProjekat/court/handlers"
	habb "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	l := log.New(os.Stdout, "standard-api", log.LstdFlags)
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	repo, err := Repo.New(timeoutContext, l)
	if err != nil {
		l.Println(err)
	}
	defer func(accommodationRepo *Repo.Repo, ctx context.Context) {
		err := accommodationRepo.Disconnect(ctx)
		if err != nil {

		}
	}(repo, timeoutContext)

	// NoSQL: Checking if the connection was established
	repo.Ping()

	//Initialize the handler and inject said logger
	hh := handlers.NewBorderhendler(l, repo)

	router := mux.NewRouter()
	router.StrictSlash(true)
	//profile
	router.HandleFunc("/profile/{email}", hh.GetProfile).Methods("GET")
	router.HandleFunc("/newrequest", hh.NewRequest).Methods("POST")
	router.HandleFunc("/checkifuserexists", hh.CheckIfUserExists).Methods("GET")
	router.HandleFunc("/getrequest/{id}", hh.GetRequest).Methods("GET")
	router.HandleFunc("/adddata", hh.NewUser).Methods("POST")
	router.HandleFunc("/getallrequests", hh.GetallRequests).Methods("GET")
	router.HandleFunc("/getallcausings", hh.GetallCausings).Methods("GET")
	router.HandleFunc("/newcausing", hh.NewCausing).Methods("POST")

	headersOk := habb.AllowedHeaders([]string{"Content-Type", "jwt", "Authorization"})
	originsOk := habb.AllowedOrigins([]string{"http://localhost:4200"}) // Replace with your frontend origin
	methodsOk := habb.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Use the CORS middleware
	corsRouter := habb.CORS(originsOk, headersOk, methodsOk)(router)

	// Start the server
	srv := &http.Server{Addr: ":9198", Handler: corsRouter}
	go func() {
		log.Println("server starting border")
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}
	}()
	<-quit
	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
