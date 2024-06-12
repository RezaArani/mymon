package webmonitor

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/RezaArani/mymon/config"
	"github.com/RezaArani/mymon/wlogs"
)

// performance Metrics type
type QuickTest struct {
	Url                string `json:"url" cql:"url"`
	Interval           int    `json:"interval" cql:"interval"`
	AlarmFailTestCount int    `json:"alarmfailtestcount" cql:"alarmfailtestcount"`
	HttpMethod         string `json:"httpmethod" cql:"httpmethod"`
	IgnoreCertificate  bool   `json:"ignorecertificate" cql:"ignorecertificate"`

	ActiveFailures     int       `json:"activefailures" cql:"activefailures"`
	TestCount          int64     `json:"testcount" cql:"testcount"`
	IsHealthy          bool      `json:"healthy" cql:"healthy"`
	LastError          time.Time `json:"lasterror" cql:"lasterror"`
	TestTime           time.Time `json:"testtime" cql:"testtime"`
	HasError           bool      `json:"error" cql:"error"`
	StatusCode         string    `json:"statuscode" cql:"statuscode"`
	DnsTime            int       `json:"dnstime" cql:"dnstime"`
	SslTime            int       `json:"ssltime" cql:"ssltime"`
	TlsVersion         string    `tlsversion:"tlsversion" cql:"tlsversion"`
	WaitTime           int       `json:"waittime" cql:"waittime"`
	ReceiveTime        int       `json:"receivetime" cql:"receivetime"`
	ErrorText          string    `json:"errortext" cql:"errortext"`
	Proto              string    `json:"Proto" cql:"Proto"`
	PageSize           int32     `json:"Pagesize" cql:"pagesize"`
	CompressedbodySize int32     `json:"compressedbodysize" cql:"compressedbodysize"`
	Compressed         bool      `json:"compressed" cql:"compressed"`
	LastServerId       string    `json:"lastserverid" cql:"lastserverid"`
	ConnectTime        int       `json:"connecttime" cql:"connecttime"`
	DnsError           bool      `json:"dnsError" cql:"dnsError"`
	NetError           bool      `json:"netError" cql:"netError"`
	AppError           bool      `json:"appError" cql:"appError"`
	SlowLoad           bool      `json:"SlowLoad" cql:"SlowLoad"`
	Netrefused         bool      `json:"netrefused" cql:"netrefused"`
	Nettimedout        bool      `json:"nettimedout" cql:"nettimedout"`
	SslError           bool      `json:"sslerror" cql:"sslerror"`
	RedirectedUrl      string    `json:"redirectedurl" cql:"redirectedurl"`
	RedirectError      bool      `json:"redirecterror" cql:"redirecterror"`
	RedirectsFollowed  int       `json:"redirectsfollowed" cql:"redirectsfollowed"`
	Site_Ip            string    `json:"site_ip" cql:"site_ip"`
	LastWebhookState   int       //0 not called, -1 Error,1 OK
}

func isRedirect(resp *http.Response) bool {
	return resp.StatusCode == http.StatusPermanentRedirect || resp.StatusCode == http.StatusTemporaryRedirect || resp.StatusCode == 301 || resp.StatusCode == 302
}

// func readResponseBody(req *http.Request, resp *http.Response) string {
// 	if isRedirect(resp) || req.Method == http.MethodHead {
// 		return ""
// 	}
// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(resp.Body)
// 	respBytes := buf.String()
// 	respString := string(respBytes)
// 	return respString // msg
// }

func createBody(body string) io.Reader {
	if strings.HasPrefix(body, "@") {
		filename := body[1:]
		f, err := os.Open(filename)
		if err != nil {
			// log.Fatalf("failed to open data file %s: %v", filename, err)
		}
		return f
	}
	return strings.NewReader(body)
}

func headerKeyValue(h string) (string, string) {
	i := strings.Index(h, ":")
	if i == -1 {
		log.Fatalf("Header '%s' has invalid format, missing ':'", h)
	}
	return strings.TrimRight(h[:i], " "), strings.TrimLeft(h[i:], " :")
}

