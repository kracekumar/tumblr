package tumblr

import (
	"fmt"
	"testing"
)

func TestTumblrCallbackURLFail(t *testing.T) {
	tum := &Tumblr{}
	tum.CallbackURL = "http://tumblr.com"
	tum.ConsumerKey = "consumerkey"
	tum.ConsumerSecret = "consumersecret"
	fmt.Println(tum)
	//t.Log(Tumblr{})
	//tumblr := Tumblr{ConsumerSecret: "xxxx", ConsumerKey: "yyyyy", CallbackURL: "http"}
	//t.Log(tumblr)
}

func TestTumblrBlog(t *testing.T) {
	tum := &Tumblr{}
	tum.BlogName = "kracekumar.com"
	tum.ConsumerKey = CONSUMERKEY //comes from settings.go
	rsp := tum.info()
	fmt.Println(FetchResponseBody(rsp))
	//fmt.Println("test")
	//fmt.Println(rsp)
}
