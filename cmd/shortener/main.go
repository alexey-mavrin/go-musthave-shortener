package main

import (
	"log"
	"os"

	"github.com/alexey-mavrin/go-musthave-shortener/internal/app"
)

func main() {
	c := app.Config{
		ServerAddress: ":8080",
		BaseURL:       "http://localhost:8080/",
	}

	if sa, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		c.ServerAddress = sa
	}

	if bu, ok := os.LookupEnv("BASE_URL"); ok {
		c.BaseURL = bu
	}

	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
}
