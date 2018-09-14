package boltstore

import (
	"time"

	"github.com/asdine/storm/q"

	"github.com/asdine/storm"
	"github.com/google/uuid"
	"github.com/oskanberg/dbo/graph/model"
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

	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	deets := FilmTODTO(model.Film{
		ID:    uuid.String(),
		BomID: &BOMID,
		Title: &title,
	})

	err = d.Save(&deets)
	if err == storm.ErrAlreadyExists {
		detail := FilmDTO{}
		d.One("BOMID", BOMID, &detail)
		deets.ID = detail.ID
		err = d.Update(&deets)
	}
	if err != nil {
		return "", err
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

// GetGross gets the gross between two ranges
func (s *Store) GetGross(ID string, from, to time.Time) ([]model.DailyGross, error) {
	db := s.DB.From("gross").From("daily")

	records := make([]DailyGrossDTO, 0)
	query := db.Select(
		q.And(
			q.Eq("ID", ID),
			q.Lte("Date", to),
			q.Gte("Date", from),
		),
	)
	err := query.Find(&records)

	if err != nil {
		return nil, err
	}

	result := make([]model.DailyGross, len(records))
	for i, v := range records {
		g := v.Gross
		result[i] = model.DailyGross{
			ID:    v.ID,
			Date:  v.Date,
			Gross: &g,
		}
	}

	return result, err
}
