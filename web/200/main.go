package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func main() {
	jar, _ := cookiejar.New(nil)
	c := &http.Client{Jar: jar}

	res, err := c.Get("http://a5540f88cd4c911e8857e06d564d8652-327183582.us-west-2.elb.amazonaws.com:4738")
	if err != nil {
		panic(err)
	}

	cookie := []byte(res.Cookies()[0].Value)

	u, _ := url.Parse("http://a5540f88cd4c911e8857e06d564d8652-327183582.us-west-2.elb.amazonaws.com:4738")

	for i := range cookie {
		cookie := c.Jar.Cookies(u)
		cookie[0].Value = strings.Repeat("a", i) + cookie[0].Value[i:]
		c.Jar.SetCookies(u, cookie)
		fmt.Println(cookie)

		res, err := c.Get("http://a5540f88cd4c911e8857e06d564d8652-327183582.us-west-2.elb.amazonaws.com:4738")
		if err != nil || res.StatusCode == 500 {
			fmt.Println(i, err)
			fmt.Println(res)
		}
	}
}
