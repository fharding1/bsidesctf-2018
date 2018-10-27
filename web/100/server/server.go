package main

import (
	"encoding/base64"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

func indexOf(arr []string, s string) int {
	for i, v := range arr {
		if v == s {
			return i
		}
	}

	return -1
}

func shuffle(slice []string) {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

const badBody = `<!doctype html>
<title>Knock knock knock....</title>
<h1>Go away</h1>
<p>I have not heard the secret knock yet</p>

<!-- /static/dict.txt -->`

const flag = `BSidesPDX{ThIs_ChAlLeNgE_iSn'T_hArD}`

func main() {
	rand.Seed(time.Now().Unix()) // vvv secure

	knocks := make(map[string]int)

	dictBytes, err := ioutil.ReadFile("dict.txt")
	if err != nil {
		panic(err)
	}

	dictWordlist := strings.Split(string(dictBytes), "\n")
	shuffle(dictWordlist)
	wordlist := dictWordlist[:1000]

	indicies := make(map[string]int, len(dictWordlist))
	for _, word := range dictWordlist {
		if index := indexOf(wordlist, word); index != -1 {
			indicies[word] = index
		}
	}

	http.HandleFunc("/static/dict.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write(dictBytes)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var session string

		if sessionCookie, err := r.Cookie("session"); err != nil {
			generatedSession := uuid.Must(uuid.NewV4()).String()
			w.Header().Add("Set-Cookie", "session="+generatedSession)
			session = generatedSession
		} else {
			session = sessionCookie.Value
		}

		n, ok := knocks[session]
		if !ok {
			knocks[session] = 0
			n = 0
		}

		if knocks[session] == len(wordlist) {
			w.Write([]byte(flag))
			return
		}

		word := strings.TrimLeft(r.URL.String(), "/")
		if word == wordlist[n] {
			knocks[session]++
		} else {
			knocks[session] = 0
		}

		index, ok := indicies[word]
		if !ok {
			w.Write([]byte(badBody))
			return
		}

		w.Write([]byte(base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(index + 1)))))
	})

	http.ListenAndServe(":43478", nil)
}
