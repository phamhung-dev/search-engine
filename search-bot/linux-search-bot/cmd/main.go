package main

import (
	"context"
	"linux-search-bot/business"
	"linux-search-bot/databases"
	"linux-search-bot/storage"
)

func main() {
	db := databases.NewPostgreSQLConnection()

	searchServiceStorage := storage.NewStorage(db)

	monitorBusiness := business.NewMonitorBusiness(searchServiceStorage)

	crawlBusiness := business.NewCrawlBusiness(searchServiceStorage)

	go func() {
		monitorBusiness.Monitor(context.Background())
	}()

	crawlBusiness.Crawl(context.Background())

	select {}
}
