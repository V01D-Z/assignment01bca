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
func NewBlock(nonce int, previousHash, hash string, transactions []*Transaction) *Block {
	block := &Block{
		Timestamp:    time.Now(),
		Transactions: transactions,
		Nonce:        nonce,
		PreviousHash: previousHash,
		Hash:         hash, // Set the hash directly from ProofOfWork
	}
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

// Calculate block hash (includes nonce and transaction hashes)
func CalculateBlockHash(nonce int, previousHash string, transactions []*Transaction) string {
	transactionsHash := ""
	for _, transaction := range transactions {
		transactionsHash += CalculateTransactionHash(transaction)
	}

	// Include the nonce in the block hash calculation
	data := previousHash + transactionsHash + strconv.Itoa(nonce)
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
		fmt.Println(strings.Repeat("-", 50))
	}
}

func (bc *Blockchain) DisplayBlock(index int) {
	// Check if the index is valid
	if index < 0 || index >= len(bc.Chain) {
		fmt.Println("Invalid block index.")
		return
	}

	// Get the block at the specified index
	block := bc.Chain[index]

	// Display the block details
	fmt.Printf("Block %d:\n", index)
	fmt.Printf("Timestamp: %s\n", block.Timestamp)
	fmt.Printf("Nonce: %d\n", block.Nonce)
	fmt.Printf("Previous Hash: %s\n", block.PreviousHash)
	fmt.Printf("Hash: %s\n", block.Hash)

	// Display transactions as JSON
	txs, _ := json.MarshalIndent(block.Transactions, "", "  ")
	fmt.Printf("Transactions: %s\n", string(txs))
	fmt.Println(strings.Repeat("-", 50))
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
		recalculatedHash := CalculateBlockHash(block.Nonce, block.PreviousHash, block.Transactions)
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

// Proof of work function (returns nonce and hash)
func ProofOfWork(previousHash string, transactions []*Transaction, difficulty int) (int, string) {
	nonce := 0
	target := strings.Repeat("0", difficulty)
	var hash string
	for {
		hash = CalculateBlockHash(nonce, previousHash, transactions)
		if strings.HasPrefix(hash, target) {
			break
		}
		nonce++
	}
	return nonce, hash
}

// Mine a new block
func (bc *Blockchain) MineBlock() {
	previousBlock := bc.Chain[len(bc.Chain)-1]
	nonce, hash := ProofOfWork(previousBlock.Hash, bc.TransactionPool, 2)
	newBlock := NewBlock(nonce, previousBlock.Hash, hash, bc.TransactionPool)
	bc.Chain = append(bc.Chain, newBlock)
	bc.TransactionPool = []*Transaction{} // Empty the transaction pool
}

// Initialize blockchain with the genesis block
func InitializeBlockchain() *Blockchain {
	genesisBlock := NewBlock(0, "0", "genesis_hash", []*Transaction{})
	return &Blockchain{
		Chain: []*Block{genesisBlock},
	}
}
