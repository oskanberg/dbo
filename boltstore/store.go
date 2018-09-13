package boltstore

import (
	"time"

	"github.com/asdine/storm"
	"github.com/oskanberg/dbo/model"
)

type Closer func() error

type Store struct {
	DB *storm.DB
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

func (s *Store) UpsertDetails(BOMID string, title string) (string, error) {
	d := s.DB.From("details")

	deets := FilmTODTO(model.Film{
		BomID: &BOMID,
		Title: &title,
	})

	err := d.Save(&deets)
	if err == storm.ErrAlreadyExists {
		detail := FilmDTO{}
		d.One("BOMID", BOMID, &detail)
		deets.ID = detail.ID
		err = d.Update(&deets)
	}

	return deets.ID, err
}

func (s *Store) AddDailyGross(id string, date time.Time, gross int) error {
	db := s.DB.From("gross").From("daily")

	g := DailyGrossToDTO(
		model.DailyGross{
			ID:    id,
			Date:  date,
			Gross: &gross,
		})
	err := db.Save(&g)

	return err
}

func (s *Store) GetDetails(ID string) (model.Film, error) {
	d := s.DB.From("details")

	deets := FilmDTO{}
	err := d.One("ID", ID, &deets)
	if err != nil {
		return model.Film{}, err
	}

	return DTOToFilm(deets), nil
}

func (s *Store) GetGross(ID int) ([]model.DailyGross, error) {
	db := s.DB.From("gross").From("daily")

	records := make([]DailyGrossDTO, 0)
	err := db.Find("ID", ID, &records)
	if err != nil {
		return nil, err
	}

	result := make([]model.DailyGross, len(records))
	for i, v := range records {
		result[i] = model.DailyGross{
			ID:    v.ID,
			Date:  v.Date,
			Gross: &v.Gross,
		}
	}

	return result, err
}
