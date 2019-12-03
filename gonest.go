package gonest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Handler struct {
	apiKey  string
	apiURL  string
	lastGet time.Time
	nest    Nest
	mqtt    mqtt.Client
}

type Nest struct {
	Devices    Device               `json:"devices"`
	Structures map[string]Structure `json:"structures"`
}

var (
	apiURL = "https://developer-api.nest.com/"
)

func New() *Handler {
	nestAPIKey, ok := os.LookupEnv("NEST_API_KEY")
	if !ok {
		panic("missing environment key: NEST_API_KEY")
	}

	log.Printf("NEST_API_KEY: %s*****", nestAPIKey[0:5])

	return &Handler{
		apiKey: nestAPIKey,
		apiURL: apiURL,
	}
}

func (h *Handler) ClearCache() {
	h.lastGet = time.Time{}
}

func (h *Handler) Get() (Nest, error) {
	if h.lastGet.Add(1 * time.Minute).After(time.Now()) {
		log.Printf("nest: using cached values")
		return h.nest, nil
	}
	log.Printf("nest: calling api for new values")

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

	res, err := client.Do(req)
	if err != nil {
		return h.nest, err
	}

	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&h.nest)
	log.Printf("nest: GET body response from api: %v", h.nest)
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

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	log.Printf("nest: PUT body response from api: %v", body)

	return nil
}
