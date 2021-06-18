package ftx

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const FTXURL = "https://ftx.com/api/"

type Response struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
}

type Client struct {
	Client     *http.Client
	key        string
	secret     []byte
	Subaccount string
}

func NewClient(cl *http.Client, key, secret, sub string) *Client {
	c := &Client{cl, key, []byte(secret), sub}
	return c
}

func processResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error processing response: %v", err)
		return err
	}
	err = json.Unmarshal(body, result)
	if err != nil {
		log.Printf("Error processing response: %v", err)
		return err
	}
	return nil
}

func (pr *Client) sign(signaturePayload string) string {
	mac := hmac.New(sha256.New, pr.secret)
	mac.Write([]byte(signaturePayload))
	return hex.EncodeToString(mac.Sum(nil))
}

func (pr *Client) signRequest(method string, path string, body []byte) *http.Request {
	ts := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	signaturePayload := ts + method + "/api/" + path + string(body)
	signature := pr.sign(signaturePayload)
	req, _ := http.NewRequest(method, FTXURL+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("FTX-KEY", pr.key)
	req.Header.Set("FTX-SIGN", signature)
	req.Header.Set("FTX-TS", ts)
	if pr.Subaccount != "" {
		req.Header.Set("FTX-SUBACCOUNT", pr.Subaccount)
	}
	return req
}

func (client *Client) get(path string, body []byte) (*http.Response, error) {
	preparedRequest := client.signRequest("GET", path, body)
	resp, err := client.Client.Do(preparedRequest)
	return resp, err
}

func (client *Client) post(path string, body []byte) (*http.Response, error) {
	preparedRequest := client.signRequest("POST", path, body)
	resp, err := client.Client.Do(preparedRequest)
	return resp, err
}

func (client *Client) delete(path string, body []byte) (*http.Response, error) {
	preparedRequest := client.signRequest("DELETE", path, body)
	resp, err := client.Client.Do(preparedRequest)
	return resp, err
}
