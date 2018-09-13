package boltstore

import (
	"time"

	"github.com/oskanberg/dbo/model"
)

type DailyGrossDTO struct {
	IDDate string `storm:"id"`
	ID     string `storm:"index"`
	Date   time.Time
	Gross  int
}

type FilmDTO struct {
	ID    string `storm:"id"`
	BOMID string `storm:"unique"`
	Title string
}

func FilmTODTO(d model.Film) FilmDTO {
	return FilmDTO{
		ID:    d.ID,
		BOMID: *d.BomID,
		Title: *d.Title,
	}
}

func DTOToFilm(d FilmDTO) model.Film {
	return model.Film{
		ID:    d.ID,
		BomID: &d.BOMID,
		Title: &d.Title,
	}
}

func DailyGrossToDTO(d model.DailyGross) DailyGrossDTO {
	idDate := d.ID + "-" + d.Date.String()
	return DailyGrossDTO{
		ID:     d.ID,
		Date:   d.Date,
		IDDate: idDate,
		Gross:  *d.Gross,
	}
}

func DTOToDailyGross(d DailyGrossDTO) model.DailyGross {
	return model.DailyGross{
		ID:    d.ID,
		Date:  d.Date,
		Gross: &d.Gross,
	}
}
