package http_test

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	http "github.com/zMrKrabz/fhttp"
	"github.com/zMrKrabz/fhttp/httptrace"
)

// Basic http test with Header Order
func TestExample(t *testing.T) {
	c := http.Client{}
	req, err := http.NewRequest("GET", "https://httpbin.org/headers", strings.NewReader(""))
	req.Header = http.Header{
		"sec-ch-ua":                 {"\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\""},
		"sec-ch-ua-mobile":          {"?0"},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36"},
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"sec-fetch-site":            {"none"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-user":            {"?1"},
		"sec-fetch-dest":            {"document"},
		"accept-encoding":           {"gzip, deflate, br"},
		http.HeaderOrderKey: {
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"upgrade-insecure-requests",
			"user-agent",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-user",
			"sec-fetch-dest",
			"accept-encoding",
		},
	}

	if err != nil {
		t.Errorf(err.Error())
	}

	resp, err := c.Do(req)

	if err != nil {
		t.Errorf(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %v", resp.StatusCode)
	}

	var data interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		t.Errorf(err.Error())
	}
}

// Test with Charles cert + proxy
func TestWithCert(t *testing.T) {
	home, err := os.UserHomeDir()

	if err != nil {
		t.Errorf(err.Error())
	}

	caCert, err := os.ReadFile(fmt.Sprintf("%v/charles_cert.pem", home))

	if err != nil {
		t.Errorf("Could not find charles cert, %v", err.Error())
		return
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	proxyURL, err := url.Parse("http://localhost:8888")

	if err != nil {
		t.Errorf(err.Error())
	}

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
			Proxy:             http.ProxyURL(proxyURL),
			ForceAttemptHTTP2: true,
		},
	}

	req, err := http.NewRequest("GET", "https://httpbin.org/headers", strings.NewReader(""))

	if err != nil {
		t.Errorf(err.Error())
	}

	req.Header = http.Header{
		"sec-ch-ua":                 {"\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\""},
		"sec-ch-ua-mobile":          {"?0"},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36"},
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"sec-fetch-site":            {"none"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-user":            {"?1"},
		"sec-fetch-dest":            {"document"},
		"accept-encoding":           {"gzip, deflate, br"},
		http.HeaderOrderKey: {
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"upgrade-insecure-requests",
			"user-agent",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-user",
			"sec-fetch-dest",
			"accept-encoding",
		},
	}

	trace := &httptrace.ClientTrace{
		TLSHandshakeDone: func(cs tls.ConnectionState, e error) {
			fmt.Printf("TLS Handshake: %v", cs)
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	resp, err := client.Do(req)

	if err != nil {
		t.Errorf(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %v", resp.StatusCode)
	}

	var data interface{}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Error(err.Error())
	}
}
