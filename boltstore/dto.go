package boltstore

import (
	"strconv"
	"time"

	"github.com/oskanberg/dbo"
)

type DailyGrossDTO struct {
	IDDate string `storm:"id"`
	ID     int    `storm:"index"`
	Date   time.Time
	Gross  int
}

func DailyGrossToDTO(d dbo.DailyGross) DailyGrossDTO {
	idDate := strconv.Itoa(d.ID) + "-" + d.Date.String()
	return DailyGrossDTO{
		ID:     d.ID,
		Date:   d.Date,
		IDDate: idDate,
		Gross:  d.Gross,
	}
}

func DTOToDailyGross(d DailyGrossDTO) dbo.DailyGross {
	return dbo.DailyGross{
		ID:    d.ID,
		Date:  d.Date,
		Gross: d.Gross,
	}
}
