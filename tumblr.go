package tumblr

import (
	"errors"
	"github.com/kurrik/oauth1a"
	"io/ioutil"
	"net/http"
	"net/url"
)

//API docs: http://www.tumblr.com/docs/en/api/v2#posts

const (
	TUMBLR_BASE_URL   = "https://api.tumblr.com/v2/blog/"
	REQUEST_URL       = "http://www.tumblr.com/oauth/request_token"
	AUTHORIZATION_URL = "http://www.tumblr.com/oauth/authorize"
	ACCESS_URL        = "http://www.tumblr.com/oauth/access_token"
)

type Jar struct {
	cookies []*http.Cookie
}

type Tumblr struct {
	RequestURL     string
	AuthorizeURL   string
	AccessURL      string
	ConsumerSecret string
	ConsumerKey    string
	CallbackURL    string
	BlogName       string
	Signer         oauth1a.Signer
	client         *http.Client
	jar            *Jar
}

func FetchResponseBody(resp *http.Response) []byte {
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		errors.New("status code: " + string(resp.StatusCode))
	}
	rbody, reader_error := ioutil.ReadAll(resp.Body)
	check_error(reader_error)
	return rbody
}

func check_error(err error) {
	if err != nil {
		panic(err)
	}
}

//Return ApiURL like http://api.tumblr.com/v2/blog/kracekumar
func (t Tumblr) ApiURL() string {
	if t.BlogName == "" {
		errors.New("Tumblr BlogName is missing")
	}
	return TUMBLR_BASE_URL + t.BlogName
}

func (t Tumblr) IsConsumerKey() bool {
	if t.ConsumerKey != "" {
		errors.New("Consumer Key is missing")
	}
	return true
}

func (t Tumblr) info() *http.Response {
	t.IsConsumerKey()
	u, err := url.Parse(t.ApiURL() + "/info")
	if err != nil {
		errors.New(err.Error())
	}
	q := u.Query()
	q.Set("api_key", t.ConsumerKey)
	u.RawQuery = q.Encode()
	resp, err := t.client.Get(u.String())
	check_error(err)
	return resp
}
