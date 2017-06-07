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
	"database/sql"
	"runtime"
	"time"
	"strconv"

	_ "github.com/go-sql-driver/mysql"	
	instagot "github.com/kahalKevin/instagot"
	gotrends "github.com/kahalKevin/gotrends/gotrends"
)

var Db *sql.DB //For database things

type WebEntity struct {
    EntityId        string		`json:"entityId"`
    Score  			float32		`json:"score"`
    Description		string		`json:"description"`
}
type WebEntities []WebEntity
type WebDetection struct {
	WebEnt 		WebEntities		`json:"webEntities"`
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

type Behaviour struct{
	Jenis 		string 		`jenis`
	Merk 		string 		`merk`
	Waktu		time.Time 	`waktu`
}
type Behaviours []Behaviour

type IG_url struct{
	Image_url 		string  	`json:"image_url"`
}

func init() {
    runtime.GOMAXPROCS(runtime.NumCPU())

	var err error
	Db, err = sql.Open("mysql", "root:mel@tcp(127.0.0.1:3306)/snapper?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	loadMeta()
	http.HandleFunc("/label", giveLabel)
	http.HandleFunc("/behaviour", getBehaviourByWeeks)
	http.HandleFunc("/image", GetImage)
	http.HandleFunc("/urlimage", GetUrlImage)	
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

	var topBrand Pair
	var topColor Pair
	var topVariant Pair
	var topWear Pair
	
	wg.Add(4)
	go func(){topBrand = findBrand(webdec)}()
	go func(){topColor = findColor(webdec)}()
	go func(){topVariant = findVariant(webdec)}()
	go func(){topWear = findWear(webdec)}()
	wg.Wait()

	assembled := assembleKeyword(topBrand, topColor, topVariant, topWear)
	result := ResultForApp{
		AllPair: 	assembled}

	if result.AllPair == nil {
		log.Println("Using google trends")
		topScore := findTopScore(webdec)
		gotrendWord, gotrendScore := gotrends.SearchWithKeyword(topScore.Keyword)
		gotrendPair:= Pair{
			Keyword: 	gotrendWord,
			Score: 		float32(gotrendScore)}
		_tempAllPair := result.AllPair
		_tempAllPair = append(_tempAllPair, gotrendPair)
		result.AllPair = _tempAllPair
	}

	if(topBrand.Keyword!="" || topColor.Keyword!=""){
		go logSearchBehaviour(topWear.Keyword,topBrand.Keyword)}

	log.Printf("result: %v", result)
	json.NewEncoder(w).Encode(result)
}

func assembleKeyword(topBrand Pair, topColor Pair, topVariant Pair, topWear Pair) Pairs{
	var assembled Pairs
	var detailedA Pair
	var detailedB Pair
	var detailedC Pair
	var detailedD Pair
	var detailedE Pair
		
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

	if(topVariant.Score>0 && topBrand.Score>0){
		detailedE.Score = (topVariant.Score + topBrand.Score)/2
		detailedE.Keyword = topVariant.Keyword + " " + topBrand.Keyword

		assembled = append(assembled,detailedE)
	}

	if(topBrand.Score>0 && detailedC==Pair{} && detailedE==Pair{}){
		assembled = append(assembled,topBrand)
	}

	if(topWear.Score>0 && detailedC==Pair{} && detailedD==Pair{}){
		assembled = append(assembled,topWear)
	}

	return assembled
}

func findTopScore(webdec WebDetection) Pair{
	topScore := Pair{
		Keyword: 	"",
		Score: 		0}

	for no := range webdec.WebEnt {
		if webdec.WebEnt[no].Score > topScore.Score {
			topScore.Keyword = webdec.WebEnt[no].Description
			topScore.Score = webdec.WebEnt[no].Score
		}
	}
	return topScore
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

func getBehaviourByWeeks(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var weeks uint32
	week_in_string := r.URL.Query().Get("weeks")
	if len(week_in_string) != 0 {
	    if week_converted, err := strconv.ParseUint(week_in_string,10,32); err == nil {
		    weeks = uint32(week_converted)
		}else{
			weeks = 1
		}
	}else{
		weeks = 1
	}

	err := Db.Ping()
	if err != nil {
		return
	}
	query_behaviour_by_week := `SELECT jenis, merk , datetime
							    FROM log_behaviour
								WHERE datetime>=(DATE_SUB(NOW(), INTERVAL ` +strconv.Itoa(int(7*weeks))+ ` DAY))`

	rows, err := Db.Query(query_behaviour_by_week)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var _Bhvs Behaviours
	var bhv1 Behaviour
	for rows.Next() {
		err := rows.Scan(&bhv1.Jenis, &bhv1.Merk, &bhv1.Waktu)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(bhv1)
		_Bhvs = append(_Bhvs, bhv1)
	}
	json.NewEncoder(w).Encode(_Bhvs)
}

func logSearchBehaviour(jenis string, merk string){
	err := Db.Ping()
	if err != nil {
		return
	}
	query_insert := "insert into log_behaviour(`jenis`, `merk`) values ('"+jenis+"', '"+merk+"')"
	_, err = Db.Exec(query_insert)
	if err != nil {
		log.Fatal(err)
	}
}

func GetImage(w http.ResponseWriter, r *http.Request){
	log.Println("Getting image from Instagram")
	ig_url := r.URL.Query().Get("ig_url")
	img := instagot.GetImage(ig_url)
	instagot.WriteImageToResponseWriter(w, &img)
}

func GetUrlImage(w http.ResponseWriter, r *http.Request){
	log.Println("Getting image url from Instagram")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	ig_url := r.URL.Query().Get("ig_url")
	img_url := instagot.GetUrlImage(ig_url)

	json.NewEncoder(w).Encode(
		IG_url{Image_url: 	img_url})
}