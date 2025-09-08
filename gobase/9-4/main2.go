package main

import (
	"log"

	"github.com/phper95/tinydocker/pkg/db"
)

const TableName = "mytable"

func init() {
	err := db.InitBoltDBClients(db.DefaultBoltDBClientName, "mydb.db")
	if err != nil {
		panic(err)
	}
}
func main() {

	err := db.GetBoltDBClient(db.DefaultBoltDBClientName).Put(TableName, "key1", []byte("value1"))
	if err != nil {
		panic(err)
	}
	val, err := db.GetBoltDBClient(db.DefaultBoltDBClientName).Get(TableName, "key1")
	if err != nil {
		panic(err)
	}
	log.Println(string(val))
}
