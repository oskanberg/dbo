package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/oskanberg/dbo/boltstore"
	"github.com/oskanberg/dbo/graph"
	"github.com/oskanberg/dbo/graph/resolvers"
	"github.com/vektah/gqlgen/handler"

	"github.com/oskanberg/dbo"
)

const (
	day  = time.Hour * 24
	year = day * 365
)

func fetchAndStore(db *boltstore.Store, date time.Time) {
	gross, err := dbo.GetDailyGross(date)
	if err != nil {
		log.Printf("failed to get daily gross: %v", err)
		return
	}

	for _, record := range gross {
		id, err := db.UpsertDetails(
			record.BOMID,
			record.Title,
		)
		if err != nil {
			log.Printf("Unable to store details: %v", err)
		}

		err = db.AddDailyGross(id, date, record.Gross)
		if err != nil {
			log.Printf("failed to add daily gross: %v", err)
		}
	}

}

func main() {
	dbFile := flag.String("db", "gross.db", "database file")
	serve := flag.Bool("serve", false, "serve the database via an API")
	portInt := flag.Int("port", 8080, "serve port")
	portStr := strconv.Itoa(*portInt)

	flag.Parse()

	store, close, err := boltstore.New(*dbFile)
	if err != nil {
		log.Fatalf("error creating store, %s", err)
	}
	defer close()

	log.Println("fetching today's gross")

	date := time.Now().Truncate(24 * time.Hour)
	fetchAndStore(store, date)

	if !*serve {
		os.Exit(0)
	}

	r := resolvers.NewRootResolver(store)

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(graph.NewExecutableSchema(graph.Config{Resolvers: r})))

	log.Printf("connect to http://localhost:%v/ for GraphQL playground", portStr)
	log.Fatal(http.ListenAndServe(":"+portStr, nil))

}
