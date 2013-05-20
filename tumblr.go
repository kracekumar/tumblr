package tumblr

import (
	"encoding/json"
	"errors"
	"github.com/kurrik/oauth1a"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
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

type BlogMeta struct {
	Status int
	Msg    string
}

type TumblrBlog struct {
	Title       string
	Posts       int
	Name        string
	Url         string
	Updated     uint64
	Description string
	Ask         bool
	ASkAnon     bool
	Likes       int
}

type Info struct {
	Meta     BlogMeta
	Response struct {
		Blog TumblrBlog
	}
}

type TumblrPost struct {
	BlogName  string
	Id        uint64
	PostUrl   string
	Type      string
	Date      string
	Timestamp uint64
	State     string
	Format    string
	ReblogKey string
	Tags      []string
	NoteCount int
	Title     string
	Body      string
}

type TumblrTextPosts struct {
	Meta     BlogMeta
	Response struct {
		Blog       TumblrBlog
		Posts      []TumblrPost
		TotalPosts int
	}
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

func NewTumblr() *Tumblr {
	t := &Tumblr{}
	t.client = &http.Client{}
	return t
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

func (t Tumblr) info() *Info {
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
	var blog_info Info
	err = json.Unmarshal(FetchResponseBody(resp), &blog_info)
	check_error(err)
	return &blog_info
}

func (t Tumblr) posts(post_type string, limit, offset int, id int, reblog_info, notes_info bool, filter string) *TumblrTextPosts {
	t.IsConsumerKey()
	if post_type != "text" {
		errors.New("Not supported post type in this version")
	}

	u, err := url.Parse(t.ApiURL() + "/posts/") // + post_type)
	if err != nil {
		errors.New(err.Error())
	}
	q := u.Query()
	q.Set("api_key", t.ConsumerKey)
	q.Set("reblog_info", strconv.FormatBool(reblog_info))
	q.Set("notes_info", strconv.FormatBool(notes_info))
	if filter != "HTML" {
		q.Set("filter", filter)
	}
	// since go is not dynamic like Python, using negative to indicate not to set
	if id > -1 {
		q.Set("id", strconv.Itoa(id))
	}
	if limit > -1 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset > -1 {
		q.Set("offset", strconv.Itoa(offset))
	}
	q.Set("type", post_type)
	u.RawQuery = q.Encode()
	resp, err := t.client.Get(u.String())
	log.Println(u.String())
	check_error(err)
	//TumblrPosts
	var t_posts TumblrTextPosts
	if resp.StatusCode != 200 {
		log.Fatal("Response code is " + resp.Status)
		errors.New("Response code is " + resp.Status)
	}
	err = json.Unmarshal(FetchResponseBody(resp), &t_posts)
	check_error(err)
	return &t_posts

}

func (t Tumblr) GetAllTextPosts(post_type, filter string, reblog_info, notes_info bool, limit int) *TumblrTextPosts {
	/*
		param post_type: string
		["text", "quote", "link", "answer", "video", "audio", "photo", "cat"]

		param filter: string
		["HTML", "text", "raw"]
	*/
	if post_type != "text" {
		errors.New("As of Now only text post type is supported")
	}
	count := 0
	filter_types := []string{"HTML", "text", "raw"}
	for _, val := range filter_types {
		if val == filter {
			count += 1
		}
	}
	if count == 0 {
		errors.New("filter should be `HTML` or `text` or `raw` ")
	}
	//posts(post_type string, limit, offset int, id int, reblog_info, notes_info bool, filter string) *TumblrTextPosts {
	/*var *text_posts TumblrTextPosts
	text_posts = t.posts("text", -1, -1, -1, reblog_info, notes_info, filter)
	return &text_posts*/
	return t.posts("text", -1, 1, -1, reblog_info, notes_info, filter)
}
