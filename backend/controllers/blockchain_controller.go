package controllers

import (
	"context"
	"crypto-wallet-backend/crypto"
	"crypto-wallet-backend/database"
	"crypto-wallet-backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getBlockCollection() *mongo.Collection {
	return database.GetCollection("blocks")
}

// GetBlockchainStats returns statistics about the blockchain
func GetBlockchainStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Count blocks
	totalBlocks, _ := getBlockCollection().CountDocuments(ctx, bson.M{})

	// Count all transactions in blocks
	var totalTransactions int64
	cursor, err := getBlockCollection().Find(ctx, bson.M{})
	if err == nil {
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var block models.Block
			if cursor.Decode(&block) == nil {
				totalTransactions += int64(block.TransactionCount)
			}
		}
	}

	// Get last block
	var lastBlock models.Block
	err = getBlockCollection().FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.M{"index": -1})).Decode(&lastBlock)
	
	lastBlockHash := ""
	lastBlockTime := ""
	currentDifficulty := models.DefaultBlockchainConfig.InitialDifficulty
	
	if err == nil {
		lastBlockHash = lastBlock.Hash
		lastBlockTime = lastBlock.Timestamp.Format(time.RFC3339)
		currentDifficulty = lastBlock.Difficulty
	}

	// Count pending transactions
	pendingCount, _ := getTransactionCollection().CountDocuments(ctx, bson.M{"status": models.TxStatusPending})

	// Calculate total mining rewards
	var totalRewards float64
	cursor2, err := getBlockCollection().Find(ctx, bson.M{})
	if err == nil {
		defer cursor2.Close(ctx)
		for cursor2.Next(ctx) {
			var block models.Block
			if cursor2.Decode(&block) == nil {
				totalRewards += block.MiningReward
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": models.BlockchainStats{
			TotalBlocks:         totalBlocks,
			TotalTransactions:   totalTransactions,
			CurrentDifficulty:   currentDifficulty,
			LastBlockHash:       lastBlockHash,
			LastBlockTime:       lastBlockTime,
			TotalMiningRewards:  totalRewards,
			PendingTransactions: int(pendingCount),
		},
	})
}

// GetBlocks returns a list of blocks with pagination
func GetBlocks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	skip := (page - 1) * limit

	// Get total count
	total, _ := getBlockCollection().CountDocuments(ctx, bson.M{})

	// Get blocks sorted by index descending (newest first)
	cursor, err := getBlockCollection().Find(ctx, bson.M{},
		options.Find().
			SetSort(bson.M{"index": -1}).
			SetSkip(int64(skip)).
			SetLimit(int64(limit)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blocks"})
		return
	}
	defer cursor.Close(ctx)

	var blocks []models.Block
	if err := cursor.All(ctx, &blocks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse blocks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"blocks": blocks,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// GetBlock returns a specific block by hash or index
func GetBlock(c *gin.Context) {
	identifier := c.Param("identifier")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var block models.Block
	var err error

	// Try to parse as index first
	if index, parseErr := strconv.ParseInt(identifier, 10, 64); parseErr == nil {
		err = getBlockCollection().FindOne(ctx, bson.M{"index": index}).Decode(&block)
	} else {
		// Try as hash
		err = getBlockCollection().FindOne(ctx, bson.M{"hash": identifier}).Decode(&block)
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Block not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"block": block})
}

// GetLatestBlock returns the most recent block
func GetLatestBlock(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var block models.Block
	err := getBlockCollection().FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.M{"index": -1})).Decode(&block)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "No blocks found. Genesis block not created yet."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"block": block})
}

