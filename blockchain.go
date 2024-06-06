package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type BlockData string

type Block struct {
	Hash         string
	PreviousHash string
	Timestamp    time.Time
	Pow          int
	Data         BlockData
}

func (b *Block) calculateHash() string {
	data, err := json.Marshal(b.Data)
	if err != nil {
		log.Fatalf("Block data corrupted")
		panic(1)
	}

	hashString := b.PreviousHash + string(data) + b.Timestamp.String() + strconv.Itoa(b.Pow)
	// hash := sha256.Sum256([]byte(hashString))
	hash := fmt.Sprintf("%x", sha512.Sum512([]byte(hashString)))
	// log.Println(hash)
	return hash
}

func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.Hash, strings.Repeat("0", difficulty)) {
		b.Pow++
		b.Hash = b.calculateHash()
	}
}

func (b *Block) isValid() bool {
	return b.Hash != b.calculateHash()
}

type Blockchain struct {
	GenesisBlock Block
	Chain        []Block
	Difficulty   int
}

func (b *Blockchain) mineBlock(blockData BlockData) {
	lastBlock := b.Chain[len(b.Chain)-1]
	newBlock := Block{
		Data:         blockData,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now(),
	}
	newBlock.mine(b.Difficulty)
	b.Chain = append(b.Chain, newBlock)
}

func (b *Blockchain) isValid() bool {
	for i := range b.Chain[1:] {
		prevBlock := b.Chain[i]
		block := b.Chain[i+1]
		if block.PreviousHash != prevBlock.Hash || block.isValid() {
			return false
		}
	}

	return true
}

func CreateBlockChain(difficulty int) Blockchain {
	genesisBlock := Block{
		Hash:      "0",
		Timestamp: time.Now(),
	}

	return Blockchain{
		GenesisBlock: genesisBlock,
		Chain:        []Block{genesisBlock},
		Difficulty:   difficulty,
	}
}

func main() {
	blockChain := CreateBlockChain(4)

	block1, _ := json.Marshal(map[string]interface{}{
		"from":   "Savy",
		"to":     "josy",
		"amount": 100,
		"prof":   "sdfds934msdcsdcd",
	})

	block2, _ := json.Marshal(map[string]interface{}{
		"from":   "day",
		"to":     "josy",
		"amount": 100,
		"prof":   "sdfds934msdcsdcd",
	})

	blockChain.mineBlock(BlockData(block1))

	blockChain.mineBlock(BlockData(block2))

	// blockChain.Chain[1].Data = "tempered data"

	if !blockChain.isValid() {
		log.Fatalln("blockchain not valid")
		panic(1)
	}

	log.Println(blockChain.Chain[0])
	log.Println(blockChain.Chain[1])
	log.Println(blockChain.Chain[2])
}
