# Mini public chain demo 

The code here is modification based on https://github.com/GenaroSanada/Course. I wrote and updated it occasionally for learning and playing with block chain infrastructure, with no guarantee on correctness and completeness.

*   It has very simple blocks, transactions, accounts structures. It can mine new blocks, do simple transactions and special transactions.

*   The concensus supported includes PoS, PoW, and PoX (PoS + PoW). The PoX is using PoS to choose small number of winners to do PoW, trying to get benefits of both PoS and PoW.

*   The governance is typical governance methods using by PoS, PoW system, with reward for mining new blocks and packaging transactions, slash for bad behaviors, and dynamically increase the number of PoW winners as number of miner peers increase significantly.

*   I added -i flag to specify a non 127.0.0.1 ip for the peer. If empty, the code will try to find the outbound ip for the peer. Currently it seems to find 192.168... ip instead of actual outbound ip.

*   Use --demo flag for manual input result in each round of mining and watching block mining process more slowly. Default is false, and does not require manul input result.

Test commands that I used on different terminals (assume in cmd directory) are

```
// (Optional) create local wallet and accounts. Each peer should have its own wallet (replace "desiree" in createwallet command) and accounts.
// Run this command multiple times to create multiple accounts for one wallet.
// Later, -s and -a flags in the chain mode commands should use those wallet and accounts addresses.
go run main.go -c account desiree createwallet

// Bootstrap peer with a wallet (-s flag) and a miner account (-a flag).
go run main.go -c chain -s desiree -l 8080 -a 1FghRtifoTLuMsFRacRBpBYD2VwLmGoAhW -p pox -i ""

// Another peer directly dialing to the bootstrap peer (-d flag) with a different wallet (-s flag) and another miner account (-a flag).
// Each peer can have its own wallet and accounts. -s and -a flags should be stable after having wallet and accounts created.
// The address in -d flag could be a bit dynamic, but depends on the address print out by the bootstrap peer.
//
go run main.go  -c chain -s lzx -l 8082 -a 1Hn94smEVwEd3kPfvF39ozhqCQKGqce5qc -d /ip4/192.168.1.6/tcp/8080/ipfs/QmaHAkUhArtD2UwW4PMRGzZxCSVV6SAssy2C6XJRjonUWR -p pox

// Generate transactions among the running peers' accounts. This outputs a test transaction file.
// --wallets is for comma separated wallet suffix.
// -n is for number of transactions to generate.
// Optionally, --address is to specify http post address, default is 127.0.0.1:8081.
//
go run gen_test_txns.go --wallets desiree,lzx -n 8

// Now, run the generated transactions from the test transaction file.
. test_http_post_txns

// Of course, you can manually send transactions.
curl -i --request POST --header 'Content-Type: application/json' --data '{"From":"1FghRtifoTLuMsFRacRBpBYD2VwLmGoAhW","To":"1BvT54va6zRhos2rVkT4DDSMexTCtT4q6J","Value":100,"Data":"message2"}' http://127.0.0.1:8081/txpool
```

