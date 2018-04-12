package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/nov1n/dumpertijdmachine/storage"
	"github.com/nov1n/dumpertijdmachine/types"
)

type topper struct {
	Title       string `json:"title"`
	Url         string `json:"url"`
	Thumb       string `json:"thumb"`
	Mediatype   string `json:"mediatype"`
	Date        string `json:"date"`
	Views       string `json:"views"`
	Kudos       string `json:"kudos"`
	Description string `json:"description"`
}

type top5 struct {
	Top5 []topper `json:"top5"`
}
type StandaloneScraper struct {
	db storage.Storage
}

func NewStandaloneScraper(db storage.Storage) *StandaloneScraper {
	return &StandaloneScraper{
		db: db,
	}
}

func (sc *StandaloneScraper) Scrape() error {
	day := &types.Day{
		Date: storage.KeyForDate(time.Now()),
	}

	url := "https://www.dumpert.nl/js/snippets.js.php"
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("getting toppers: %s", err)
	}
	defer res.Body.Close()

	toppers := &top5{}
	err = json.NewDecoder(res.Body).Decode(toppers)
	if err != nil {
		return fmt.Errorf("can't deserialize toppers: %v", err)
	}

	for i, t := range toppers.Top5 {
		views, err := strconv.ParseInt(t.Views, 10, 32)
		if err != nil {
			return err
		}
		kudos, err := strconv.ParseInt(t.Kudos, 10, 32)
		if err != nil {
			return err
		}

		v := &types.Video{
			Title: t.Title,
			Rank:  int32(i),
			Views: int32(views),
			Kudos: int32(kudos),
			Url:   t.Url,
			Thumb: t.Thumb,
		}

		err = sc.db.PutImage(t.Thumb)
		if err != nil {
			return fmt.Errorf("save image error: %s", err)
		}

		day.Videos = append(day.Videos, v)
	}

	// Save to storage
	todayKey := storage.KeyForDate(time.Now())
	err = sc.db.PutDay(todayKey, day)
	if err != nil {
		return fmt.Errorf("could not put day: %s", err)
	}
	return nil
}
