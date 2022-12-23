package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/alexey-mavrin/go-musthave-shortener/internal/app"
)

func parseConfig() app.Config {
	c := app.Config{
		ServerAddress: ":8080",
		BaseURL:       "http://localhost:8080",
	}

	saFlag := flag.String("a", ":8080", "server address")
	buFlag := flag.String("b", "http://localhost:8080", "base url")
	stFlag := flag.String("f", "", "storage file")
	flag.Parse()

	if saEnv, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		c.ServerAddress = saEnv
	}
	if saFlag != nil {
		c.ServerAddress = *saFlag
	}

	if buEnv, ok := os.LookupEnv("BASE_URL"); ok {
		c.BaseURL = buEnv
	}
	if buFlag != nil {
		c.BaseURL = *buFlag
	}

	if stEnv, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		c.FileStoragePath = stEnv
	}
	if stFlag != nil {
		c.FileStoragePath = *stFlag
	}

	conf, _ := json.Marshal(c)
	log.Printf("server arguments: %s", string(conf))
	return c
}

func main() {
	c := parseConfig()

	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
}