// CreateGenesisBlock creates the first block in the chain
func CreateGenesisBlock(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if genesis already exists
	count, _ := getBlockCollection().CountDocuments(ctx, bson.M{"index": 0})
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Genesis block already exists"})
		return
	}

	// Create genesis block
	now := time.Now()
	genesisHash := crypto.GetGenesisBlockHash()

	genesis := models.Block{
		Index:            0,
		Hash:             genesisHash,
		PreviousHash:     "0000000000000000000000000000000000000000000000000000000000000000",
		Timestamp:        now,
		Transactions:     []models.Transaction{},
		TransactionCount: 0,
		MerkleRoot:       crypto.CalculateMerkleRoot([]string{}),
		Nonce:            0,
		Difficulty:       models.DefaultBlockchainConfig.InitialDifficulty,
		MinerWalletID:    "system",
		MiningReward:     0,
		Size:             0,
	}

	_, err := getBlockCollection().InsertOne(ctx, genesis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create genesis block"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Genesis block created successfully!",
		"block":   genesis,
	})
}

// MineBlock mines a new block with pending transactions
func MineBlock(c *gin.Context) {
	userID := c.GetString("userId")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // Longer timeout for mining
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get miner's wallet
	var minerWallet models.Wallet
	err = getWalletCollection().FindOne(ctx, bson.M{"userId": objID}).Decode(&minerWallet)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found. Please generate a wallet first."})
		return
	}

	// Get the latest block
	var lastBlock models.Block
	err = getBlockCollection().FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.M{"index": -1})).Decode(&lastBlock)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No genesis block. Please create genesis block first."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Get pending transactions
	cursor, err := getTransactionCollection().Find(ctx, bson.M{"status": models.TxStatusPending},
		options.Find().SetLimit(int64(models.DefaultBlockchainConfig.MaxTransactionsPerBlock)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending transactions"})
		return
	}
	defer cursor.Close(ctx)

	var pendingTxs []models.Transaction
	if err := cursor.All(ctx, &pendingTxs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse transactions"})
		return
	}

	// Build transaction ID list for merkle root
	txIDs := make([]string, len(pendingTxs))
	for i, tx := range pendingTxs {
		txIDs[i] = tx.TransactionID
	}

	// Calculate merkle root
	merkleRoot := crypto.CalculateMerkleRoot(txIDs)

	// Prepare new block
	newIndex := lastBlock.Index + 1
	timestamp := time.Now()
	difficulty := lastBlock.Difficulty

	// Adjust difficulty every N blocks
	if newIndex%int64(models.DefaultBlockchainConfig.DifficultyAdjustment) == 0 && newIndex > 0 {
		// Simple difficulty adjustment - could be improved
		difficulty = calculateNewDifficulty(ctx, difficulty)
	}

	// Mine the block (find valid nonce)
	maxIterations := int64(10000000) // 10 million attempts max
	nonce, hash, found := crypto.MineBlock(newIndex, lastBlock.Hash, timestamp, merkleRoot, difficulty, maxIterations)

	if !found {
		c.JSON(http.StatusRequestTimeout, gin.H{
			"error":   "Mining timeout - could not find valid hash",
			"message": "Try again or reduce difficulty",
		})
		return
	}

	// Create the block
	miningReward := models.DefaultBlockchainConfig.BlockReward
	newBlock := models.Block{
		Index:            newIndex,
		Hash:             hash,
		PreviousHash:     lastBlock.Hash,
		Timestamp:        timestamp,
		Transactions:     pendingTxs,
		TransactionCount: len(pendingTxs),
		MerkleRoot:       merkleRoot,
		Nonce:            nonce,
		Difficulty:       difficulty,
		MinerWalletID:    minerWallet.WalletID,
		MiningReward:     miningReward,
		Size:             int64(len(pendingTxs) * 500), // Approximate size
	}

	// Start a session for atomic operations
	session, err := database.GetClient().StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database session"})
		return
	}
	defer session.EndSession(ctx)

	// Execute atomically
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Insert the block
		_, err := getBlockCollection().InsertOne(sessCtx, newBlock)
		if err != nil {
			return nil, err
		}

		// Update transactions to confirmed
		now := time.Now()
		for _, tx := range pendingTxs {
			_, err := getTransactionCollection().UpdateOne(sessCtx,
				bson.M{"transactionId": tx.TransactionID},
				bson.M{"$set": bson.M{
					"status":      models.TxStatusConfirmed,
					"blockHash":   hash,
					"blockHeight": newIndex,
					"confirmedAt": now,
				}})
			if err != nil {
				return nil, err
			}

			// If this is a zakat transaction, update the zakat_payments status
			if tx.Type == models.TxTypeZakat {
				_, err = database.GetCollection("zakat_payments").UpdateOne(sessCtx,
					bson.M{"transactionId": tx.TransactionID},
					bson.M{"$set": bson.M{
						"status":      "confirmed",
						"confirmedAt": now,
					}})
				if err != nil {
					return nil, err
				}
			}

			// Mark UTXOs as confirmed
			for _, output := range tx.Outputs {
				_, err := getUTXOCollection().UpdateMany(sessCtx,
					bson.M{"transactionId": tx.TransactionID, "walletId": output.WalletID},
					bson.M{"$set": bson.M{
						"isConfirmed": true,
						"blockHash":   hash,
					}})
				if err != nil {
					return nil, err
				}
			}
		}

		// Create coinbase UTXO for mining reward
		coinbaseTxID := crypto.GenerateTransactionID(minerWallet.WalletID, nil, []crypto.OutputData{{WalletID: minerWallet.WalletID, Amount: miningReward}}, timestamp.Unix())
		coinbaseUTXO := models.UTXO{
			TransactionID: coinbaseTxID,
			OutputIndex:   0,
			WalletID:      minerWallet.WalletID,
			Amount:        miningReward,
			PublicKey:     minerWallet.PublicKey,
			IsSpent:       false,
			IsConfirmed:   true,
			BlockHash:     hash,
			CreatedAt:     now,
		}
		_, err = getUTXOCollection().InsertOne(sessCtx, coinbaseUTXO)
		if err != nil {
			return nil, err
		}

		// Create coinbase Transaction record for mining reward (so it shows in transaction history)
		coinbaseTx := models.Transaction{
			TransactionID: coinbaseTxID,
			Type:          models.TxTypeCoinbase,
			SenderWallet:  "",
			Outputs: []models.TransactionOutput{
				{
					WalletID:  minerWallet.WalletID,
					Amount:    miningReward,
					PublicKey: minerWallet.PublicKey,
				},
			},
			TotalInput:  0,
			TotalOutput: miningReward,
			Fee:         0,
			Status:      models.TxStatusConfirmed,
			BlockHash:   hash,
			Message:     "Mining Reward",
			Timestamp:   now,
			ConfirmedAt: &now,
		}
		_, err = getTransactionCollection().InsertOne(sessCtx, coinbaseTx)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save block", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Block mined successfully! 🎉",
		"block":        newBlock,
		"miningReward": miningReward,
		"nonce":        nonce,
		"hash":         hash,
	})
}

