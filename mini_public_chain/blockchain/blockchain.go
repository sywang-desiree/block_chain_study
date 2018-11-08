package blockchain

import (
	"bufio"
	"crypto/sha256"
        "encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
        "path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/GenaroSanada/Course/wallet"
)

// Per node wallets.
// Benefiary if this node mines a new block.
var Wallets wallet.Wallets

var DataFileName string = "chainData.txt"

// Governance
const Difficulty = 1

// Reward for mining a new block.
const NewBlockReward = 1000

// Reward for packaging a transaction in a new block.
// Assume penalty is deducting the same amount as transaction reward.
const TransactionReward = 1

// Increment numPowWinners every this number of miner peers.
const PowWinnerIncrementInterval = 50

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int `json:"index"`
	Timestamp string `json:"timestamp"`
	Result       int `json:"result"`
	Hash      string `json:"hash"`
	PrevHash  string `json:"prevhash"`
	Proof        uint64           `json:"proof"`

        // Used for PoW
        Difficulty   int `json: "difficulty"`
        Nonce	  string `json:"nonce"`

        // Used for PoS
        Validator string `json:"validator"`  

	Transactions map[string]Transaction `json:"transactions"`
	Accounts   map[string]Account  `json:"accounts"`
	MinerReward	uint64 `json:"minerreward"`
}

type Account struct {
	Balance uint64 `json:"balance"`
	State   uint64 `json:"state"`
}

type SpecialPayload struct {
	ResetAllAccountsState	bool `json:"reset_all_accounts_state"`
}

type Transaction struct {
	Amount    uint64    `json:"amount"`
	Recipient string `json:"recipient"`
	// If no sender, assume it is reward or penalty (negative Amount).
	Sender    string `json:"sender"`
	Timestamp string `json:"timestamp"`
	Processed bool `json:"processed"`
	Data      []byte `json:"data"`
	// Special transaction beyond transfer amount
	SpecialPayload	SpecialPayload	`json:"special_payload"`
}

// Changed to store all transactions instead of just pending transaction for the sake of correct synchronization of
// whether a transaction is processed.
type TxPool struct {
	AllTx	map[string]Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		AllTx:   make(map[string]Transaction),
	}
}


func (p *TxPool)Clear() bool {
	if len(p.AllTx) == 0 {
		return true
	}
	p.AllTx = make(map[string]Transaction)
	return true
}

// Blockchain is a series of validated Blocks
type Blockchain struct {
	Blocks []Block
	TxPool *TxPool
	DataDir string
        // Proof of XXX: pow, pos, or pox
        Proof 	string
	NumPowWinners int
}

// For communication among peers
type BlockchainMessage struct {
	Blocks []Block	`json:"blocks"`
	Transactions map[string]Transaction	`json:"transactions"`
}

func (t *Blockchain) NewTransaction(sender string, recipient string, amount uint64, data []byte) *Transaction {
	transaction := new(Transaction)
	transaction.Sender = sender
	transaction.Recipient = recipient
	transaction.Amount = amount
	transaction.Timestamp = time.Now().String()
	transaction.Data = data
	transaction.Processed = false

	return transaction
}

func (t *Blockchain) NewSpecialTransaction(data []byte, specialPayload SpecialPayload) *Transaction {
	transaction := new(Transaction)
	transaction.Timestamp = time.Now().String()
	transaction.Data = data
	transaction.SpecialPayload = specialPayload
	transaction.Processed = false
	return transaction
}

var mutex = &sync.Mutex{}

func (t *Blockchain)AddTxPool(tx *Transaction) int {
	mutex.Lock()
	t.TxPool.AllTx[CalculateTransactionHash(*tx)] = *tx
	mutex.Unlock()
	return len(t.TxPool.AllTx)
}

func (t *Blockchain) LastBlock() Block {
	return t.Blocks[len(t.Blocks)-1]
}

func (t *Blockchain) GetBalance(address string) uint64 {
	accounts := t.LastBlock().Accounts
	if value, ok := accounts[address]; ok {
		return value.Balance
	}
	return 0
}


