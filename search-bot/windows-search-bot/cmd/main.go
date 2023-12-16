package main

import (
	"context"
	"windows-search-bot/business"
	"windows-search-bot/databases"
	"windows-search-bot/storage"
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
