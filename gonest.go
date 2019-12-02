package gonest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	apiKey  string
	apiURL  string
	lastGet time.Time
	nest    Nest
}

type Nest struct {
	Devices    Device               `json:"devices"`
	Structures map[string]Structure `json:"structures"`
}

var (
	apiURL = "https://developer-api.nest.com/"
)

func New(apikey string) *Handler {
	return &Handler{
		apiKey: apikey,
		apiURL: apiURL,
	}
}

func (h *Handler) ClearCache() {
	h.lastGet = time.Time{}
}

func (h *Handler) Get() (Nest, error) {
	if h.lastGet.Add(1 * time.Minute).After(time.Now()) {
		log.Printf("reusing value")
		return h.nest, nil
	}
	log.Printf("new requeset")

	req, _ := http.NewRequest("GET", h.apiURL, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.apiKey))

	client := &http.Client{}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		for key, val := range via[0].Header {
			req.Header[key] = val
		}
		return nil
	}

	log.Printf("req: %+v", req)
	res, err := client.Do(req)
	if err != nil {
		return h.nest, err
	}

	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&h.nest)
	log.Printf("Nest: %+v", h.nest)
	h.lastGet = time.Now()
	return h.nest, nil
}

func (h *Handler) Set(path, data string) error {
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s%s", h.apiURL, path), strings.NewReader(data))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.apiKey))

	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		for key, val := range via[0].Header {
			req.Header[key] = val
		}
		return nil
	}

	log.Printf("req: %+v", req)
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	log.Printf("result: %s", body)

	return nil
}