func (t *Blockchain)PackageTx(newBlock *Block) {
	(*newBlock).Transactions = t.TxPool.AllTx
	AccountsMap := t.LastBlock().Accounts
	for k1, v1 := range AccountsMap {
		fmt.Println(k1, "--", v1)
	}

	unusedTx := make([]Transaction,0)

	for k, v := range t.TxPool.AllTx{
		if v.Processed {
			continue
		}

		if v.SpecialPayload.ResetAllAccountsState {
			log.Println("Sadly we are going to reset all accounts' states")
			for ak, av := range AccountsMap {
				av.State = 0
				AccountsMap[ak] = av
			}
			v.Processed = true
			t.TxPool.AllTx[k] = v
			continue
		}		

		if value, ok := AccountsMap[v.Sender]; ok {
			if value.Balance < v.Amount{
				unusedTx = append(unusedTx, v)
				continue
			}
			value.Balance -= v.Amount
			value.State += 1
			v.Processed = true
			t.TxPool.AllTx[k] = v			
			AccountsMap[v.Sender] = value
		}

		if value, ok := AccountsMap[v.Recipient]; ok {
			value.Balance += v.Amount
			AccountsMap[v.Recipient] = value
		}else {

			newAccount := new(Account)
			newAccount.Balance = v.Amount
			newAccount.State = 0
			AccountsMap[v.Recipient] = *newAccount
		}
		(*newBlock).MinerReward += uint64(TransactionReward)
	}

    //t.TxPool.Clear()
    //余额不够的交易放回交易池
    //if len(unusedTx) > 0 {
//		for _, v := range unusedTx{
//			t.AddTxPool(&v)
//		}
//	}

	(*newBlock).Accounts = AccountsMap
}

func (t *Blockchain)WriteData2File() {
	if t.DataDir == "" {
		return
	}

	joinPath := filepath.Join(t.DataDir, DataFileName)

	file,err := os.OpenFile(joinPath,os.O_WRONLY|os.O_CREATE,0755)  //以写方式打开文件
	if err != nil {
		log.Println("can't write data to file, open file fail err:",err)
		return
	}
	defer file.Close()
	enc := gob.NewEncoder(file)

	if err := enc.Encode(t); err != nil {
		log.Fatal("encode error:", err)
	}

	fmt.Println()
	fmt.Printf("\n%sfile:%s\n>", "已配置数据存储目录，写入当前数据到存储目录中.", filepath.Join(t.DataDir,DataFileName))
}

func (t *Blockchain)ReadDataFromFile() {
	if t.DataDir == "" {
		return
	}

	joinPath := filepath.Join(t.DataDir, DataFileName)

	if !IsExist(joinPath) {
		return
	}

	file,err := os.Open(joinPath)  //以写方式打开文件
	if err != nil {
		log.Println("can't read data from file, open file fail err:",err)
		return
	}
	defer file.Close()
	dec := gob.NewDecoder(file)

	var blockchainInstanceFromFile Blockchain
	if err := dec.Decode(&blockchainInstanceFromFile); err != nil {
		log.Fatal("decode error:", err)
	}

	BlockchainInstance = blockchainInstanceFromFile
}

var BlockchainInstance Blockchain = Blockchain{
	TxPool : NewTxPool(),
}

func Lock(){
	mutex.Lock()
}

func UnLock(){
	mutex.Unlock()
}

func PrintBlockchain() {
	var tMessage BlockchainMessage
	tMessage.Blocks = BlockchainInstance.Blocks
	tMessage.Transactions = BlockchainInstance.TxPool.AllTx

	bytes, err := json.MarshalIndent(tMessage, "", "  ")
	if err != nil {
		log.Fatal(err)
        }
        // Green console color:         \x1b[32m
        // Reset console color:         \x1b[0m
        fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
}

