package scraper

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/gocolly/colly"
	"github.com/hashicorp/go-multierror"
	"github.com/nov1n/dumpertijdmachine/storage"
	"github.com/nov1n/dumpertijdmachine/types"
)

type CollyScraper struct {
	db storage.Storage
}

func NewCollyScraper(db storage.Storage) *CollyScraper {
	return &CollyScraper{
		db: db,
	}
}

func (sc *CollyScraper) Scrape() error {
	c := colly.NewCollector()

	day := &types.Day{
		Date: storage.KeyForDate(time.Now()),
	}

	c.OnError(func(res *colly.Response, err error) {
		log.Printf("scraping error: %v (%v)", err, res)
	})

	c.OnHTML("body", func(h *colly.HTMLElement) {
		h.ForEach(".dump-cnt a.dumpthumb", func(i int, e *colly.HTMLElement) {
			views, kudos, err := getViewsAndKudos(e)
			if err != nil {
				log.Printf("getViewsAndKudos error: %v", err)
				return
			}

			v := &types.Video{
				Title: e.ChildText("h1"),
				Rank:  int32(i),
				Views: int32(views),
				Kudos: int32(kudos),
				Url:   e.Attr("href"),
				Thumb: e.ChildAttr("img", "src"),
			}

			day.Videos = append(day.Videos, v)

			err = sc.db.PutImage(e.ChildAttr("img", "src"))
			if err != nil {
				log.Printf("save image error: %s", err)
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting %s", r.URL)
	})
	c.Visit("https://www.dumpert.nl/toppers")

	// Save to storage
	todayKey := storage.KeyForDate(time.Now())
	err := sc.db.PutDay(todayKey, day)
	if err != nil {
		return fmt.Errorf("could not put day: %s", err)
	}
	return nil
}

func getViewsAndKudos(e *colly.HTMLElement) (views int, kudos int, result error) {
	re := regexp.MustCompile(`^views: (\d*) kudos: (\d*)$`)
	match := re.FindStringSubmatch(e.ChildText(".stats"))
	if len(match) < 3 {
		return 0, 0, fmt.Errorf("stats field did not match regex: %v, size was %d", re, len(match))
	}
	views, err := strconv.Atoi(match[1])
	if err != nil {
		result = multierror.Append(result, err)
	}
	kudos, err = strconv.Atoi(match[2])
	if err != nil {
		result = multierror.Append(result, err)
	}
	return
}
