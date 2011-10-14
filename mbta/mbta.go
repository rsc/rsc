package main

import (
	"io/ioutil"
	"log"
	"os"
	
	"goprotobuf.googlecode.com/hg/proto"
	"rsc.googlecode.com/hg/gtfs"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) != 2 {
		log.Fatal("usage: mbta file.pb")
	}
	pb, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	
	var feed gtfs.FeedMessage
	if err := proto.Unmarshal(pb, &feed); err != nil {
		log.Fatal(err)
	}
	
	proto.MarshalText(os.Stdout, &feed)
}
