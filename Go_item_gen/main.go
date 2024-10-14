package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Item struct {
	Caption string  `json:"caption"`
	Weight  float64 `json:"weight"`
	Number  int32   `json:"number"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func randFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}

func SendItem() Item {
	var item = Item{
		Caption: randStringRunes(5),
		Weight:  randFloats(10, 200, 1)[0],
		Number:  rand.Int31(),
	}

	jsonItem, err := json.Marshal(item)
	if err != nil {
		fmt.Println(err)
		return item
	}

	posturl := "http://localhost:8080/item"
	req, err := http.NewRequest("POST", posturl, bytes.NewBuffer(jsonItem))
	if err != nil {
		fmt.Println(err)
		return item
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return item
	}

	defer resp.Body.Close()

	return item

}

func main() {
	itemCount := 5
	items := []Item{}
	for i := 0; i < itemCount; i++ {
		item := SendItem()
		items = append(items, item)
	}

	for _, i := range items {
		geturl := fmt.Sprintf("http://localhost:8080/item/%s", i.Caption)
		resp, err := http.Get(geturl)
		if err != nil {
			fmt.Println(err)
			return
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		item := Item{}
		json.Unmarshal(body, &item)
		fmt.Printf("%v\n", item)
	}
}
