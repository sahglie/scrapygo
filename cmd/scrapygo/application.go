package scrapygo

import (
	"context"
	"scrapygo/internal/config"
	"scrapygo/internal/database"
	"scrapygo/internal/services"
	"time"
)

type Application struct {
	*config.AppConfig
}

func (app *Application) ScrapeOnInterval(quit chan bool, seconds time.Duration) {
	service := services.NewConfig(app.DB, app.Logger)

	ticker := time.NewTicker(seconds)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			app.Logger.Info("shutting down scraper...")
			return
		case <-ticker.C:
			app.Logger.Info("scraping feeds...")
			feeds, err := app.DB.GetNextFeedsToScrape(context.TODO())
			if err != nil {
				app.Logger.Error("failed to get next feeds", "err", err)
				continue
			}

			for _, f := range feeds {
				go func(feed database.Feed) {
					app.Logger.Info("scraping feed: ", "f", f)
					scrapeErr := service.ScrapeFeed(f, nil)
					if scrapeErr == nil {
						app.DB.MarkFeedFetched(context.TODO(), f.ID)
					}
				}(f)
			}
		}
	}
}
