package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/natiiix/gochain/proto_types"
)

const chainFile = "chain.dat"

func main() {
	chain := &proto_types.Blockchain{}
	err := unmarshal(chainFile, chain)
	if err != nil {
		chain = createBlockchain()
	}

	if !validateChain(chain) {
		log.Fatalln("Invalid initial chain!")
	}

	b1 := createBlock("first", chain)
	fmt.Println(blockToString(b1))

	b2 := createBlock("second", chain)
	fmt.Println(blockToString(b2))

	b3 := createBlock("third", chain)
	fmt.Println(blockToString(b3))

	fmt.Println("----------------------------------------------------------------")

	for i, b := range chain.Blocks {
		fmt.Printf("[%v] %v\n", i, blockToString(b))
	}

	if !validateChain(chain) {
		log.Fatalln("Invalid final chain!")
	}

	marshal(chainFile, chain)
}

func createBlockchain() *proto_types.Blockchain {
	chain := &proto_types.Blockchain{
		Blocks: []*proto_types.Block{},
	}

	createBlock("Genesis block", chain)
	return chain
}

func createBlockHashless(data string, prevHash []byte) *proto_types.Block {
	return &proto_types.Block{
		Timestamp:     time.Now().Unix(),
		Data:          data,
		PrevBlockHash: prevHash,
	}
}

func createBlock(data string, chain *proto_types.Blockchain) *proto_types.Block {
	prevHash := []byte{}
	if len(chain.Blocks) > 0 {
		prevHash = chain.Blocks[len(chain.Blocks)-1].Hash
	}

	b := createBlockHashless(data, prevHash)
	runWork(b)
	chain.Blocks = append(chain.Blocks, b)
	return b
}

func runWork(block *proto_types.Block) {
	for block.Nonce = 0; ; block.Nonce++ {
		valid, hash := validateBlock(block)

		if valid {
			block.Hash = hash[:]
			return
		}

		if block.Nonce == math.MaxInt64 {
			log.Fatalln("Unable to find nonce for block:", block)
		}
	}
}

func validateChain(chain *proto_types.Blockchain) bool {
	for _, block := range chain.Blocks {
		if valid, _ := validateBlock(block); !valid {
			return false
		}
	}

	return true
}

func validateBlock(block *proto_types.Block) (bool, [32]byte) {
	hash := getBlockHash(block)
	return hashCondition(hash), hash
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

func blockToString(block *proto_types.Block) string {
	return fmt.Sprintf("{ Timestamp: %v, Data: \"%v\", PrevBlockHash: %v, Nonce: %v, Hash: %v }",
		time.Unix(block.Timestamp, 0),
		block.Data,
		hex.EncodeToString(block.PrevBlockHash),
		block.Nonce, hex.EncodeToString(block.Hash))
}

func int64ToBytes(value int64) []byte {
	arr := make([]byte, 8)
	binary.LittleEndian.PutUint64(arr, uint64(value))
	return arr
}

func hashCondition(hash [32]byte) bool {
	return hash[0] == 0 && hash[1] == 0
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

func unmarshal(file string, msg proto.Message) error {
	in, err := ioutil.ReadFile(file)

	if err != nil {
		// log.Fatalln("Error reading file:", err)
		return err
	}

	if err := proto.Unmarshal(in, msg); err != nil {
		log.Fatalln("Failed to parse address book:", err)
		// return err
	}

	return nil
}