func WriteBlockchainJson(rw *bufio.ReadWriter) {
	var tMessage BlockchainMessage
	
	mutex.Lock()
	tMessage.Blocks = BlockchainInstance.Blocks
	tMessage.Transactions = BlockchainInstance.TxPool.AllTx
	bytes, err := json.Marshal(tMessage)
	if err != nil {
		log.Println(err)
	}
	rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
	rw.Flush()
	mutex.Unlock()
}

func ProcessBlocksDiff(blocks []Block) bool {
	chain_length := len(BlockchainInstance.Blocks)
	if len(blocks) > chain_length {
		BlockchainInstance.Blocks = blocks
		return true
	}

	// Process same length chain with different last block. Here maybe using BFT is better. 
	if len(blocks) == chain_length && blocks[len(blocks) - 1].Timestamp < BlockchainInstance.Blocks[chain_length - 1].Timestamp {
		BlockchainInstance.Blocks = blocks
		return true
	}
	return false
}

func ProcessTransactionsDiff(coming_tx map[string]Transaction) bool {
	has_diff := false
	for k, v := range coming_tx {
		if value, ok := BlockchainInstance.TxPool.AllTx[k]; ok {
			if v.Processed && !value.Processed {
				has_diff = true
				log.Println("Processed transaction found")
				BlockchainInstance.TxPool.AllTx[k] = v
			}
		} else if !v.Processed {
			has_diff = true
			log.Println("New pending transaction found")
			BlockchainInstance.TxPool.AllTx[k] = v
		}
	}
	return has_diff
}

func ReadData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {

			var tMessage BlockchainMessage
			if err := json.Unmarshal([]byte(str), &tMessage); err != nil {
				log.Fatal(err)
			}

			has_diff := false
			mutex.Lock()
			has_diff = ProcessBlocksDiff(tMessage.Blocks) || ProcessTransactionsDiff(tMessage.Transactions)
			if has_diff {
				PrintBlockchain()
                                BlockchainInstance.WriteData2File()
			}
			mutex.Unlock()
		}
	}
}

func WriteData(rw *bufio.ReadWriter) {
	if BlockchainInstance.Proof == "pos" || BlockchainInstance.Proof == "pox" {
		go pos()
	}

	go func() {
		for {
			time.Sleep(10 * time.Second)
			WriteBlockchainJson(rw)

		}
	}()

        stdReader := bufio.NewReader(os.Stdin)
	validator := ""

	for {	
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		sendData = strings.Replace(sendData, "\n", "", -1)
		_result, err := strconv.Atoi(sendData)
		if err != nil {
			log.Fatal(err)
		}

		if BlockchainInstance.Proof == "pos" || BlockchainInstance.Proof == "pox" {
                	validator = GenPosValidator(stdReader)
        	}

		newBlock := GenerateBlock(BlockchainInstance.Blocks[len(BlockchainInstance.Blocks)-1], _result, BlockchainInstance.Proof, validator)
		
		if len(BlockchainInstance.TxPool.AllTx) > 0 {
			BlockchainInstance.PackageTx(&newBlock)
		}else {
			newBlock.Accounts = BlockchainInstance.LastBlock().Accounts
			newBlock.Transactions = make(map[string]Transaction)
		}

		prevBlock := BlockchainInstance.Blocks[len(BlockchainInstance.Blocks)-1]
		if IsBlockValid(newBlock, prevBlock) {
 			log.Printf("newBlock is valid. Total %d blocks\n", len(BlockchainInstance.Blocks))
			mutex.Lock()
			if BlockchainInstance.Proof == "pow" {
				// Should finalize block and reward here.
				BlockchainInstance.Blocks = append(BlockchainInstance.Blocks, newBlock)
			} else if BlockchainInstance.Proof == "pos" || BlockchainInstance.Proof == "pox" {
				candidateBlocks <- newBlock
			}
			mutex.Unlock()
		}
		if BlockchainInstance.Proof == "pos" || BlockchainInstance.Proof == "pox" {
			<-hasbeenValid
 			blockValid := true
			if BlockchainInstance.Proof == "pox" {
				newBlock = pow(newBlock)
				// prevBlock may change at this point.
				prevBlock = BlockchainInstance.Blocks[len(BlockchainInstance.Blocks)-1]
				blockValid = IsBlockValid(newBlock, prevBlock)
			}
			if blockValid {
				// Finalize and calculate reward
				mutex.Lock()
				newBlock.MinerReward += uint64(NewBlockReward)
				if Wallets.DefaultMinerAccount != "" {
					minerAccount := newBlock.Accounts[Wallets.DefaultMinerAccount]
					minerAccount.Balance += newBlock.MinerReward
					newBlock.Accounts[Wallets.DefaultMinerAccount] = minerAccount
				}
				// TODO: put this reward into transactions as well.

                                BlockchainInstance.Blocks = append(BlockchainInstance.Blocks, newBlock)
                                mutex.Unlock()
				
				// Apply MinerReward
				if Wallets.DefaultMinerAccount != "" {
                        		wall := Wallets.GetWallet(Wallets.DefaultMinerAccount)
                        		wall.Balance += newBlock.MinerReward
                        		wall.State += 1
				}
			}
		}

		BlockchainInstance.WriteData2File()

		PrintBlockchain()
		WriteBlockchainJson(rw)
	}

}

