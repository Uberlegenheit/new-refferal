package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func SendRequest(res interface{}, u string, m string, b io.Reader, hdrs map[string]string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(m, u, b)
	if err != nil {
		return err
	}

	for k, v := range hdrs {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request: %s", resp.Status)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, res)
	if err != nil {
		return err
	}

	return nil
}

func HTTPGet(fullURL string, params map[string]string, dst interface{}) error {
	if len(params) != 0 {
		values := url.Values{}
		for key, value := range params {
			values.Add(key, value)
		}
		fullURL = fmt.Sprintf("%s?%s", fullURL, values.Encode())
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fullURL)
	if err != nil {
		return fmt.Errorf("http.Get: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request: %s", resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %s", err.Error())
	}
	err = json.Unmarshal(data, dst)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	return nil
}

func ValidateID(id uint64) error {
	if id == 0 {
		return fmt.Errorf("id can't be 0")
	}
	return nil
}

func ValidateEmail(email string) error {
	if len(email) < 5 || !strings.ContainsRune(email, '@') || !strings.ContainsRune(email, '.') {
		return fmt.Errorf("invalid email")
	}
	return nil
}
