package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io/ioutil"
	"log"
	"math"
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

func createBlock(data string, chain *proto_types.Blockchain) *proto_types.Block {
	b := createBlockHashless(data, chain.Blocks[len(chain.Blocks)-1].Hash)
	runWork(b)
	chain.Blocks = append(chain.Blocks, b)
	return b
}

func getBlockHash(block *proto_types.Block) [32]byte {
	blockBytes := bytes.Join([][]byte{
		int64ToBytes(block.Timestamp),
		[]byte(block.Data),
		block.PrevBlockHash,
		int64ToBytes(block.Nonce),
	}, []byte{})

	return sha256.Sum256(blockBytes)
}

func runWork(block *proto_types.Block) {
	block.Nonce = 0

	for {
		hash := getBlockHash(block)

		if hashCondition(hash) {
			block.Hash = hash[:]
			break
		}

		if block.Nonce == math.MaxInt64 {
			log.Fatalln("Unable to find nonce for block:", block)
		} else {
			block.Nonce++
		}
	}
}

func int64ToBytes(value int64) []byte {
	arr := make([]byte, 8)
	binary.LittleEndian.PutUint64(arr, uint64(value))
	return arr
}

func hashCondition(hash [32]byte) bool {
	return hash[0] == 0
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
