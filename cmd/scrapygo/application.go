package scrapygo

import (
	"scrapygo/internal/config"
	"time"
)

type Application struct {
	*config.AppConfig
}

func (app *Application) ScrapeOnInterval(quit chan bool, seconds time.Duration) {
	go func() {
		for {
			select {
			case <-quit:
				app.Logger.Info("shutting down scraper...")
				return
			default:
				app.Logger.Info("scraping data...")
				app.Logger.Info("sleeping for 30 seconds")
				time.Sleep(seconds)
			}
		}
	}()

	return
}
