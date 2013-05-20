package tumblr

import (
	"fmt"
	//"io/ioutil"
	"testing"
)

func TestTumblrBlog(t *testing.T) {
	tum := NewTumblr()
	tum.BlogName = "kracekumar.com"
	tum.ConsumerKey = CONSUMERKEY //comes from settings.go
	rsp := tum.info()
	fmt.Println(rsp)
}

func TestGetAllTextPosts(t *testing.T) {
	tum := NewTumblr()
	tum.BlogName = "kracekumar.com"
	tum.ConsumerKey = CONSUMERKEY
	rsp := tum.GetAllTextPosts("text", "raw", false, false, 52)
	fmt.Println(rsp)
}
