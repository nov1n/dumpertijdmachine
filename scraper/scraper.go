package scraper

var (
	dbName     = "/home/robert/dumpertijdmachine/storage.db"
	bucketName = "videos-by-day"
)

type Scraper interface {
	Scrape() error
}