func GenPosValidator(stdReader *bufio.Reader) string {
	fmt.Print("Enter token balance:")
	scanBalance, err := stdReader.ReadString('\n')
	if err != nil {
        	log.Fatal(err)
        }
        scanBalance = strings.Replace(scanBalance, "\n", "", -1)
        balance, err := strconv.Atoi(scanBalance)
        if err != nil {
        	log.Printf("%v not a number token balance: %v\n", scanBalance, err)
                return ""
        }

        t := time.Now()
        validator := CalculateHash(t.String())
        validators[validator] = balance
        log.Printf("validator %s\n", validator)

	return validator
}

// make sure block is valid by checking index, and comparing the hash of the previous block
func IsBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if CalculateBlockHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// SHA256 hashing
func CalculateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// SHA256 hashing
func CalculateBlockHash(block Block) string {
        record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.Result) + block.PrevHash + block.Nonce
        return CalculateHash(record)
}

func CalculateTransactionHash(tx Transaction) string {
	record := tx.Sender + tx.Recipient + strconv.Itoa(int(tx.Amount)) + tx.Timestamp + string(tx.Data)
	return CalculateHash(record)
}

// create a new block using previous block's hash.
// Proof specifies whether to use pow or pos for reaching consensus.
// validator is only used in pos.
func GenerateBlock(oldBlock Block, Result int, Proof string, validator string) Block {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Result = Result
	newBlock.PrevHash = oldBlock.Hash
        newBlock.Difficulty = Difficulty
        newBlock.Validator = validator
	newBlock.MinerReward = 0

        log.Printf("Use %s to generate new block\n", Proof)
        if Proof == "pow" {
		newBlock = pow(newBlock)
        } else if Proof == "pos" || Proof == "pox" {
		newBlock.Hash = CalculateBlockHash(newBlock)
	}
	return newBlock
}

// pow method could only update Nonce and Hash of newBlock. 
func pow(newBlock Block) Block {
	for i := 0; ; i++ {
        	hex := fmt.Sprintf("%x", i)
                newBlock.Nonce = hex
                if !isHashValid(CalculateBlockHash(newBlock), newBlock.Difficulty) {
                	fmt.Println(CalculateBlockHash(newBlock), " do more work!")
                        time.Sleep(time.Second)
                        continue
                } else {
                	fmt.Println(CalculateBlockHash(newBlock), " work done!")
                        newBlock.Hash = CalculateBlockHash(newBlock)
                        break
                }

	}
	return newBlock
}

func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

func IsExist(file string) bool {
	fi, e := os.Stat(file)
	if e != nil {
		return os.IsExist(e)
	}
	return !fi.IsDir()
}
