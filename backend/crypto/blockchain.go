package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// HashBlock creates a SHA-256 hash of block data
func HashBlock(index int64, previousHash string, timestamp time.Time, merkleRoot string, nonce int64, difficulty int) string {
	data := fmt.Sprintf("%d%s%d%s%d%d", index, previousHash, timestamp.Unix(), merkleRoot, nonce, difficulty)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// CalculateMerkleRoot calculates the Merkle root from transaction IDs
func CalculateMerkleRoot(transactionIDs []string) string {
	if len(transactionIDs) == 0 {
		// Empty block - hash of empty string
		hash := sha256.Sum256([]byte(""))
		return hex.EncodeToString(hash[:])
	}

	if len(transactionIDs) == 1 {
		hash := sha256.Sum256([]byte(transactionIDs[0]))
		return hex.EncodeToString(hash[:])
	}

	// Build Merkle tree
	hashes := make([]string, len(transactionIDs))
	for i, txID := range transactionIDs {
		hash := sha256.Sum256([]byte(txID))
		hashes[i] = hex.EncodeToString(hash[:])
	}

	// Keep hashing pairs until we have one root
	for len(hashes) > 1 {
		var newLevel []string

		for i := 0; i < len(hashes); i += 2 {
			var combined string
			if i+1 < len(hashes) {
				combined = hashes[i] + hashes[i+1]
			} else {
				// Odd number - duplicate the last hash
				combined = hashes[i] + hashes[i]
			}
			hash := sha256.Sum256([]byte(combined))
			newLevel = append(newLevel, hex.EncodeToString(hash[:]))
		}

		hashes = newLevel
	}

	return hashes[0]
}

// ValidateBlockHash checks if a hash meets the difficulty requirement
func ValidateBlockHash(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

// MineBlock attempts to find a valid nonce for the block
// Returns the nonce and hash if successful
func MineBlock(index int64, previousHash string, timestamp time.Time, merkleRoot string, difficulty int, maxIterations int64) (int64, string, bool) {
	prefix := strings.Repeat("0", difficulty)

	for nonce := int64(0); nonce < maxIterations; nonce++ {
		hash := HashBlock(index, previousHash, timestamp, merkleRoot, nonce, difficulty)
		if strings.HasPrefix(hash, prefix) {
			return nonce, hash, true
		}
	}

	return 0, "", false
}

// GetGenesisBlockHash returns the hash for the genesis block
func GetGenesisBlockHash() string {
	genesisData := "Genesis Block - Crypto Wallet Blockchain - 2025"
	hash := sha256.Sum256([]byte(genesisData))
	return hex.EncodeToString(hash[:])
}

// CalculateDifficulty adjusts difficulty based on block times
func CalculateDifficulty(currentDifficulty int, lastBlockTimes []time.Duration, targetBlockTime time.Duration) int {
	if len(lastBlockTimes) < 2 {
		return currentDifficulty
	}

	// Calculate average block time
	var totalTime time.Duration
	for _, t := range lastBlockTimes {
		totalTime += t
	}
	avgTime := totalTime / time.Duration(len(lastBlockTimes))

	// Adjust difficulty
	if avgTime < targetBlockTime/2 {
		// Blocks are coming too fast, increase difficulty
		return currentDifficulty + 1
	} else if avgTime > targetBlockTime*2 {
		// Blocks are too slow, decrease difficulty
		if currentDifficulty > 1 {
			return currentDifficulty - 1
		}
	}

	return currentDifficulty
}

// ValidateChainLink validates that two blocks are properly linked
func ValidateChainLink(previousBlock, currentBlock struct {
	Index        int64
	Hash         string
	PreviousHash string
}) bool {
	// Check index is sequential
	if currentBlock.Index != previousBlock.Index+1 {
		return false
	}

	// Check hash link
	if currentBlock.PreviousHash != previousBlock.Hash {
		return false
	}

	return true
}
