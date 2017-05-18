package main

import (
	"log"
	_"io/ioutil"
	"encoding/json"
	"net/http"
)

type WebEntity struct {
    EntityId        string		`json:"entityId"`
    Score  			float32		`json:"score"`
    Description		string		`json:"description"`
}
type WebEntities []WebEntity
type WebDetection struct {
	WebEnt 		WebEntities			`json:"webEntities"`
}

func main() {
	http.HandleFunc("/label", giveLabel)
	http.ListenAndServe(":8000", nil)
}

func giveLabel(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp, err := http.Get("http://0.0.0.0:8080/ws")
	if err != nil {
		log.Printf("1 error: %v", err)
	}
	defer resp.Body.Close()
	
	var webdec WebDetection

	json.NewDecoder(resp.Body).Decode(&webdec)
	// log.Printf("result: %v", webdec)
	json.NewEncoder(w).Encode(webdec)
}
