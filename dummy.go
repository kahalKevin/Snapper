package main

import (
	_ "log"
	"encoding/json"
	"net/http"
)

type WebEntity struct {
    EntityId        string		`json:"entityId"`
    Score  			float32		`json:"score"`
    Description		string		`json:"description"`
}
type WebEntities []WebEntity

type FullMatchingImage struct {
	Url		string		`json:"url"`
	Dummy 	int32		`json:"dumm"`
}
type FullMatchingImages []FullMatchingImage

type WebDetection struct {
	WebEnt 		WebEntities			`json:"webEntities"`
	FullMatch 	FullMatchingImages	`json:"fullMatchingImages"`
}

func main() {
	http.HandleFunc("/ws", defaultjob)
	http.ListenAndServe(":8080", nil)
}

func defaultjob(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	ent1 := WebEntity{
		EntityId:		"/m/06rrc",
		Score:			0.85754,
		Description:	"Shoe"}

	ent2 := WebEntity{
		EntityId:		"/m/0lwkh",
		Score:			0.49217,
		Description:	"Nike"}

	ent3 := WebEntity{
		EntityId:		"/m/019sc",
		Score:			0.3900155,
		Description:	"Black"}

	ent4 := WebEntity{
		EntityId:		"/m/019rjn",
		Score:			0.37383,
		Description:	"Futsal"}

	ent5 := WebEntity{
		EntityId:		"/g/12dpwwx05",
		Score:			0.31789,
		Description:	"Nike Hypervenom"}

	ent6 := WebEntity{
		EntityId:		"/m/027sf8d",
		Score:			0.30877,
		Description:	"Nike Mercurial Vapor"}

	ent7 := WebEntity{
		EntityId:		"/m/04lbp",
		Score:			0.25490758,
		Description:	"Leather"}

	ent8 := WebEntity{
		EntityId:		"/m/01sdr",
		Score:			0.18925,
		Description:	"Color"}

	ent9 := WebEntity{
		EntityId:		"/g/11dylp_1v",
		Score:			0.18442,
		Description:	"Sneakers"}

	ent10 := WebEntity{
		EntityId:		"/m/01g5v",
		Score:			0.17937,
		Description:	"Blue"}

	ent11 := WebEntity{
		EntityId:		"/m/092sx5",
		Score:			0.1774,
		Description:	"Pricing strategies"}

	ent12 := WebEntity{
		EntityId:		"/m/038hg",
		Score:			0.17397,
		Description:	"Green"}

	ent13 := WebEntity{
		EntityId:		"/m/05t5gr",
		Score:			0.17268,
		Description:	"Promotion"}

	ent14 := WebEntity{
		EntityId:		"/m/06ntj",
		Score:			0.17248,
		Description:	"Sports"}

	ent15 := WebEntity{
		EntityId:		"/m/083jv",
		Score:			0.17213,
		Description:	"White"}

	var ents WebEntities
	ents = append(ents, ent1)
	ents = append(ents, ent2)
	ents = append(ents, ent3)
	ents = append(ents, ent4)
	ents = append(ents, ent5)
	ents = append(ents, ent6)
	ents = append(ents, ent7)
	ents = append(ents, ent8)
	ents = append(ents, ent9)
	ents = append(ents, ent10)
	ents = append(ents, ent11)
	ents = append(ents, ent12)
	ents = append(ents, ent13)
	ents = append(ents, ent14)
	ents = append(ents, ent15)

	fullMatch1 := FullMatchingImage{
		Url:	"https://s0.bukalapak.com/img/037091828/m-1000-1000/IMG_20170104_WA0004_scaled.jpg",
		Dummy:	1}

	fullMatch2 := FullMatchingImage{
		Url:	"https://s1.bukalapak.com/img/1583582511/m-1000-1000/14326926_61f6a647_1150_4d16_a4f3_bba7f259b753_444_444.jpg",
		Dummy:	2}

	var fullmatchs FullMatchingImages
	fullmatchs = append(fullmatchs, fullMatch1)
	fullmatchs = append(fullmatchs, fullMatch2)

	webdect := WebDetection{
		WebEnt:		ents,
		FullMatch: 	fullmatchs}

	json.NewEncoder(w).Encode(webdect)
}
