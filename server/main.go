package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jasonlvhit/gocron"
	"github.com/nov1n/dumpertijdmachine/scraper"
	"github.com/nov1n/dumpertijdmachine/storage"
)

var (
	//dbName       = "/home/robert/dumpertijdmachine/storage.db"
	bucketName   = "videos-by-day"
	projectID    = "dumpertijdmachine"
	templatePath = "index.tmpl"
	startDate    = "2018.04.08"
	db           storage.Storage
	numVids      = 5
	cache        = make(map[time.Time]string)
)

func main() {
	// Initiate storage
	var err error
	db, err = storage.NewGCS(projectID, bucketName)
	if err != nil {
		log.Fatalf("could not create storage: %s", err)
	}
	defer db.Close()

	// Schedule crawler
	sc := scraper.NewScraper(db)
	gocron.Every(1).Day().At("00:05").Do(sc.Scrape)
	gocron.Every(2).Hours().Do(sc.Scrape)
	go func() {
		log.Print("Scheduling scraper every 2 hours.")
		<-gocron.Start()
	}()

	// Setup webserver
	r := mux.NewRouter()
	r.HandleFunc("/", dayHandler)
	r.HandleFunc(`/{date:\d\d\d\d\.\d\d\.\d\d}`, dayHandler)
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(":8000", loggedRouter))
}

func dayHandler(w http.ResponseWriter, r *http.Request) {
	dateS := mux.Vars(r)["date"]
	date, _ := time.Parse("2006.01.02", dateS)

	if r.URL.Path == "/" {
		date = time.Now()
	} else {
		// Check if it is cached
		if val, ok := cache[date]; ok {
			log.Printf("Response from cache for: %s", date)
			fmt.Fprintf(w, val)
			return
		}
	}
	day, err := db.GetDay(date)
	if err != nil {
		log.Printf("get day: %s", err)
		http.Error(w, "Hier is niks.", 500)
		return
	}

	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"prevDate": func(d string) string {
			t, _ := time.Parse("2006.01.02", d)
			return storage.KeyForDate(t.AddDate(0, 0, -1))
		},
		"nextDate": func(d string) string {
			t, _ := time.Parse("2006.01.02", d)
			return storage.KeyForDate(t.AddDate(0, 0, 1))
		},
		"hasNextDate": func(d string) bool {
			t, _ := time.Parse("2006.01.02", d)
			year1, month1, day1 := t.Date()
			year2, month2, day2 := time.Now().Date()
			return !(year1 == year2 && month1 == month2 && day1 == day2)
		},
		"hasPrevDate": func(d string) bool {
			t, _ := time.Parse("2006.01.02", d)
			t2, _ := time.Parse("2006.01.02", startDate)
			return t.After(t2)
		},
		"htmlDateFormat": func(d string) string {
			t, _ := time.Parse("2006.01.02", d)
			return t.Format("2006-01-02")
		},
		"base64": func(url string) string {
			name := storage.ImageName(url)
			img, err := db.GetImage(name)
			if err != nil {
				log.Printf("error getting image: %s", err)
				return `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`
			}
			return base64.StdEncoding.EncodeToString(img.Bytes())
		},
	}
	t := template.Must(template.New("index.tmpl").Funcs(funcMap).ParseFiles(templatePath))
	day.Videos = day.Videos[0:numVids]

	// Add entry to cache
	tpl := &bytes.Buffer{}
	err = t.Execute(tpl, day)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Printf("Add cache entry for: %s", date)
	cache[date] = tpl.String()
	fmt.Fprintf(w, tpl.String())
}
