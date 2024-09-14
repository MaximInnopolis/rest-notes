package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var spellerAPIURL = "https://speller.yandex.net/services/spellservice.json/checkText"

// SpellerService represents service responsible for interacting with Speller API
type SpellerService struct {
	client *http.Client
}

// NewSpellerService creates new instance of SpellerService with predefined HTTP client timeout
func NewSpellerService() *SpellerService {
	return &SpellerService{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// CheckText sends text to Speller API to check for spelling errors
// It takes text to check, the language code, and options as parameters
// Returns spelling errors as list of maps or error if the request fails
func (s *SpellerService) CheckText(text, lang string, options int) ([]map[string]interface{}, error) {
	params := url.Values{}
	params.Add("text", text)
	params.Add("lang", lang)
	params.Add("options", string(options))
	params.Add("format", "plain")

	// Create new POST request with the encoded parameters
	req, err := http.NewRequest("POST", spellerAPIURL, bytes.NewBufferString(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get valid response from speller service")
	}

	// Unmarshal the JSON response into a slice of maps
	var result []map[string]interface{}
	if err = json.Unmarshal(body, &result); err != nil {
		log.Printf("Ошибка при парсинге JSON: %v", err)
		return nil, errors.New("failed to parse JSON response from speller service")
	}

	return result, nil
}
