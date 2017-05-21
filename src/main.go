package main

import (
	"log"
	"io"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"os"
	"bufio"
	"strings"
	"sync"
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

var wg sync.WaitGroup
var ColorList = make(map[string]string)
var VariantList = make(map[string]string)
var WearList = make(map[string]string)
type Brand []string
var BrandList Brand

type Pair struct {
	Keyword		string 		`json:"keyword"`
	Score 		float32		`json:"score"`
}
type Pairs []Pair
type ResultForApp struct{
	AllPair		Pairs 		`json:"pairs"`
}

func main() {
	loadMeta()
	http.HandleFunc("/label", giveLabel)
	http.ListenAndServe(":8000", nil)
}

func giveLabel(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// limit the size of incoming req body
	body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	// resp, err := http.Get("http://0.0.0.0:8080/ws")
	// if err != nil {
	// 	log.Printf("1 error: %v", err)
	// }
	// defer resp.Body.Close()
	
	var webdec WebDetection

	// json.NewDecoder(resp.Body).Decode(&webdec)
	json.Unmarshal(body, &webdec)

	wg.Add(4)
	topBrand := findBrand(webdec)
	topColor := findColor(webdec)
	topVariant := findVariant(webdec)
	topWear := findWear(webdec)
	wg.Wait()

	assembled := assembleKeyword(topBrand, topColor, topVariant, topWear)
	result := ResultForApp{
		AllPair: 	assembled}

	log.Printf("result: %v", result)
	json.NewEncoder(w).Encode(result)
}

func assembleKeyword(topBrand Pair, topColor Pair, topVariant Pair, topWear Pair) Pairs{
	var assembled Pairs
	var detailedA Pair
	var detailedB Pair
	var detailedC Pair
	var detailedD Pair
	if(topBrand.Score>0 && topColor.Score>0 && topVariant.Score>0 && topWear.Score>0){
		detailedA.Score = (topBrand.Score + topColor.Score + topVariant.Score + topWear.Score)/4
		detailedA.Keyword = topBrand.Keyword + " " + topColor.Keyword + " " + topVariant.Keyword + " " + topWear.Keyword

		assembled = append(assembled,detailedA)
	}

	if(topBrand.Score>0 && topVariant.Score>0 && topWear.Score>0){
		detailedB.Score = (topBrand.Score + topVariant.Score + topWear.Score)/3
		detailedB.Keyword = topBrand.Keyword + " " + topVariant.Keyword + " " + topWear.Keyword

		assembled = append(assembled,detailedB)
	}

	if(topBrand.Score>0 && topWear.Score>0){
		detailedC.Score = (topBrand.Score + topWear.Score)/2
		detailedC.Keyword = topBrand.Keyword + " " + topWear.Keyword

		assembled = append(assembled,detailedC)
	}

	if(topVariant.Score>0 && topWear.Score>0){
		detailedD.Score = (topVariant.Score + topWear.Score)/2
		detailedD.Keyword = topVariant.Keyword + " " + topWear.Keyword

		assembled = append(assembled,detailedD)
	}

	return assembled
}

func findBrand(webdec WebDetection) Pair{
	topBrand := Pair{
		Keyword: 	"",
		Score: 		0}

	for no := range webdec.WebEnt {
		if (stringIsInSlice(webdec.WebEnt[no].Description, BrandList)){
			if webdec.WebEnt[no].Score > topBrand.Score {
				topBrand.Keyword = webdec.WebEnt[no].Description
				topBrand.Score = webdec.WebEnt[no].Score
			}
		}
	}
	wg.Done()
	return topBrand
}

func findColor(webdec WebDetection) Pair{
	topColor := Pair{
		Keyword: 	"",
		Score: 		0}

	for no := range webdec.WebEnt {
		if translate, ok := ColorList[webdec.WebEnt[no].Description]; ok{
			if webdec.WebEnt[no].Score > topColor.Score {
				topColor.Keyword = translate
				topColor.Score = webdec.WebEnt[no].Score
			}
		}
	}
	wg.Done()
	return topColor
}

func findVariant(webdec WebDetection) Pair{
	topVariant := Pair{
		Keyword: 	"",
		Score: 		0}

	for no := range webdec.WebEnt {
		if translate, ok := VariantList[webdec.WebEnt[no].Description]; ok{
			if webdec.WebEnt[no].Score > topVariant.Score {
				topVariant.Keyword = translate
				topVariant.Score = webdec.WebEnt[no].Score
			}
		}
	}
	wg.Done()
	return topVariant
}

func findWear(webdec WebDetection) Pair{
	topWear := Pair{
		Keyword: 	"",
		Score: 		0}

	for no := range webdec.WebEnt {
		if translate, ok := WearList[webdec.WebEnt[no].Description]; ok{
			if webdec.WebEnt[no].Score > topWear.Score {
				topWear.Keyword = translate
				topWear.Score = webdec.WebEnt[no].Score
			}
		}
	}
	wg.Done()
	return topWear
}

func loadMeta(){
	wg.Add(4)
	go loadBrand("../util/brand.txt")
	go loadColor("../util/color.txt")
	go loadVariant("../util/variant.txt")
	go loadWear("../util/wear.txt")
	wg.Wait()
}

func loadBrand(file string) {
	f, _ := os.Open(file)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
	    line := scanner.Text()	    
		BrandList = append(BrandList,line)
	}
	wg.Done()
}

func loadColor(file string) {
	f, _ := os.Open(file)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
	    line := scanner.Text()	    
	    words := strings.Split(line,":")
		ColorList[words[0]] = words[1]
	}
	wg.Done()
}

func loadVariant(file string) {
	f, _ := os.Open(file)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
	    line := scanner.Text()	    
	    words := strings.Split(line,":")
		VariantList[words[0]] = words[1]
	}
	wg.Done()
}

func loadWear(file string) {
	f, _ := os.Open(file)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
	    line := scanner.Text()	    
	    words := strings.Split(line,":")
		WearList[words[0]] = words[1]
	}
	wg.Done()
}

func stringIsInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}