package main

import (
	"fight-alerts-backend/datastore"
	handler "fight-alerts-backend/lambda_handler"
	"fight-alerts-backend/scraper"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	host := os.Getenv("RDS_HOST")
	password := os.Getenv("RDS_PASSWORD")
	username := os.Getenv("RDS_USERNAME")
	d := &datastore.Datastore{Host: host, Password: password, User: username, Dbname: "FightAlertsDb", Port: 5432}

	s := scraper.Scraper{Url: scraper.MmaUrl}

	lh := handler.Handler{Scraper: s, Datastore: d}

	lambda.Start(lh.HandleRequest)
}
