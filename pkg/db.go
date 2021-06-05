package pkg

import (
	"fmt"
	"log"

	"github.com/xujiajun/nutsdb"
)


func OpenDB(){

	opt := nutsdb.DefaultOptions
	fmt.Println(opt)

	opt.RWMode = nutsdb.MMap;
	fmt.Println(opt)
	
	opt.Dir = "./data/nutsdb"
	db, err := nutsdb.Open(opt)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}