package handle_request

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/36kone/gosend/internal/utils"
)

func HandleHttpProtocol(protocol string, url string) (string, error) {
	switch protocol {
	case "http":
		url = "http://" + url
	case "https":
		url = "https://" + url
	default:
		return "", fmt.Errorf("unsupported protocol: %s", protocol)
	}

	return url, nil
}

func HandleRequestWithHeaders(url string, method string, body any, contentType string, headers map[string]string) (string, int, error) {
	var reqBody io.Reader

	if body != nil {
		switch contentType {
		case "application/json":
			jsonBody, err := json.Marshal(body)
			if err != nil {
				return "", 0, err
			}
			reqBody = bytes.NewBuffer(jsonBody)
		case "application/x-www-form-urlencoded":
			reqBody = strings.NewReader(body.(string))
		}
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return "", 0, err
	}

	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	var response bytes.Buffer
	if err := json.Indent(&response, respBody, "", "  "); err != nil {
		response.Write(respBody)
	}

	return response.String(), resp.StatusCode, nil
}

func HandleRequest(url string, method string, body any, contentType string) (string, int, error) {
	var reqBody io.Reader

	if body != nil {
		switch contentType {
		case "application/json":
			jsonBody, err := json.Marshal(body)
			if err != nil {
				return "", 0, err
			}
			reqBody = bytes.NewBuffer(jsonBody)
		case "application/x-www-form-urlencoded":
			reqBody = strings.NewReader(body.(string))
		}
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return "", 0, err
	}

	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	var response bytes.Buffer
	if err := json.Indent(&response, respBody, "", "  "); err != nil {
		response.Write(respBody)
	}

	return response.String(), resp.StatusCode, nil
}

func HandleJsonBody(reader *bufio.Reader) (any, string) {
	var body any
	var contentType string

	fmt.Print("Send JSON or Form? (json/form): ")
	inputType, _ := reader.ReadString('\n')
	inputType = strings.TrimSpace(strings.ToLower(inputType))

	switch inputType {
	case "json":
		fmt.Print("Enter JSON Body: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		json.Unmarshal([]byte(input), &body)
		contentType = "application/json"
	case "form":
		fmt.Print("Enter Form Data (key=value&key2=value2): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		body = input
		contentType = "application/x-www-form-urlencoded"
	default:
		fmt.Println("Invalid input type:", inputType)
		return nil, ""
	}

	return body, contentType
}

func ReadURL(reader *bufio.Reader) (string, error) {
	fmt.Print("Enter URL: ")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)

	if url == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		fmt.Print("\nEnter Protocol (http/https) or press Enter to set https: ")
		protocol, _ := reader.ReadString('\n')
		protocol = strings.TrimSpace(protocol)

		if protocol == "" {
			protocol = "https"
		}

		reqUrl, err := HandleHttpProtocol(protocol, url)
		if err != nil {
			fmt.Println("Error:", err)
			return "", err
		}
		url = reqUrl
	}

	return url, nil
}

func HandleHttpMethod(reader *bufio.Reader) (string, any, string, error) {
	var method string
	var body any
	var contentType string

	fmt.Print("\nEnter HTTP Method (GET/POST/PUT/DELETE) or press Enter to set GET: ")
	method, _ = reader.ReadString('\n')
	method = strings.TrimSpace(strings.ToUpper(method))

	switch strings.ToUpper(method) {
	case "":
		method = "GET"
	case "GET":
		method = "GET"
	case "DELETE":
		method = "DELETE"
	case "POST", "PUT":
		body, contentType = HandleJsonBody(reader)
	default:
		fmt.Println("\nInvalid method:", method)
		return "", "", "", fmt.Errorf("invalid HTTP method: %s", method)
	}

	return method, body, contentType, nil
}

func ReadHeaders(reader *bufio.Reader) (map[string]string, bool) {
	headers := make(map[string]string)

	fmt.Println("\nDo you want to add HTTP headers? (y/n): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(strings.ToLower(choice))

	if choice == "n" {
		return headers, false
	}

	fmt.Println("Enter headers one per line in the format Key:Value")
	fmt.Println("Leave empty line to finish")

	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid header format. Use Key:Value")
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		headers[key] = value
	}

	return headers, true
}

func Run() {
	var response string
	var statusCode int

	reader := bufio.NewReader(os.Stdin)

	url, err := ReadURL(reader)

	if err != nil {
		fmt.Println(err)
		return
	}

	method, body, contentType, err := HandleHttpMethod(reader)

	if err != nil {
		fmt.Println(err)
		return
	}

	headers, hasHeaders := ReadHeaders(reader)

	if hasHeaders {
		response, statusCode, err = HandleRequestWithHeaders(url, method, body, contentType, headers)

		if err != nil {
			fmt.Println(err)
			return
		}

	} else {
		response, statusCode, err = HandleRequest(url, method, body, contentType)

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Response:", utils.ColorizeJSON(response))
	fmt.Println("\nStatus Code:", utils.ColorizeStatusCode(statusCode))
}
