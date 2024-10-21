package assignment01bca

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	TransactionID              string
	SenderBlockchainAddress    string
	RecipientBlockchainAddress string
	Value                      float32
}

type Block struct {
	Index        int
	Timestamp    time.Time
	Transactions []*Transaction
	Nonce        int
	PreviousHash string
	Hash         string
}

type Blockchain struct {
	Chain           []*Block
	TransactionPool []*Transaction
}

// Create a new block
func NewBlock(nonce int, previousHash string, transactions []*Transaction) *Block {
	block := &Block{
		Timestamp:    time.Now(),
		Transactions: transactions,
		Nonce:        nonce,
		PreviousHash: previousHash,
	}
	block.Hash = CalculateBlockHash(block)
	return block
}

// AddTransaction creates a new transaction and adds it to the pool
func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.TransactionPool = append(bc.TransactionPool, t)
}

// Create a new transaction
func NewTransaction(sender string, recipient string, value float32) *Transaction {
	transaction := &Transaction{
		SenderBlockchainAddress:    sender,
		RecipientBlockchainAddress: recipient,
		Value:                      value,
	}
	transaction.TransactionID = CalculateTransactionHash(transaction)
	return transaction
}

// Calculate transaction hash
func CalculateTransactionHash(transaction *Transaction) string {
	data := transaction.SenderBlockchainAddress + transaction.RecipientBlockchainAddress + strconv.FormatFloat(float64(transaction.Value), 'f', 2, 32)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Calculate block hash (includes transaction hashes)
func CalculateBlockHash(block *Block) string {
	transactionsHash := ""
	for _, transaction := range block.Transactions {
		transactionsHash += CalculateTransactionHash(transaction)
	}

	data := strconv.Itoa(block.Index) + block.Timestamp.String() + strconv.Itoa(block.Nonce) + block.PreviousHash + transactionsHash
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// List all blocks and their transactions
func (bc *Blockchain) ListBlocks() {
	for i, block := range bc.Chain {
		fmt.Printf("Block %d:\n", i)
		fmt.Printf("Timestamp: %s\n", block.Timestamp)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Previous Hash: %s\n", block.PreviousHash)
		fmt.Printf("Hash: %s\n", block.Hash)

		// Display transactions as JSON
		txs, _ := json.MarshalIndent(block.Transactions, "", "  ")
		fmt.Printf("Transactions: %s\n", string(txs))
	}
}

// Display a specific block
func (bc *Blockchain) DisplayBlock(index int) {
	if index < len(bc.Chain) {
		block := bc.Chain[index]
		fmt.Printf("Block %d: Hash: %s, Transactions: %+v\n", block.Index, block.Hash, block.Transactions)
	}
}

// Tamper with a specific transaction
func (bc *Blockchain) Tamper(blockIndex int, txIndex int, newRecipient string) {
	if blockIndex < len(bc.Chain) && txIndex < len(bc.Chain[blockIndex].Transactions) {
		bc.Chain[blockIndex].Transactions[txIndex].RecipientBlockchainAddress = newRecipient
		fmt.Printf("Block %d, Transaction %d has been tampered with!\n", blockIndex, txIndex)
	}
}

// Verify if the blockchain is valid
func (bc *Blockchain) VerifyChain() bool {
	for i, block := range bc.Chain {
		recalculatedHash := CalculateBlockHash(block)
		if block.Hash != recalculatedHash {
			fmt.Printf("Block %d is invalid! Stored Hash: %s, Recalculated Hash: %s\n", i, block.Hash, recalculatedHash)
			return false
		}
		if i > 0 && block.PreviousHash != bc.Chain[i-1].Hash {
			fmt.Printf("Block %d has an invalid previous hash!\n", i)
			return false
		}
	}
	return true
}

// Proof of work function
func ProofOfWork(previousHash string, transactions []*Transaction, difficulty int) int {
	nonce := 0
	target := strings.Repeat("0", difficulty)
	var hash string
	for {
		block := Block{Nonce: nonce, PreviousHash: previousHash, Transactions: transactions}
		hash = CalculateBlockHash(&block)
		if strings.HasPrefix(hash, target) {
			break
		}
		nonce++
	}
	return nonce
}

// Mine a new block
func (bc *Blockchain) MineBlock() {
	previousBlock := bc.Chain[len(bc.Chain)-1]
	nonce := ProofOfWork(previousBlock.Hash, bc.TransactionPool, 2)
	newBlock := NewBlock(nonce, previousBlock.Hash, bc.TransactionPool)
	bc.Chain = append(bc.Chain, newBlock)
	bc.TransactionPool = []*Transaction{} // Empty the transaction pool
}

// Initialize blockchain with the genesis block
func InitializeBlockchain() *Blockchain {
	genesisBlock := NewBlock(0, "0", []*Transaction{})
	return &Blockchain{
		Chain: []*Block{genesisBlock},
	}
}
