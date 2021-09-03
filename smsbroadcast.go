package smsbroadcast

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const sdkVersion = "0.1.0"
const baseUrlString = "https://www.smsbroadcast.com.au/api-adv.php"

type Client struct {
	BaseUrl *url.URL `url:"-"`

	Username   string       `url:"username"`
	Password   string       `url:"password"`
	httpClient *http.Client `url:"-"`
	userAgent  string       `url:"-"`
}

type ClientOptions struct {
	HttpClient *http.Client
}

type Message struct {
	To      string `url:"to"`
	From    string `url:"from"`
	Message string `url:"message"`
	Ref     string `url:"ref,omitempty"`
}

type MsgResponse struct {
	Status    int
	Summary   string
	Recipient string
	Reference string
}

func NewClient(username, password string, options *ClientOptions) (sbcClient *Client, err error) {
	sbcClient = &Client{}

	if len(username) == 0 {
		username = os.Getenv("SMSBROADCAST_USERNAME")
	}
	if len(password) == 0 {
		password = os.Getenv("SMSBROADCAST_PASSWORD")
	}

	sbcClient.Username = username
	sbcClient.Password = password
	sbcClient.userAgent = fmt.Sprintf("%s/%s (Go: %s)", "smsbroadcast-go", sdkVersion, runtime.Version())

	baseUrl, err := url.Parse(baseUrlString)

	sbcClient.BaseUrl = baseUrl
	sbcClient.httpClient = &http.Client{
		Timeout: time.Minute,
	}

	if options.HttpClient != nil {
		sbcClient.httpClient = options.HttpClient
	}
	return
}

func (sbcClient *Client) newRequest(msg Message) (request *http.Request, err error) {
	var values url.Values
	requestUrl := *sbcClient.BaseUrl

	var buffer = new(bytes.Buffer)
	if values, err = query.Values(msg); err != nil {
		return
	}

	values.Add("username", sbcClient.Username)
	values.Add("password", sbcClient.Password)

	requestUrl.RawQuery = values.Encode()
	request, err = http.NewRequest("POST", requestUrl.String(), buffer)
	request.Header.Add("User-Agent", sbcClient.userAgent)

	return
}

func (sbcClient *Client) Send(msg Message) (result MsgResponse, err error) {
	var response *http.Response
	request, err := sbcClient.newRequest(msg)
	if err != nil {
		return
	}
	response, err = sbcClient.httpClient.Do(request)
	if err != nil {
		return
	}
	result, err = parseResponse(response)
	return
}

func parseResponse(response *http.Response) (result MsgResponse, err error) {
	result.Status = response.StatusCode
	// ReadCloser to bytes.Buffer (and later to string)
	bytesBuffer := new(bytes.Buffer)
	_, err = bytesBuffer.ReadFrom(response.Body)
	if err != nil {
		return
	}
	resParams := strings.Split(strings.ReplaceAll(bytesBuffer.String(), "\n", ""), ":")
	result.Summary = resParams[0]
	result.Recipient = resParams[1]
	result.Reference = resParams[2]

	return
}
