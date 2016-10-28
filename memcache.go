package main


import (
	"github.com/bradfitz/gomemcache/memcache"
	"fmt"
)

func main() {
	mC := memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11211", "10.0.0.4:11211")
	mC.Set(&memcache.Item{Key: "data", Value: []byte("30") })
	dataItem, _ := mC.Get("data")
	fmt.Println("Data",dataItem)
}