package scraper

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/schollz/progressbar/v3"
	"os"
	"strings"
	"sync"
	"time"
)

// ScraperService defines the service for scraping data
type ScraperService struct {
	collector *colly.Collector
	results   [][]string
}

// NewScraperService creates a new ScraperService instance
func NewScraperService() *ScraperService {
	c := colly.NewCollector(
		colly.AllowedDomains("baloto.com"),
	)
	return &ScraperService{
		collector: c,
		results:   make([][]string, 0),
	}
}

// ScrapeAndSaveResults performs the scraping and saves the results to a CSV file
func (s *ScraperService) ScrapeAndSaveResults() error {
	var resultsMutex sync.Mutex

	// Callback to handle data extraction from each page
	s.collector.OnHTML("#results-table", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			if el.ChildText("thead.text-uppercase.text-center") != "" {
				// Skip this record
				return
			}

			var draw, date, result, superBall string

			if el.ChildAttr("img[src='https://d314ivgy8nq27r.cloudfront.net/static/img/baloto-kind.png']", "src") != "" {
				draw = "baloto"
			} else if el.ChildAttr("img[src='https://d314ivgy8nq27r.cloudfront.net/static/img/revancha-kind.png']", "src") != "" {
				draw = "revancha"
			}

			// Get and transform the date
			dateStr := el.ChildText("td.creation-date-results:not([style='font-weight: bold'])")
			if dateStr == "" {
				// Skip this record if the date is blank or null
				return
			}
			date, err := transformDate(dateStr)
			if err != nil {
				fmt.Println("Error transforming the date:", err)
				return
			}

			// Get the result
			el.ForEach("td[style='font-weight: bold'].creation-date-results", func(_ int, elem *colly.HTMLElement) {
				result = elem.Text
				// Remove spaces between dashes and the last dash
				result = strings.ReplaceAll(result, " - ", "-")
				if strings.HasSuffix(result, "-") {
					result = result[:len(result)-1]
				}
				// Take only the first 5 numbers
				numbers := strings.Split(result, "-")
				if len(numbers) > 5 {
					result = strings.Join(numbers[:5], "-")
				}
			})

			// Get the super ball
			superBall = el.ChildText("span.balota-red-results")

			// Format the rows correctly
			record := []string{
				draw,
				date,
				result,
				superBall,
			}
			resultsMutex.Lock()
			s.results = append(s.results, record)
			resultsMutex.Unlock()
		})
	})

	// Number of pages to visit
	totalPages := 71

	// Create a progress bar
	bar := progressbar.NewOptions(totalPages,
		progressbar.OptionSetDescription("Crawling pages"),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "#", SaucerHead: ">", SaucerPadding: "-", BarStart: "[", BarEnd: "]"}),
	)

	// Visit each page
	for i := 1; i <= totalPages; i++ {
		url := fmt.Sprintf("https://baloto.com/resultados?page=%d", i)
		s.collector.Visit(url)
		bar.Add(1)
	}

	// Create the CSV file
	file, err := os.Create("baloto_results.csv")
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	writer.Write([]string{"Draw", "Date", "Result", "Super Ball"})

	// Write results
	for _, result := range s.results {
		writer.Write(result)
	}

	return nil
}

// transformDate transforms the date to the format YYYY-MM-DD
func transformDate(dateStr string) (string, error) {
	// Create a map of months in Spanish to months in English
	months := map[string]string{
		"Enero":      "January",
		"Febrero":    "February",
		"Marzo":      "March",
		"Abril":      "April",
		"Mayo":       "May",
		"Junio":      "June",
		"Julio":      "July",
		"Agosto":     "August",
		"Septiembre": "September",
		"Octubre":    "October",
		"Noviembre":  "November",
		"Diciembre":  "December",
	}

	// Split the date string
	parts := strings.Split(dateStr, " ")
	if len(parts) != 5 {
		return "", fmt.Errorf("invalid date format: %s", dateStr)
	}
	day := parts[0]
	month := months[parts[2]]
	year := parts[4]

	// Create a date string in English format
	dateInStr := fmt.Sprintf("%s %s %s", day, month, year)

	// Parse the date string to the time.Time format
	date, err := time.Parse("2 January 2006", dateInStr)
	if err != nil {
		return "", err
	}

	// Return the date in the format YYYY-MM-DD
	return date.Format("2006-01-02"), nil
}
