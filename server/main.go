package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

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
	db           storage.Storage
	numVids      = 5
	cache        = make(map[string]string)
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
	gocron.Every(1).Day().At("23:55").Do(sc.Scrape)
	go func() {
		log.Print("Scheduling scraper every day at 23:55")
		<-gocron.Start()
	}()

	// Setup webserver
	r := mux.NewRouter()
	r.HandleFunc("/", dayHandler)
	r.HandleFunc(`/{date:\d\d\d\d/\d\d/\d\d}`, dayHandler)

	log.Fatal(http.ListenAndServe(":8000", r))
}
func dayHandler(w http.ResponseWriter, r *http.Request) {
	dateKey := mux.Vars(r)["date"]
	if r.URL.Path == "/" {
		dateKey = storage.KeyForDate(time.Now())
	} else {
		// Check if it is cached
		if val, ok := cache[string(dateKey)]; ok {
			log.Printf("Response from cache for: %s", dateKey)
			fmt.Fprintf(w, val)
			return
		}
	}
	day, err := db.GetDay(dateKey)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"prevDate": func(d string) string {
			t, _ := time.Parse("2006/01/02", d)
			return string(storage.KeyForDate(t.AddDate(0, 0, -1)))
		},
		"nextDate": func(d string) string {
			t, _ := time.Parse("2006/01/02", d)
			return string(storage.KeyForDate(t.AddDate(0, 0, 1)))
		},
		"hasNextDate": func(d string) bool {
			t, _ := time.Parse("2006/01/02", d)
			year1, month1, day1 := t.Date()
			year2, month2, day2 := time.Now().Date()
			return !(year1 == year2 && month1 == month2 && day1 == day2)
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
	log.Printf("Add cache entry for: %s", dateKey)
	cache[string(dateKey)] = tpl.String()
	fmt.Fprintf(w, tpl.String())
}
