package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Element struct {
	Status Status `json:"status"`
}

func outputHTML(w http.ResponseWriter, filename string, data interface{}) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	element := readFromFile()
	statusWater, classWater, statusWind, classWind := checkStatus(element.Status.Water, element.Status.Wind)
	myvar := map[string]interface{}{
		"statusWater": statusWater,
		"statusWind":  statusWind,
		"water":       element.Status.Water,
		"wind":        element.Status.Wind,
		"classWater":  classWater,
		"classWind":   classWind,
	}
	outputHTML(w, "./login.html", myvar)
}

func main() {
	runtime.GOMAXPROCS(2)

	http.HandleFunc("/", loginHandler)
	go doEvery(15*time.Second, writeToFile)

	var address = "localhost:8080"
	fmt.Printf("server started at %s\n", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func randomNumber() int {
	min := 1
	max := 100
	return rand.Intn(max-min) + min
}

func writeToFile(t time.Time) {
	status := Status{
		Water: randomNumber(),
		Wind:  randomNumber(),
	}

	element := Element{
		Status: status,
	}

	content, err := json.Marshal(element)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("file.json", content, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Water :", status.Water, " Wind :", status.Wind, ", time :", t)
}

func readFromFile() Element {
	content, err := ioutil.ReadFile("file.json")
	if err != nil {
		log.Fatal(err)
	}

	element := Element{}
	err = json.Unmarshal(content, &element)
	if err != nil {
		log.Fatal(err)
	}

	return element
}

func checkStatus(water, wind int) (statusWater, classWater, statusWind, classWind string) {
	if water < 5 {
		statusWater = fmt.Sprintf("Aman")
		classWater = "success"
	} else if water < 8 {
		statusWater = fmt.Sprintf("Siaga")
		classWater = "warning"
	} else {
		statusWater = fmt.Sprintf("Bahaya")
		classWater = "danger"
	}

	if wind < 6 {
		statusWind = fmt.Sprintf("Aman")
		classWind = "success"
	} else if wind < 15 {
		statusWind = fmt.Sprintf("Siaga")
		classWind = "warning"
	} else {
		statusWind = fmt.Sprintf("Bahaya")
		classWind = "danger"
	}

	return
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}
