package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Block represents a block in the blockchain
type Block struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Index            int64              `json:"index" bson:"index"`                       // Block height/number
	Hash             string             `json:"hash" bson:"hash"`                         // Hash of this block
	PreviousHash     string             `json:"previousHash" bson:"previousHash"`         // Hash of previous block
	Timestamp        time.Time          `json:"timestamp" bson:"timestamp"`               // When block was mined
	Transactions     []Transaction      `json:"transactions" bson:"transactions"`         // Transactions in this block
	TransactionCount int                `json:"transactionCount" bson:"transactionCount"` // Number of transactions
	MerkleRoot       string             `json:"merkleRoot" bson:"merkleRoot"`             // Merkle root of transactions
	Nonce            int64              `json:"nonce" bson:"nonce"`                       // Proof-of-work nonce
	Difficulty       int                `json:"difficulty" bson:"difficulty"`             // Mining difficulty (number of leading zeros)
	MinerWalletID    string             `json:"minerWalletId" bson:"minerWalletId"`       // Wallet that mined this block
	MiningReward     float64            `json:"miningReward" bson:"miningReward"`         // Reward for mining this block
	Size             int64              `json:"size" bson:"size"`                         // Block size in bytes (approximate)
}

// BlockHeader contains just the header info for lighter queries
type BlockHeader struct {
	Index        int64     `json:"index"`
	Hash         string    `json:"hash"`
	PreviousHash string    `json:"previousHash"`
	Timestamp    time.Time `json:"timestamp"`
	MerkleRoot   string    `json:"merkleRoot"`
	Nonce        int64     `json:"nonce"`
	Difficulty   int       `json:"difficulty"`
}

// GenesisBlock creates the first block in the chain
type GenesisBlockInfo struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// MiningJob represents a job for mining a new block
type MiningJob struct {
	BlockIndex       int64         `json:"blockIndex"`
	PreviousHash     string        `json:"previousHash"`
	Transactions     []Transaction `json:"transactions"`
	MerkleRoot       string        `json:"merkleRoot"`
	Difficulty       int           `json:"difficulty"`
	MiningReward     float64       `json:"miningReward"`
	TargetPrefix     string        `json:"targetPrefix"` // The hash must start with this
}

// MiningResult is returned when a block is successfully mined
type MiningResult struct {
	Success bool   `json:"success"`
	Block   *Block `json:"block,omitempty"`
	Message string `json:"message"`
	Nonce   int64  `json:"nonce,omitempty"`
	Hash    string `json:"hash,omitempty"`
}

// BlockchainStats provides statistics about the blockchain
type BlockchainStats struct {
	TotalBlocks         int64   `json:"totalBlocks"`
	TotalTransactions   int64   `json:"totalTransactions"`
	CurrentDifficulty   int     `json:"currentDifficulty"`
	LastBlockHash       string  `json:"lastBlockHash"`
	LastBlockTime       string  `json:"lastBlockTime"`
	AverageBlockTime    float64 `json:"averageBlockTime"`    // in seconds
	TotalMiningRewards  float64 `json:"totalMiningRewards"`
	PendingTransactions int     `json:"pendingTransactions"`
}

// ChainValidation represents the result of validating the blockchain
type ChainValidation struct {
	IsValid       bool     `json:"isValid"`
	BlocksChecked int64    `json:"blocksChecked"`
	Errors        []string `json:"errors,omitempty"`
}

// BlockchainConfig holds configuration for the blockchain
type BlockchainConfig struct {
	InitialDifficulty     int     `json:"initialDifficulty"`     // Starting difficulty
	BlockReward           float64 `json:"blockReward"`           // Mining reward
	DifficultyAdjustment  int     `json:"difficultyAdjustment"`  // Adjust every N blocks
	TargetBlockTime       int     `json:"targetBlockTime"`       // Target time in seconds
	MaxTransactionsPerBlock int   `json:"maxTransactionsPerBlock"`
}

// Default blockchain configuration
var DefaultBlockchainConfig = BlockchainConfig{
	InitialDifficulty:       4,     // 4 leading zeros
	BlockReward:             50.0,  // 50 coins per block
	DifficultyAdjustment:    10,    // Adjust every 10 blocks
	TargetBlockTime:         30,    // 30 seconds target
	MaxTransactionsPerBlock: 100,   // Max 100 transactions per block
}
