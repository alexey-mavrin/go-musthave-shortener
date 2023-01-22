package app

type URL struct {
	URL string `json:"url"`
}

type Result struct {
	Result string `json:"result"`
}

type Config struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	sh              store
}