// calculateNewDifficulty adjusts difficulty based on recent block times
func calculateNewDifficulty(ctx context.Context, currentDifficulty int) int {
	// Get last N blocks
	cursor, err := getBlockCollection().Find(ctx, bson.M{},
		options.Find().SetSort(bson.M{"index": -1}).SetLimit(int64(models.DefaultBlockchainConfig.DifficultyAdjustment)))
	if err != nil {
		return currentDifficulty
	}
	defer cursor.Close(ctx)

	var blocks []models.Block
	if err := cursor.All(ctx, &blocks); err != nil || len(blocks) < 2 {
		return currentDifficulty
	}

	// Calculate average time between blocks
	var totalTime time.Duration
	for i := 0; i < len(blocks)-1; i++ {
		diff := blocks[i].Timestamp.Sub(blocks[i+1].Timestamp)
		totalTime += diff
	}
	avgTime := totalTime / time.Duration(len(blocks)-1)
	targetTime := time.Duration(models.DefaultBlockchainConfig.TargetBlockTime) * time.Second

	// Adjust difficulty
	if avgTime < targetTime/2 {
		return currentDifficulty + 1
	} else if avgTime > targetTime*2 && currentDifficulty > 1 {
		return currentDifficulty - 1
	}

	return currentDifficulty
}

