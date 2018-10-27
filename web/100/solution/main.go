package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
)

const (
	url      = "http://localhost:43478/"
	dictPath = "dict.txt"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	bytes, _ := ioutil.ReadFile(dictPath)
	words := strings.Split(string(bytes), "\n")

	jar, _ := cookiejar.New(nil)
	http.DefaultClient.Jar = jar

	foundWords := make([]string, 1000)

	fmt.Println("Initial ordering requests: ")
	for i, word := range words {
		fmt.Printf("\rProgress: %.2f%%", float64(i)/float64(len(words))*100)

		res, err := http.Get(url + word)
		check(err)

		bytes, err := ioutil.ReadAll(res.Body)
		check(err)
		check(res.Body.Close())

		numString, err := base64.StdEncoding.DecodeString(string(bytes))
		if err != nil || len(numString) == 0 {
			continue
		}

		n, err := strconv.Atoi(string(numString))
		check(err)

		foundWords[n-1] = word
	}
	fmt.Println()

	fmt.Println("Killing by 1000 curls:")
	for i, word := range foundWords {
		fmt.Printf("\rProgress: %.2f%%", float64(i)/float64(len(foundWords))*100)
		_, err := http.Get(url + word)
		check(err)
	}
	fmt.Println()

	res, err := http.Get(url)
	check(err)
	_, err = io.Copy(os.Stdout, res.Body)
	check(err)
	check(res.Body.Close())
}
