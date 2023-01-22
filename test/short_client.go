package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/alexey-mavrin/go-musthave-shortener/internal/app"
)

func main() {
	var endpoint string

	var useJSON = flag.Bool("json", false, "use JSON request")
	var useGzip = flag.Bool("gzip", false, "use Gzip compression")
	flag.Parse()
	fmt.Println(*useJSON)

	fmt.Println("Введите длинный URL")
	reader := bufio.NewReader(os.Stdin)
	long, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		fmt.Println(err)
		os.Exit(1)
	}
	long = strings.TrimSuffix(long, "\n")
	client := &http.Client{}

	var r io.Reader

	if *useJSON {
		b, err := json.Marshal(app.URL{URL: long})
		if err != nil {
			log.Fatal(err)
		}
		var buf = bytes.NewBuffer(b)

		r = io.Reader(buf)
		endpoint = "http://localhost:8080/api/shorten"
	} else {
		r = strings.NewReader(long)
		endpoint = "http://localhost:8080/"
	}

	request, err := http.NewRequest(http.MethodPost, endpoint, r)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// в заголовках запроса сообщаем, что данные кодированы стандартной URL-схемой
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(long)))
	if *useGzip {
		request.Header.Add("Content-Encoding", "gzip")
	}
	// отправляем запрос и получаем ответ
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// печатаем код ответа
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	// читаем поток из тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// и печатаем его
	fmt.Println(string(body))
}
