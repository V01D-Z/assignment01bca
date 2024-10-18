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

// Calculate the hash for a block
func CalculateBlockHash(block *Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp.String() + strconv.Itoa(block.Nonce) + block.PreviousHash
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Calculate the hash for a transaction
func CalculateTransactionHash(transaction *Transaction) string {
	record := transaction.SenderBlockchainAddress + transaction.RecipientBlockchainAddress + fmt.Sprintf("%f", transaction.Value)
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

// VerifyChain checks if the blockchain is still valid
func (bc *Blockchain) VerifyChain() bool {
	for i := 1; i < len(bc.Chain); i++ {
		previousBlock := bc.Chain[i-1]
		currentBlock := bc.Chain[i]

		// Recalculate the block's hash to see if it matches
		if currentBlock.Hash != CalculateBlockHash(currentBlock) {
			fmt.Println("Blockchain has been tampered! Block hash mismatch detected.")
			return false
		}

		// Check if the current block's PreviousHash matches the previous block's hash
		if currentBlock.PreviousHash != previousBlock.Hash {
			fmt.Println("Blockchain has been tampered! Previous block's hash doesn't match.")
			return false
		}

		// Verify each transaction's integrity in the current block
		for _, tx := range currentBlock.Transactions {
			calculatedTxHash := CalculateTransactionHash(tx)
			if tx.TransactionID != calculatedTxHash {
				fmt.Println("Blockchain has been tampered! Transaction ID mismatch detected.")
				return false
			}
		}
	}

	fmt.Println("Blockchain is valid.")
	return true
}

// Proof of Work to derive the correct nonce
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

// Method to mine a new block, derive nonce, and add block to chain
func (bc *Blockchain) MineBlock() {
	previousBlock := bc.Chain[len(bc.Chain)-1]
	nonce := ProofOfWork(previousBlock.Hash, bc.TransactionPool, 2)
	newBlock := NewBlock(nonce, previousBlock.Hash, bc.TransactionPool)
	bc.Chain = append(bc.Chain, newBlock)
	// Empty the transaction pool after mining
	bc.TransactionPool = []*Transaction{}
}

// Initialize blockchain with the genesis block
func InitializeBlockchain() *Blockchain {
	genesisBlock := NewBlock(0, "0", []*Transaction{})
	return &Blockchain{
		Chain: []*Block{genesisBlock},
	}
}
