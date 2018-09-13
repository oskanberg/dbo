package boltstore

import (
	"fmt"
	"log"
	"time"

	"github.com/asdine/storm"
	"github.com/oskanberg/dbo"
)

type Closer func() error

type Store struct {
	DB *storm.DB
}

type Details struct {
	ID    int    `storm:"id,increment"`
	BOMID string `storm:"unique"`
	Title string
}

func New(dbLocation string) (*Store, Closer, error) {
	db, err := storm.Open(dbLocation)
	if err != nil {
		return nil, nil, err
	}

	d := db.From("gross").From("daily")
	d.Init(DailyGrossDTO{})

	return &Store{db}, db.Close, nil
}

func (s *Store) UpsertDetails(BOMID string, title string) (int, error) {
	d := s.DB.From("details")

	deets := Details{
		BOMID: BOMID,
		Title: title,
	}

	err := d.Save(&deets)
	if err == storm.ErrAlreadyExists {
		detail := Details{}
		d.One("BOMID", BOMID, &detail)
		deets.ID = detail.ID
		err = d.Update(&deets)
	}

	return deets.ID, err
}

func (s *Store) AddDailyGross(id int, date time.Time, gross int) error {
	db := s.DB.From("gross").From("daily")

	g := DailyGrossToDTO(
		dbo.DailyGross{
			ID:    id,
			Date:  date,
			Gross: gross,
		})
	err := db.Save(&g)

	return err
}

func (s *Store) GetDetails(ID int) (dbo.Details, error) {
	d := s.DB.From("details")

	deets := Details{}
	err := d.One("ID", ID, &deets)
	if err != nil {
		return dbo.Details{}, err
	}

	return dbo.Details{
		ID:    deets.ID,
		BOMID: deets.BOMID,
		Title: deets.Title,
	}, nil
}

func (s *Store) GetGross(ID int) ([]dbo.DailyGross, error) {
	db := s.DB.From("gross").From("daily")

	records := make([]DailyGrossDTO, 0)
	err := db.Find("ID", ID, &records)
	if err != nil {
		return nil, err
	}

	result := make([]dbo.DailyGross, len(records))
	for i, v := range records {
		result[i] = dbo.DailyGross{
			ID:    v.ID,
			Date:  v.Date,
			Gross: v.Gross,
		}
	}

	return result, err
}

func (s *Store) ListDailyGross() {
	var records []Details

	db := s.DB.From("details")
	err := db.All(&records)
	if err != nil {
		panic(err)
	}

	for _, r := range records {
		gross, err := s.GetGross(r.ID)
		if err != nil {
			log.Println(r.ID, err)
			fmt.Printf("%v: no gross\n", r.Title)
			continue
		}
		fmt.Printf("%v:%v\n", r.Title, gross)
	}
}
