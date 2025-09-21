package main

import (
	"encoding/json"
	"github.com/phper95/tinydocker/enum"
	"github.com/phper95/tinydocker/network"
	"github.com/phper95/tinydocker/pkg/db"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := db.InitBoltDBClient(db.DefaultBoltDBClientName, enum.DefaultNetworkDBPath)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer func() {
		err := db.GetBoltDBClient(db.DefaultBoltDBClientName).Close()
		if err != nil {
			log.Println(err)
		}
	}()

	networks, err := db.GetBoltDBClient(db.DefaultBoltDBClientName).GetAll(enum.DefaultNetworkTable)
	if err != nil {
		log.Println(err)
		return
	}
	for name, data := range networks {
		log.Println("Network Name:  ", name, "Data: ", string(data))
		nw := &network.Network{}
		err := json.Unmarshal(data, nw)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("Network: %+v \n ", *nw)
	}
	allocatedIP, err := db.GetBoltDBClient(db.DefaultBoltDBClientName).GetAll(enum.AllocatedIPKey)
	if err != nil {
		log.Println(err)
		return
	}
	for name, data := range allocatedIP {
		log.Println("AllocatedIP: name", name, " data:", string(data))
	}
}
