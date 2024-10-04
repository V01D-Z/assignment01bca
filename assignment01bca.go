package assignment01bca

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

type Block struct {
	Transaction  string
	Nonce        int
	PreviousHash string
	Hash         string
}

func NewBlock(transaction string, nonce int, previousHash string) *Block {
	block := &Block{
		Transaction:  transaction,
		Nonce:        nonce,
		PreviousHash: previousHash,
	}
	block.Hash = CalculateHash(block.Transaction + strconv.Itoa(block.Nonce) + block.PreviousHash)
	return block
}

var Blockchain []Block

func ListBlocks() {
	for i, block := range Blockchain {
		fmt.Printf("Block %d:\n", i)
		fmt.Printf("Transaction: %s\n", block.Transaction)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Previous Hash: %s\n", block.PreviousHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Println(strings.Repeat("-", 50))
	}
}

func ChangeBlock(index int, newTransaction string) {
	if index < len(Blockchain) {
		Blockchain[index].Transaction = newTransaction
		Blockchain[index].Hash = CalculateHash(Blockchain[index].Transaction + strconv.Itoa(Blockchain[index].Nonce) + Blockchain[index].PreviousHash)
	}
}

func VerifyChain() bool {
	for i := 1; i < len(Blockchain); i++ {
		previousBlock := Blockchain[i-1]
		currentBlock := Blockchain[i]
		if currentBlock.Hash != CalculateHash(currentBlock.Transaction+strconv.Itoa(currentBlock.Nonce)+currentBlock.PreviousHash) {
			fmt.Println("Blockchain has been tampered!")
			return false
		}
		if currentBlock.PreviousHash != previousBlock.Hash {
			fmt.Println("Previous block's hash doesn't match!")
			return false
		}
	}
	fmt.Println("Blockchain is valid.")
	return true
}

func CalculateHash(stringToHash string) string {
	hash := sha256.New()
	hash.Write([]byte(stringToHash))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}
