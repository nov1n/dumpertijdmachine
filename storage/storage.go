package storage

import (
	"fmt"
	"time"

	"github.com/nov1n/dumpertijdmachine/types"
)

type Storage interface {
	PutDay(date time.Time, day *types.Day) error
	GetDay(date time.Time) (*types.Day, error)
	Close() error
}

func KeyForDate(date time.Time) string {
	return fmt.Sprintf("%d.%02d.%02d", date.Year(), date.Month(), date.Day())
}