// ValidateBlockchain validates the entire blockchain
func ValidateBlockchain(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := getBlockCollection().Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"index": 1}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blocks"})
		return
	}
	defer cursor.Close(ctx)

	var blocks []models.Block
	if err := cursor.All(ctx, &blocks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse blocks"})
		return
	}

	var errors []string
	for i := 1; i < len(blocks); i++ {
		// Check hash link
		if blocks[i].PreviousHash != blocks[i-1].Hash {
			errors = append(errors, "Block "+strconv.FormatInt(blocks[i].Index, 10)+" has invalid previous hash link")
		}

		// Check index
		if blocks[i].Index != blocks[i-1].Index+1 {
			errors = append(errors, "Block "+strconv.FormatInt(blocks[i].Index, 10)+" has invalid index")
		}

		// Verify hash meets difficulty
		if !crypto.ValidateBlockHash(blocks[i].Hash, blocks[i].Difficulty) {
			errors = append(errors, "Block "+strconv.FormatInt(blocks[i].Index, 10)+" hash doesn't meet difficulty")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"validation": models.ChainValidation{
			IsValid:       len(errors) == 0,
			BlocksChecked: int64(len(blocks)),
			Errors:        errors,
		},
	})
}

// GetMiningStatus returns current mining info
func GetMiningStatus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get pending transaction count
	pendingCount, _ := getTransactionCollection().CountDocuments(ctx, bson.M{"status": models.TxStatusPending})

	// Get latest block
	var lastBlock models.Block
	err := getBlockCollection().FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.M{"index": -1})).Decode(&lastBlock)

	difficulty := models.DefaultBlockchainConfig.InitialDifficulty
	lastBlockHash := ""
	lastBlockIndex := int64(-1)

	if err == nil {
		difficulty = lastBlock.Difficulty
		lastBlockHash = lastBlock.Hash
		lastBlockIndex = lastBlock.Index
	}

	c.JSON(http.StatusOK, gin.H{
		"pendingTransactions": pendingCount,
		"difficulty":          difficulty,
		"miningReward":        models.DefaultBlockchainConfig.BlockReward,
		"lastBlockHash":       lastBlockHash,
		"lastBlockIndex":      lastBlockIndex,
		"targetPrefix":        getTargetPrefix(difficulty),
	})
}

func getTargetPrefix(difficulty int) string {
	prefix := ""
	for i := 0; i < difficulty; i++ {
		prefix += "0"
	}
	return prefix
}

// GetMyMinedBlocks returns blocks mined by the authenticated user
func GetMyMinedBlocks(c *gin.Context) {
	userID := c.GetString("userId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get user's wallet
	var wallet models.Wallet
	err = getWalletCollection().FindOne(ctx, bson.M{"userId": objID}).Decode(&wallet)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	// Get blocks mined by this wallet
	cursor, err := getBlockCollection().Find(ctx, bson.M{"minerWalletId": wallet.WalletID},
		options.Find().SetSort(bson.M{"index": -1}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blocks"})
		return
	}
	defer cursor.Close(ctx)

	var blocks []models.Block
	if err := cursor.All(ctx, &blocks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse blocks"})
		return
	}

	// Calculate total rewards
	var totalRewards float64
	for _, block := range blocks {
		totalRewards += block.MiningReward
	}

	c.JSON(http.StatusOK, gin.H{
		"blocks":       blocks,
		"count":        len(blocks),
		"totalRewards": totalRewards,
	})
}
