package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var url string = "https://example.com/test?param=FUZZ"
var method string = "GET"
var threads int = 30
var headers = map[string]string{
	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
	"Accept-Language": "en-US",
}
var cookies = []*http.Cookie{}
var bodyParams = map[string]string{}
var payloads string
var statusCode, bodyLength int

var wg sync.WaitGroup

// Colors
var red = color.New(color.FgHiRed)
var cyan = color.New(color.FgCyan)
var blue = color.New(color.FgHiBlue)
var green = color.New(color.FgGreen)
var white = color.New(color.FgWhite)
var yellow = color.New(color.FgYellow)
var hiwhite = color.New(color.FgHiWhite)
var magenta = color.New(color.FgHiMagenta)

type anyType interface{}

func trim(s string) string {
	s = strings.TrimLeft(s, " ")
	return strings.TrimRight(s, " ")
}

func prettyPrint(v anyType) string {
	switch v := v.(type) {
	case map[string]string:
		s := ""
		for k, v := range v {
			s += fmt.Sprintf("%v: %v\n                       ", k, v)
		}
		return s
	case []*http.Cookie:
		s := ""
		for i, v := range v {
			s += fmt.Sprintf("%v: %v\n                       ", i, v)
		}
		return s
	default:
		return "unsupported type"
	}
}

func input() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanRunes)
	var s string
	for scanner.Scan() {
		r := scanner.Text()
		if r == "\n" {
			break
		}
		s += r
	}
	return strings.ReplaceAll(s, "\r", "")
}

func contains(s1 string, s2 string) bool {
	s1 = strings.ToLower(s1)
	return strings.Contains(s1, s2)
}

func setMethod(m string) {
	switch strings.ToLower(strings.TrimSpace(m)) {
	case "get":
		method = "GET"
	case "post":
		method = "POST"
	default:
		fmt.Println("Method may be 'GET' or 'POST'.")
	}
}

func setHeaders() {
	fmt.Print("Add header in following format name:value -> ")
	header := input()
	parts := strings.SplitN(header, ":", 2)
	if len(parts) == 2 {
		headers[parts[0]] = parts[1]
	} else {
		fmt.Println("Wrong input.")
	}
}

func setCookies() {
	for {
		fmt.Print("Add Cookie in following format name:value or press enter to continue: ")
		cookie := input()
		parts := strings.SplitN(cookie, ":", 2)
		if cookie == "" {
			break
		} else if len(parts) == 2 {
			cookies = append(cookies, &http.Cookie{Name: parts[0], Value: parts[1]})
		} else {
			fmt.Println("Wrong input.")
		}
	}
}

func setBodyParams() {
	for {
		fmt.Print("Add data parameter in following format param=value or press enter to continue: ")
		data := input()
		parts := strings.SplitN(data, "=", 2)
		if data == "" {
			break
		} else if len(parts) == 2 {
			bodyParams[parts[0]] = parts[1]
		} else {
			fmt.Println("Wrong input.")
		}
	}
}

func urlEncode(data map[string]string) string {
	var result []string
	for key, value := range data {
		result = append(result, key+"="+value)
	}
	return strings.Join(result, "&")
}

func makeRequest(url string, method string, headers map[string]string, cookies []*http.Cookie, bodyParams map[string]string) (int, int, string) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	var req *http.Request
	var err error

	if method == "POST" {
		data := urlEncode(bodyParams)
		req, err = http.NewRequest("POST", url, strings.NewReader(data))
		if err != nil {
			return 0, 0, ""
		}
		if _, ok := headers["Content-Type"]; !ok {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	} else {
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
			return 0, 0, ""
		}
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	response, err := client.Do(req)
	if err != nil {
		return 0, 0, ""
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, 0, ""
	}

	return response.StatusCode, len(body), string(body)
}

func getRange(s string) (int, int, int) {
	re := regexp.MustCompile(`^(-?[0-9]+):(-?[0-9]+):(-?[0-9]+)$`)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return 0, 0, 0
	}
	from, _ := strconv.Atoi(matches[1])
	to, _ := strconv.Atoi(matches[2])
	step, _ := strconv.Atoi(matches[3])
	return from, to, step
}

func printFuzzBanner() {
	fmt.Println("--------------------------")
	fmt.Println("Code  Length       Payload")
	fmt.Println("--------------------------")
}

