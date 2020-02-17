package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

type ArtistData struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type RelationData struct {
	RIndex []RelIndex `json:"index"`
}

type RelIndex struct {
	ID             int64               `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type BigData struct {
	GroupieArtist   []ArtistData
	GroupieRelation RelationData
}

func main() {
	fmt.Println("Starting the application...")

	respose1, err1 := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	respose2, err2 := http.Get("https://groupietrackers.herokuapp.com/api/relation")

	if err1 != nil || err2 != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err1)
	}

	data1, _ := ioutil.ReadAll(respose1.Body)
	bytes1 := []byte(data1)
	var Art []ArtistData
	json.Unmarshal(bytes1, &Art)

	data2, _ := ioutil.ReadAll(respose2.Body)
	bytes2 := []byte(data2)
	var Rel RelationData
	json.Unmarshal(bytes2, &Rel)

	ExecData := BigData{Art, Rel}

	fmt.Println(ExecData.GroupieRelation.RIndex[0])

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			tmpl, _ := template.ParseFiles("index.html")
			if r.Method == "GET" {
				tmpl.Execute(w, Art)
			}
		} else if r.URL.Path == "/relation" {
			tmpl, _ := template.ParseFiles("relation.html")
			if r.Method == "GET" {
				tmpl.Execute(w, ExecData)
			}
		} else {
			fmt.Fprintln(w, "ERROR 404")
			return
		}
	})

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)

}
