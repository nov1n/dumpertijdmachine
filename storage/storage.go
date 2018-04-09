package storage

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/nov1n/dumpertijdmachine/types"
)

type Storage interface {
	PutDay(date time.Time, day *types.Day) error
	GetDay(date time.Time) (*types.Day, error)
	PutImage(url string) error
	GetImage(url string) (*bytes.Buffer, error)
	Close() error
}

func KeyForDate(date time.Time) string {
	return fmt.Sprintf("%d.%02d.%02d", date.Year(), date.Month(), date.Day())
}

func ImageName(url string) string {
	parts := strings.Split(url, "/")
	return "pics/" + parts[len(parts)-1]
}
