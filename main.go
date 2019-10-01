package main

import (
	"io/ioutil"
	"log"

	"github.com/golang/protobuf/proto"
)

func main() {
	// TODO
}

func marshal(file string, msg proto.Message) {
	out, err := proto.Marshal(msg)

	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	}

	if err := ioutil.WriteFile(file, out, 0644); err != nil {
		log.Fatalln("Failed to write address book:", err)
	}
}

func unmarshal(file string, msg proto.Message) {
	in, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatalln("Error reading file:", err)
	}

	if err := proto.Unmarshal(in, msg); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}
}
