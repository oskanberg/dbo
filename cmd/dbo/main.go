package main

import (
	"flag"
	"log"
	"time"

	"github.com/oskanberg/dbo/boltstore"

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
	wordPtr := flag.String("db", "gross.db", "database file")
	flag.Parse()

	store, close, err := boltstore.New(*wordPtr)
	if err != nil {
		log.Fatalf("error creating store, %s", err)
	}
	defer close()

	for i := 0; i < 100; i++ {
		date := time.Now().
			Add(time.Duration(-i) * day).
			Truncate(24 * time.Hour)
		fetchAndStore(store, date)
	}

	store.ListDailyGross()
}
