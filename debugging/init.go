package debugging

import (
	"log"
	"net/http"
	"io/ioutil"
)

func ProbHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://127.0.0.1:1")
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)	
	if err != nil {
		log.Println(err)
	}

	w.Write(body)
}
