package main

import (
	"log"

	"ball_scraper/analyze"
	"ball_scraper/scraper"
)

func main() {
	// Initialize the scraper service
	scraperService := scraper.NewScraperService()

	// Execute the scraping process
	if err := scraperService.ScrapeAndSaveResults(); err != nil {
		log.Fatalf("Scraping process failed: %v", err)
	}

	// Filepath to the CSV file
	filePath := "baloto_results.csv"

	// Create a new Analyzer instance
	analyzer := analyze.NewAnalyzer(filePath)

	// Call the analyze service
	analyzer.AnalyzeResults()

	log.Println("CSV file successfully generated.")
}
