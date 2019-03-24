package tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bidease/spl/config"
)

// Request ..
func Request(path string, out interface{}) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://portal.servers.com/rest/%s", path), nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Email", config.Options.Email)
	req.Header.Set("X-User-Token", config.Options.Token)

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	err = json.Unmarshal(body, &out)
	if err != nil {
		log.Fatalln(err)
	}
}
