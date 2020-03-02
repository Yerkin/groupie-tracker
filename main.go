package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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

type ForMember struct {
	Mem       string
	GroupName []string
}

type ForLoc struct {
	Loc  string
	City []string
}

type ForAlb struct {
	Year  string
	Group string
}

type ForDate struct {
	Year  int
	Group []string
}

func ToLast(str string) string {
	tmp := ""
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] == ' ' && tmp != "" {
			break
		} else if str[i] != ' ' {
			tmp = string(str[i]) + tmp
		}
	}
	return tmp
}

func FindId(art []ArtistData, s string) int {
	for i := range art {
		if strings.ToLower(art[i].Name) == strings.ToLower(s) {
			return art[i].ID
		}
	}
	return 0
}

func FindMem(art []ArtistData, s string) []string {
	var ans []string
	for i := range art {
		for _, v := range art[i].Members {
			if v == s {
				ans = append(ans, art[i].Name)
			}
		}
	}
	return ans
}

func FindLoc(exe BigData, s string) []string {
	var ans []string
	for i := range exe.GroupieArtist {
		for key := range exe.GroupieRelation.RIndex[i].DatesLocations {
			if strings.ToLower(key) == strings.ToLower(s) {
				ans = append(ans, exe.GroupieArtist[i].Name)
			}
		}
	}
	return ans
}

func FindAlb(art []ArtistData, y string) string {
	for i := range art {
		if art[i].FirstAlbum == y {
			return art[i].Name
		}
	}
	return "-1"
}

func FindDate(art []ArtistData, y int) []string {
	var ans []string
	for i := range art {
		if art[i].CreationDate == y {
			ans = append(ans, art[i].Name)
		}
	}
	return ans
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

	// fmt.Println(ExecData.GroupieRelation.RIndex[0])

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			if r.Method == "GET" {
				tmpl, _ := template.ParseFiles("index.html")
				tmpl.Execute(w, ExecData)
			}
			if r.Method == "POST" {
				s := r.FormValue("toSearch")
				check := 0
				for _, e := range s {
					if e == '/' {
						check++
					}
				}
				if check == 2 {
					tmp := ToLast(s)
					s = s[:(len(s) - len(tmp) - 4)]

					if tmp == "Artists" {
						ArtId := FindId(Art, s)
						if ArtId > 0 && ArtId < 53 {
							tmpl, _ := template.ParseFiles("relation.html")
							tmpl.Execute(w, Art[ArtId-1])
						} else {
							tmpl, _ := template.ParseFiles("error.html")
							tmpl.Execute(w, nil)
						}
					} else if tmp == "Member" {
						var myMem []string
						myMem = FindMem(Art, s)
						if myMem != nil {
							MemData := ForMember{s, myMem}
							tmpl, _ := template.ParseFiles("mem.html")
							tmpl.Execute(w, MemData)
						} else {
							tmpl, _ := template.ParseFiles("error.html")
							tmpl.Execute(w, nil)
						}

					} else if tmp == "Location" {
						var myCit []string
						myCit = FindLoc(ExecData, s)
						if myCit != nil {
							LocData := ForLoc{s, myCit}
							tmpl, _ := template.ParseFiles("loc.html")
							tmpl.Execute(w, LocData)
						} else {
							tmpl, _ := template.ParseFiles("error.html")
							tmpl.Execute(w, nil)
						}

					} else if tmp == "Album" {
						s = s[:(len(s) - 6)]
						gn := FindAlb(Art, s)
						fmt.Println(s)
						if s != "-1" {
							AlbData := ForAlb{s, gn}
							tmpl, _ := template.ParseFiles("alb.html")
							tmpl.Execute(w, AlbData)
						} else {
							tmpl, _ := template.ParseFiles("error.html")
							tmpl.Execute(w, nil)
						}
					} else if tmp == "CreationDate" {
						var myGroup []string
						d, Derr := strconv.Atoi(s)
						if Derr != nil {
							myGroup = FindDate(Art, d)
							if myGroup != nil {
								DatData := ForDate{d, myGroup}
								tmpl, _ := template.ParseFiles("dat.html")
								tmpl.Execute(w, DatData)
							} else {
								tmpl, _ := template.ParseFiles("error.html")
								tmpl.Execute(w, nil)
							}
						} else {
							tmpl, _ := template.ParseFiles("error.html")
							tmpl.Execute(w, nil)
						}

					} else {
						tmpl, _ := template.ParseFiles("error.html")
						tmpl.Execute(w, nil)
					}
				} else {
					tmpl, _ := template.ParseFiles("error.html")
					tmpl.Execute(w, nil)
				}
			}
		} else {
			fmt.Fprintln(w, "ERROR 404")
			return
		}
	})

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)

}
