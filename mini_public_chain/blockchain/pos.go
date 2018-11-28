package blockchain

import (
	"time"
	"math/rand"
)

// Blockchain is a series of validated Blocks
var tempBlocks []Block

// candidateBlocks handles incoming blocks for validation
var candidateBlocks = make(chan Block)

var hasbeenValid = make(chan int, 1)

// validators keeps track of open validators and balances
var validators = make(map[string]int)

func pos()  {
	go func() {
		for candidate := range candidateBlocks {
			mutex.Lock()
			tempBlocks = append(tempBlocks, candidate)
			mutex.Unlock()
		}
	}()

	go func() {
		for {
			pickWinner()
		}
	}()
}

// Increment numPowWinners every this number of miner peers.
const powWinnerIncrementInterval = 50

// Check number of miner peers, and adjust pow winners accordingly.
// Use number of accounts to approximate miner peers.
func maybeAdjustPowWinners() {
  numAccounts := len(BlockchainInstance.Accounts)
  if numAccounts <= 2 {
	BlockchainInstance.NumPowWinners = 1
  } else {
  	BlockchainInstance.NumPowWinners = 2 + numAccounts / powWinnerIncrementInterval
  }
}

// pickWinner creates a lottery pool of validators and chooses the validator who gets to forge a block to the blockchain
// by random selecting from the pool, weighted by amount of tokens staked
func pickWinner() {
	// This sleep number if not set properly, may even cause deadlock.
	time.Sleep(1 * time.Second)
	mutex.Lock()
	temp := tempBlocks
	mutex.Unlock()

	lotteryPool := []string{}
	if len(temp) > 0 {

		// slightly modified traditional proof of stake algorithm
		// from all validators who submitted a block, weight them by the number of staked tokens
		// in traditional proof of stake, validators can participate without submitting a block to be forged
	OUTER:
		for _, block := range temp {
			// if already in lottery pool, skip
			for _, node := range lotteryPool {
				if block.Validator == node {
					continue OUTER
				}
			}

			// lock list of validators to prevent data race
			mutex.Lock()
			setValidators := validators
			mutex.Unlock()

			k, ok := setValidators[block.Validator]
			if ok {
				for i := 0; i < k; i++ {
					lotteryPool = append(lotteryPool, block.Validator)
				}
			}
		}

		// randomly pick winner from lottery pool
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
                
		// default pos only has 1 winner.
                winners := 1
		if BlockchainInstance.Proof == "pox" {
			maybeAdjustPowWinners()
			winners = BlockchainInstance.NumPowWinners
		}
		for i := 0; i < winners; i++ {
			lotteryWinner := lotteryPool[r.Intn(len(lotteryPool))]
			// add block of winner to blockchain and let all the other nodes know
			for _, block := range temp {
				if block.Validator == lotteryWinner {
					//mutex.Lock()
					//BlockchainInstance.Blocks = append(BlockchainInstance.Blocks, block)
					//mutex.Unlock()
					hasbeenValid<-1
					break
				}
			}
		}
	}

	mutex.Lock()
	tempBlocks = []Block{}
	mutex.Unlock()
}
