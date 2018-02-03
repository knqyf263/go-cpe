package main

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
)

// CpeList is CPE list
type CpeList struct {
	CpeItems []CpeItem `xml:"cpe-item"`
}

// CpeItem has CPE information
type CpeItem struct {
	Name      string    `xml:"name,attr"`
	Cpe23Item Cpe23Item `xml:"cpe23-item"`
}

// Cpe23Item has CPE 2.3 information
type Cpe23Item struct {
	Name string `xml:"name,attr"`
}

// Pair has fs and uri
type Pair struct {
	URI string
	FS  string
}

func main() {
	url := "http://static.nvd.nist.gov/feeds/xml/cpe/dictionary/official-cpe-dictionary_v2.3.xml.gz"
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		fmt.Printf("HTTP error. errs: %s, url: %s", err, url)
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	b := bytes.NewBufferString(string(body))
	reader, err := gzip.NewReader(b)
	defer reader.Close()
	if err != nil {
		fmt.Printf("Failed to decompress NVD feedfile. url: %s, err: %s", url, err)
		return
	}
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("Failed to Read NVD feedfile. url: %s, err: %s", url, err)
		return
	}
	cpeList := CpeList{}
	if err = xml.Unmarshal(bytes, &cpeList); err != nil {
		fmt.Printf("Failed to unmarshal. url: %s, err: %s", url, err)
		return
	}

	var uriList, fsList []string
	for _, cpeItem := range cpeList.CpeItems {
		uriList = append(uriList, cpeItem.Name)
		fsList = append(fsList, cpeItem.Cpe23Item.Name)
	}
	shuffle(fsList)

	pair := []Pair{}
	for i, uri := range uriList {
		pair = append(pair, Pair{
			URI: uri,
			FS:  fsList[i],
		})
	}
	fmt.Printf("%d data...\n", len(cpeList.CpeItems))

	fmt.Println("Generating test code...")
	t := template.Must(template.ParseFiles("dictionary_test.tmpl"))
	file, _ := os.Create(`./dictionary_test.go`)
	defer file.Close()
	t.Execute(file, map[string]interface{}{
		"Pair": pair,
	})
}

func shuffle(data []string) {
	n := len(data)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		data[i], data[j] = data[j], data[i]
	}
}