func createRequest(method string, url *url.URL, body string, cookie []string) *http.Request {
	if url.Scheme == "http" && strings.HasSuffix(url.Host, ":443") {
		url.Scheme = "https"
		url.Host = strings.Replace(url.Host, ":443", "", 1)

	}
	if url.Scheme == "https" && strings.HasSuffix(url.Host, ":443") {
		url.Host = strings.Replace(url.Host, ":443", "", 1)
	}
	req, _ := http.NewRequest(method, url.String(), createBody(body))
	// req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")

	if cookie != nil {
		cookieStr := ""
		for _, element := range cookie {
			// index is the index where we are
			// element is the element from someSlice for where we are

			s := strings.Split(element, ";Expires=")
			//req.Header.Add("Cookie", s[0])
			if cookieStr != "" {
				cookieStr += "; "
			}
			cookieStr += s[0]
		}
		req.Header.Add("Cookie", cookieStr)
		// req.Header.Add("DNT", "1")
		// req.Header.Add("Upgrade-Insecure-Requests", "1")
		// req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		// req.Header.Add("Sec-Fetch-Site", "none")
		// req.Header.Add("Sec-Fetch-Mode", "navigate")
		// req.Header.Add("Sec-Fetch-User", "?1")
		// req.Header.Add("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`)
		// req.Header.Add("sec-ch-ua-platform", `"Windows"`)
		//req.Header.Add("Accept-Encoding", "gzip")
		// req.Header.Add("Accept-Language", "en-US,en;q=0.9")
		//req.Header.Add("Host", url)
		// req.Header.Add("Sec-Fetch-Dest", "document")
		// req.Header.Add("sec-ch-ua-mobile", "?0")
		// req.Header.Add("Connection", "keep-alive")

	}

	// if err != nil {
	// 	// log.Fatalf("unable to create request: %v", err)

	// }
	for _, h := range config.GetConfig().HttpHeaders {
		k, v := headerKeyValue(h)
		if strings.EqualFold(k, "host") {
			req.Host = v
			continue
		}
		req.Header.Add(k, v)
	}
	return req
}

func connectToSrv(network string) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, _, addr string) (net.Conn, error) {
		return (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: false,
		}).DialContext(ctx, network, addr)
	}
}

func (qt *QuickTest) getCert(filename string) []tls.Certificate {
	if filename == "" {
		return nil
	}
	var (
		pkeyPem []byte
		certPem []byte
	)

	// read client certificate file (must include client private key and certificate)
	certFileBytes, err := ioutil.ReadFile(config.GetConfig().ClientCertFile)
	if err != nil {
		// log.Printf("failed to read client certificate file: %v", err)
		qt.HasError = true
		qt.SslError = true
		qt.ErrorText += fmt.Sprintf("failed to read client certificate file: %v", err)
	}

	for {
		block, rest := pem.Decode(certFileBytes)
		if block == nil {
			break
		}
		certFileBytes = rest

		if strings.HasSuffix(block.Type, "PRIVATE KEY") {
			pkeyPem = pem.EncodeToMemory(block)
		}
		if strings.HasSuffix(block.Type, "CERTIFICATE") {
			certPem = pem.EncodeToMemory(block)
		}
	}

	cert, err := tls.X509KeyPair(certPem, pkeyPem)
	if err != nil {
		log.Fatalf("unable to load client cert and key pair: %v", err)
		qt.HasError = true
		qt.SslError = true

	}
	return []tls.Certificate{cert}
}

func ParseURL(uri string) *url.URL {
	if !strings.Contains(uri, "://") && !strings.HasPrefix(uri, "//") {
		uri = "//" + uri
	}

	url, err := url.Parse(uri)
	if err != nil {
		log.Fatalf("could not parse url %q: %v", uri, err)
	}

	if url.Scheme == "" {
		url.Scheme = "http"
		if !strings.HasSuffix(url.Host, ":80") {
			url.Scheme += "s"
		}
	}
	return url
}

func resetTestIntance(qt *QuickTest) {

	qt.IsHealthy = false

	qt.HasError = false
	qt.StatusCode = ""
	qt.DnsTime = 0
	qt.SslTime = 0
	qt.TlsVersion = ""
	qt.WaitTime = 0
	qt.ReceiveTime = 0
	qt.ErrorText = ""
	qt.Proto = ""
	qt.PageSize = 0
	qt.CompressedbodySize = 0
	qt.Compressed = false

	qt.ConnectTime = 0
	qt.DnsError = false
	qt.NetError = false
	qt.AppError = false
	qt.SlowLoad = false
	qt.Netrefused = false
	qt.Nettimedout = false
	qt.SslError = false
	qt.RedirectedUrl = ""
	qt.RedirectError = false
}

