package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
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

func enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	dbFile := flag.String("db", "gross.db", "database file")
	serve := flag.Bool("serve", false, "serve the database via an API")
	fetchDate := flag.String("date", time.Now().Format("2006-01-02"), "date to load")
	lookupDays := flag.Int("days", 0, "number of days to fetch (looks back)")

	portInt := flag.Int("port", 8080, "serve port")
	portStr := strconv.Itoa(*portInt)

	flag.Parse()

	store, close, err := boltstore.New(*dbFile)
	if err != nil {
		log.Fatalf("error creating store, %s", err)
	}
	defer close()

	date, err := time.Parse("2006-01-02", *fetchDate)
	if err != nil {
		log.Println("failed to parse date", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(*lookupDays)
	for i := 0; i < *lookupDays; i++ {
		go func(offset int) {
			fetchAndStore(store, date.Add(time.Duration(-offset)*24*time.Hour).Truncate(24*time.Hour))
			wg.Done()
		}(i)
	}
	wg.Wait()

	if !*serve {
		os.Exit(0)
	}

	r := resolvers.NewRootResolver(store)

	gqlHandler := handler.GraphQL(graph.NewExecutableSchema(graph.Config{Resolvers: r}))
	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", enableCORS(gqlHandler))

	log.Printf("connect to http://localhost:%v/ for GraphQL playground", portStr)
	log.Fatal(http.ListenAndServe(":"+portStr, nil))

}