func fuzz() anyType {
	buffer := make(chan struct{}, threads)
	from, to, step := getRange(payloads)
	laps := 0
	if from == 0 && to == 0 {
		file, err := os.Open(payloads)
		if err != nil {
			fmt.Println("Payloads not set")
			return ""
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		printFuzzBanner()
		for scanner.Scan() {
			laps++
			line := scanner.Text()
			yellow.Printf("-[ATT::%v]\r", laps)
			wg.Add(1)
			buffer <- struct{}{}
			newGoroutine(line, buffer)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
	} else {
		yellow.Printf("\nIterate from %v to %v with step of %v\n\n", from, to, step)
		printFuzzBanner()
		for {
			laps++
			num := strconv.Itoa(from)
			yellow.Printf("-[ATT::%v]\r", laps)
			wg.Add(1)
			buffer <- struct{}{}
			newGoroutine(num, buffer)

			// from, to, step maybe negative numbers
			if from == to {
				break
			} else {
				from += step
			}
		}
	}
	wg.Wait()
	white.Printf("Total payloads sent: %v\n", laps)
	white.Print("Press enter to exit ")
	input()
	return ""
}

func printFinding(statusCode int, bodyLength int, payload string) {
	if statusCode >= 200 && statusCode < 300 {
		output := fmt.Sprintf("%-5v L:%-10d %v", statusCode, bodyLength, payload)
		green.Println(output)
	} else if statusCode >= 300 && statusCode < 400 {
		output := fmt.Sprintf("%-5v L:%-10d %v", statusCode, bodyLength, payload)
		cyan.Println(output)
	} else if statusCode == 403 {
		output := fmt.Sprintf("%-5v L:%-10d %v", statusCode, bodyLength, payload)
		red.Println(output)
	}
}

func newGoroutine(payload string, buffer chan struct{}) {
	go func(payload string) {
		defer wg.Done()
		defer func() { <-buffer }()
		url := strings.ReplaceAll(url, "FUZZ", payload)
		if method == "POST" {
			statusCode, bodyLength, _ = makeRequest(url, "POST", headers, cookies, bodyParams)
		} else {
			statusCode, bodyLength, _ = makeRequest(url, "GET", headers, cookies, nil)
		}
		printFinding(statusCode, statusCode, payload)

	}(payload)
}

func printState() {
	options := fmt.Sprintf(`
Method                 %v
URL                    %v
Payloads               %v
Threads                %v
Headers                %v
Cookies                %v
Body Parameters        %v

`, method, url, payloads, threads, prettyPrint(headers), prettyPrint(cookies), prettyPrint(bodyParams))
	hiwhite.Println(options)
}

func cmd() {
	magenta.Printf(">>>  ")
}

func main() {
	blue.Printf(`
____ _    ____ ___  ____ ____ 
|__, |___ |--|  /__ |__, [__]  v1.0.0

All-in-One fuzzer by @leddcode

To use GlazGo, you will need to provide the URL of the website or application that you want to test.
The URL should contain the string "FUZZ" where you want the tool to inject test data.

Keep in mind that fuzzing can generate a large number of requests and may potentially cause issues with the website or application being tested.
It is important to use caution and obtain permission before fuzzing any production systems.`)
	printState()
	for {
		var action string
		cmd()
		action = input()
		if contains(action, "run") {
			break
		} else if contains(action, "options") {
			printState()
		} else if contains(action, "set method") {
			fmt.Print("     HTTP Method: ")
			setMethod(trim(input()))
		} else if contains(action, "set url") {
			fmt.Print("     Target URL: ")
			url = trim(input())
		} else if contains(action, "set threads") {
			fmt.Print("     Number of threads: ")
			i, err := strconv.Atoi(trim(input()))
			if err != nil {
				fmt.Println(err)
				return
			}
			threads = i
		} else if contains(action, "set payloads") {
			fmt.Print("     Payloads: ")
			payloads = trim(input())
		} else if contains(action, "set header") {
			setHeaders()
		} else if contains(action, "set cookie") {
			setCookies()
		} else if contains(action, "set data") {
			setBodyParams()
		} else {
			if action != "" {
				red.Println(action, "??")
				hiwhite.Println(`Commands:
	- options       (Show current settings)
	- run           (Runs GlazGo)
	- set method    (GET or POST - default GET)
	- set url       (Target URL containig FUZZ placeholder)
	- set data      (param=value)
	- set header    (header:value)
	- set cookie    (cookie:value)
	- set threads   (Single integer - default 30)
	- set payloads  (The path to the file containing the payloads
	                 or the range of numbers in format from_number:to_number:step - ex. 1:100:1)`)
			}

			action = ""
		}
	}

	fuzz()

}
