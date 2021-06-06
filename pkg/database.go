package pkg

import (
	"fmt"
	"log"

	"github.com/xujiajun/nutsdb"
)

var database *nutsdb.DB

func OpenDB() {

	opt := nutsdb.DefaultOptions
	fmt.Println(opt)

	opt.Dir = "./data/nutsdb"
	db, err := nutsdb.Open(opt)
	if err != nil {
		log.Fatal(err)
	}

	database = db
}

func viewAllData() {
	if err := database.View(
		func(tx *nutsdb.Tx) error {
			bucket := "bucket"
			entries, err := tx.GetAll(bucket)
			if err != nil {
				return err
			}

			for _, entry := range entries {
				fmt.Println(string(entry.Key))
			}

			return nil
		}); err != nil {
		log.Println(err)
	}
}

func CloseDB() {
	defer database.Close()
}

func InsertData(key []byte, value []byte) {
	if err := database.Update(
		func(tx *nutsdb.Tx) error {
			if err := tx.Put("bucket", key, value, 50); err != nil {
				return err
			}
			return nil
		}); err != nil {
		log.Fatal(err)
	}
}
