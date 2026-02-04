package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/36kone/gosend/internal/handle_request"
	"github.com/36kone/gosend/internal/utils"
)

func main() {
	urlFlag := flag.String("u", "", "Request URL (e.g. https://example.com/api)")
	methodFlag := flag.String("X", "GET", "HTTP Method (GET, POST, PUT, DELETE)")
	dataFlag := flag.String("d", "", "Request body (JSON or form, raw string)")
	headersFlag := flag.String("H", "", "Headers (comma-separated, ex: Authorization:Bearer token,Content-Type:application/json)")
	noColorFlag := flag.Bool("n", false, "Disable colored output")
	prettyFlag := flag.Bool("p", true, "Pretty-print JSON")
	flag.Parse()

	utils.PrintBanner()

	var url string
	if *urlFlag != "" {
		url = *urlFlag
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		var err error
		url, err = handle_request.ReadURL(reader)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	method := strings.ToUpper(*methodFlag)
	var body any
	var contentType string

	if method == "POST" || method == "PUT" {
		if *dataFlag != "" {
			if strings.HasPrefix(*dataFlag, "{") && strings.HasSuffix(*dataFlag, "}") {
				json.Unmarshal([]byte(*dataFlag), &body)
				contentType = "application/json"
			} else {
				body = *dataFlag
				contentType = "application/x-www-form-urlencoded"
			}
		} else {
			reader := bufio.NewReader(os.Stdin)
			body, contentType = handle_request.HandleJsonBody(reader)
		}
	}

	headers := make(map[string]string)
	if *headersFlag != "" {
		parts := strings.Split(*headersFlag, ",")
		for _, h := range parts {
			hPair := strings.SplitN(h, ":", 2)
			if len(hPair) == 2 {
				headers[strings.TrimSpace(hPair[0])] = strings.TrimSpace(hPair[1])
			}
		}
	}

	response, status, err := handle_request.HandleRequestWithHeaders(url, method, body, contentType, headers)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *noColorFlag {
		fmt.Println("\nResponse:", response)
		fmt.Println("\nStatus Code:", status)
	} else {
		if *prettyFlag {
			fmt.Println("\nResponse:", utils.ColorizeJSON(response))
		} else {
			fmt.Println("\nResponse:", response)
		}
		fmt.Println("\nStatus Code:", utils.ColorizeStatusCode(status))
	}
}