// visit visits a url and times the interaction.
// If the response is a 30x, visit follows the redirect.
func (qt *QuickTest) Visit(vurl *url.URL, cookies []string) {
	qt.TestTime = time.Now()
	resetTestIntance(qt)
	req := createRequest(qt.HttpMethod, vurl, config.GetConfig().PostBody, cookies)

	var t0, t1, t2, t3, t4, t5, t6 time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			t0 = time.Now()

		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {

			if dnsInfo.Err != nil {
				qt.ErrorText += dnsInfo.Err.Error()

			} else {
				t1 = time.Now()
			}

		},
		ConnectStart: func(_, ADDR string) {

			qt.Site_Ip = ADDR
			if t1.IsZero() {
				// connecting to IP
				t1 = time.Now()
			}

		},
		ConnectDone: func(net, addr string, err error) {

			if err != nil {
				qt.ErrorText += "unable to connect to host"
				errText := fmt.Sprintf("unable to connect to host %v: %v", addr, err)
				wlogs.LogMessage(errText, false)
				qt.HasError = true
				qt.NetError = true

			}
			t2 = time.Now()

		},
		GotConn: func(_ httptrace.GotConnInfo) {
			t3 = time.Now()

		},
		GotFirstResponseByte: func() {
			t4 = time.Now()

		},
		TLSHandshakeStart: func() {
			t5 = time.Now()

		},
		TLSHandshakeDone: func(_ tls.ConnectionState, tlsError error) {
			t6 = time.Now()
			if tlsError != nil {
				qt.SslError = true
			}

		},
	}
	req = req.WithContext(httptrace.WithClientTrace(context.Background(), trace))
	//	Accept-Encoding: gzip, deflate
	// req.Header.Set("accept-encoding", "gzip")
	// os.Setenv("HTTP_PROXY", "http://127.0.0.1:8889")
	// os.Setenv("HTTPS_PROXY", "http://127.0.0.1:8889")
	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	}

	switch {
	case config.GetConfig().FourOnly:
		tr.DialContext = connectToSrv("tcp4")
	case config.GetConfig().SixOnly:
		tr.DialContext = connectToSrv("tcp6")
	}

	switch vurl.Scheme {
	case "https":
		host, _, err := net.SplitHostPort(req.Host)
		if err != nil {
			host = req.Host
		}

		tr.TLSClientConfig = &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: qt.IgnoreCertificate,
			Certificates:       qt.getCert(config.GetConfig().ClientCertFile),
			MinVersion:         tls.VersionTLS12,
		}
	}

	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// always refuse to follow redirects, visit does that
			// manually if required.
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)

	if err != nil {
		wlogs.LogMessage(err, false)

		qt.HasError = true
		qt.ErrorText = err.Error()
		if strings.Contains(qt.ErrorText, "x509:") {
			qt.SslError = true
		}
		if strings.Contains(qt.ErrorText, "no such host") {
			qt.DnsError = true
		}
		if strings.Contains(qt.ErrorText, "machine actively refused it") {
			qt.Netrefused = true
		}
		if strings.Contains(qt.ErrorText, "cause the connected party did not properly respond after a period of time") {
			qt.Nettimedout = true
		}
	}

	if resp == nil || qt.HasError {

		qt.DnsTime = int(t1.Sub(t0) / time.Millisecond)
		qt.ConnectTime = int(t2.Sub(t1) / time.Millisecond)
		qt.SslTime = int(t6.Sub(t5) / time.Millisecond)
		qt.WaitTime = int(t4.Sub(t3) / time.Millisecond)

		if qt.ConnectTime < 0 {
			qt.ConnectTime = 0
		}

		if qt.SslTime < 0 {
			qt.SslTime = 0
		}

		if qt.WaitTime < 0 {
			qt.WaitTime = 0
		}

		// SendDataToCommunicator(QuickTestResult)

	} else {
		if resp.StatusCode < 299 || resp.StatusCode > 399 {
			if resp.StatusCode != 200 {
				qt.HasError = true
			} else {
				qt.IsHealthy = true
			}
			// switch resp.StatusCode {
			// case 502:
			// 	 "the server acting as the gateway received an invalid response from the main server."
			// case 504:
			//   "the server acting as the gateway didn’t receive a response at all from the application main server."
			// case 403:
			//   "The HTTP 403 Forbidden client error status response code indicates that the server understands the request but refuses to authorize it. The access is permanently forbidden and tied to the application logic, such as insufficient rights to a resource."
			// case 404:
			// 	 "The HTTP 404 Not Found client error response code indicates that the server can't find the requested resource. Links that lead to a 404 page are often called broken or dead links and can be subject to link rot."
			// case 503:
			// 	 "The server is currently unable to handle the request due to a temporary overloading or maintenance of the server. Note: The existence of the 503 status code does not imply that a server must use it when becoming overloaded. Some servers may wish to simply refuse the connection."
			// case 500:
			// 	 "Unexpected server error. "
			// }

			var bodyLength int64
			var fullBodyLength int64
			if !resp.Uncompressed {

				if resp.ContentLength == -1 {

					if resp.TransferEncoding != nil && len(resp.TransferEncoding) > 0 && resp.TransferEncoding[0] == "chunked" {
						buft := new(bytes.Buffer)
						buft.ReadFrom(resp.Body)
						newStr := buft.String()
						fullBodyLength = int64(len(newStr))
					}
					// else {
					// tmpIoReader := resp.Body

					// 	buf := &bytes.Buffer{}
					// 	c, _ := io.Copy(buf, tmpIoReader)
					// 	bodyLength = c
					// 	compreader := bytes.NewReader(buf.Bytes())
					// 	gzreader, _ := gzip.NewReader(compreader)
					// 	unzippedData := &bytes.Buffer{}
					// 	b, _ := io.Copy(unzippedData, gzreader)
					// 	fullBodyLength = b
					// }
				} else {
					bodyLength = resp.ContentLength
					if resp.StatusCode == 200 {
						gzreader, gzReaderErr := gzip.NewReader(resp.Body)
						if gzReaderErr != nil {
							// log.Println(gzReaderErr.Error())
							fullBodyLength = bodyLength
						} else {
							buf := &bytes.Buffer{}
							c, _ := io.Copy(buf, gzreader)
							fullBodyLength = c
						}

					} else {
						fullBodyLength = resp.ContentLength
					}
				}

			} else {
				fullBodyLength = resp.ContentLength
			}
			resp.Body.Close()

			t7 := time.Now() // after read body

			qt.CompressedbodySize = int32(bodyLength)
			qt.PageSize = int32(fullBodyLength)

			qt.Proto = resp.Proto
			qt.StatusCode = strconv.Itoa(resp.StatusCode)
			qt.Compressed = !resp.Uncompressed
			qt.HasError = resp.StatusCode != 200

			if resp.TLS != nil {
				qt.TlsVersion = strconv.Itoa(int(resp.TLS.Version))
			}

			if t0.IsZero() {
				// we skipped DNS
				t0 = t1
			}

			switch vurl.Scheme {
			case "https":

				qt.DnsTime = int(t1.Sub(t0) / time.Millisecond)
				qt.ConnectTime = int(t2.Sub(t1) / time.Millisecond)
				qt.SslTime = int(t6.Sub(t5) / time.Millisecond)
				qt.WaitTime = int(t4.Sub(t3) / time.Millisecond)
				qt.ReceiveTime = int(t7.Sub(t4) / time.Millisecond)

			case "http":

				qt.DnsTime = int(t1.Sub(t0) / time.Millisecond)
				qt.ConnectTime = int(t3.Sub(t1) / time.Millisecond)
				qt.SslTime = int(t6.Sub(t5) / time.Millisecond)
				qt.WaitTime = int(t4.Sub(t3) / time.Millisecond)
				qt.ReceiveTime = int(t7.Sub(t4) / time.Millisecond)
			}

		}
		if config.GetConfig().FollowRedirects && isRedirect(resp) {
			loc, err := resp.Location()
			reqCookies := (resp.Header["Set-Cookie"])

			if resp.Header.Get("Set-Cookie") != "" {
				wlogs.LogMessage("Cookies Requested!", false)

				cookies = reqCookies
			}
			if err != nil {
				if err == http.ErrNoLocation {
					// 30x but no Location to follow, give up.
					return
				}
				errText := "unable to follow redirect:" + fmt.Sprintf("unable to follow redirect: %v", err)
				errText += err.Error()
				qt.HasError = true
				qt.RedirectError = true
			}
			if isRedirect(resp) {
				qt.RedirectedUrl = loc.String()
			}
			qt.RedirectsFollowed++
			if qt.RedirectsFollowed > config.GetConfig().MaxRedirects {
				//maxRedirectsStr := strconv.Itoa(maxRedirects)
				//errText := fmt.Sprintf("maximum number of redirects (%d) followed", maxRedirects)
				// log.Printf(errText)
				qt.HasError = true
				qt.RedirectError = true
				qt.ErrorText += fmt.Sprint("maximum number of redirects (%d) followed", config.GetConfig().MaxRedirects)

			}

			qt.Visit(loc, cookies)
		}
	}

}
