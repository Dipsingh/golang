package main


import (
	"fmt"
	"github.com/couchbase/go-couchbase"
)

func main(){
	conn, err := couchbase.Connect("http://localhost:8091")
	if err != nil {

	}

	pool,err := conn.GetPool("file_manager")
	if err != nil {

	}

	bucket,err := pool.GetBucket("file_manager")
	if err != nil {

	}

	fmt.Println("FMT",bucket)


}
