package main

import (
	"github.com/boltdb/bolt"
	"github.com/3DPrintShop/PrintQL/printdb"
	"log"
	"net/http"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/smithaitufe/go-graphql-upload"

	"github.com/rs/cors"
	"github.com/3DPrintShop/PrintQL/handler"
	"github.com/3DPrintShop/PrintQL/loader"
	"github.com/3DPrintShop/PrintQL/resolver"
	"github.com/3DPrintShop/PrintQL/schema"
)

func main() {

	var (
		addr              = ":8000"
		readHeaderTimeout = 1 * time.Second
		writeTimeout      = 10 * time.Second
		idleTimeout       = 90 * time.Second
		maxHeaderBytes    = http.DefaultMaxHeaderBytes
	)

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	root, err := resolver.NewRoot()
	if err != nil {
		log.Fatal(err)
	}

	db, err := bolt.Open("my.db", 0600, nil)

	printClient, err := printdb.NewClient(db)
	h := handler.GraphQL{
		Schema:  graphql.MustParseSchema(schema.String(), root),
		Loaders: loader.Initialize(printClient),
		Client:  printClient,
		DB:      db,
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	upload := http.FileServer(http.Dir("./uploads"))
	media := http.FileServer(http.Dir("./media"))
	mux.Handle("/", fs)
	mux.Handle("/assets/", http.StripPrefix("/assets/", upload))
	mux.Handle("/media/", http.StripPrefix("/media/", media))
	mux.Handle("/gql/", handler.GraphiQL{})
	mux.Handle("/gql", handler.GraphiQL{})
	mux.Handle("/graphql/", graphqlupload.Handler(handler.AddUserContext(h)))
	mux.Handle("/graphql", graphqlupload.Handler(handler.AddUserContext(h)))
	mux.HandleFunc("/auth/google/login", handler.OauthGoogleLogin)
	mux.HandleFunc("/auth/google/callback", handler.OauthGoogleCallback)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(mux)

	s := &http.Server{
		Addr:              addr,
		Handler:           corsHandler,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	log.Printf("Listening for requests on %s", s.Addr)

	if err = s.ListenAndServe(); err != nil {
		log.Println("server.ListenAndServe:", err)
	}
}
