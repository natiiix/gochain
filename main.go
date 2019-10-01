package main

import (
	"encoding/binary"
	"io/ioutil"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/natiiix/gochain/proto_types"
)

func main() {
	// TODO
}

func createBlockchain() *proto_types.Blockchain {
	return &proto_types.Blockchain{
		Blocks: []*proto_types.Block{
			createBlockHashless("Genesis block", []byte{}),
		},
	}
}

func createBlockHashless(data string, prevHash []byte) *proto_types.Block {
	return &proto_types.Block{
		Timestamp:     time.Now().Unix(),
		Data:          data,
		PrevBlockHash: prevHash,
	}
}

func int64ToBytes(value int64) []byte {
	arr := make([]byte, 8)
	binary.LittleEndian.PutUint64(arr, uint64(value))
	return arr
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
