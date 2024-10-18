package main

import (
	"fmt"

	"github.com/V01D-Z/assignment01bca/assignment01bca"
)

func main() {
	// Initialize the blockchain with the genesis block
	blockchain := assignment01bca.InitializeBlockchain()

	// Add transactions to the transaction pool
	blockchain.AddTransaction("Alice", "Bob", 1.0)
	blockchain.AddTransaction("Charlie", "Dave", 2.5)

	// Mine the first block with the current transactions
	fmt.Println("Mining Block 1...")
	blockchain.MineBlock()

	// Add more transactions to the transaction pool
	blockchain.AddTransaction("Eve", "Frank", 0.75)
	blockchain.AddTransaction("George", "Harry", 3.1)

	// Mine the second block
	fmt.Println("Mining Block 2...")
	blockchain.MineBlock()

	// List all blocks in the blockchain
	fmt.Println("\nListing All Blocks:")
	blockchain.ListBlocks()

	// Change a transaction in block 1 and verify the chain
	fmt.Println("\nTampering with the blockchain...")
	blockchain.Chain[1].Transactions[0].RecipientBlockchainAddress = "TamperedRecipient"
	blockchain.ListBlocks()
	valid := blockchain.VerifyChain()
	fmt.Println("Is Blockchain valid after tampering?", valid)
}
