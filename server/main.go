package main

import (
	"bytes"
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
			fmt.Println(year1, year2, month1, month2, day1, day2)
			return !(year1 == year2 && month1 == month2 && day1 == day2)
		},
		"hasPrevDate": func(d string) bool {
			t, _ := time.Parse("2006.01.02", d)
			t2, _ := time.Parse("2006.01.02", startDate)
			return t.After(t2)
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
