# Mini public chain demo 

The code here is modification based on https://github.com/GenaroSanada/Course. I wrote and updated it occasionally for learning and playing with block chain infrastructure, with no guarantee on correctness and completeness.

*   It has very simple blocks, transactions, accounts structures. It can mine new blocks, do simple transactions and special transactions.

*   The concensus supported includes PoS, PoW, and PoX (PoS + PoW). The PoX is using PoS to choose small number of winners to do PoW, trying to get benefits of both PoS and PoW.

*   The governance is typical governance methods using by PoS, PoW system, with reward for mining new blocks and packaging transactions, slash for bad behaviors, and dynamically increase the number of PoW winners as number of miner peers increase significantly.

Test commands that I used on different terminals are

```
go run main.go -c chain -s lzhx_ -l 8080 -a 15QuLUSY8m4B1GyMGBr82nEjwzGwV82Wvi -p pox
```

```
go run main.go  -c chain -s lzx -l 8082 -a 1Hn94smEVwEd3kPfvF39ozhqCQKGqce5qc -p pox
```

// Run this only after the peers running on local machine are connected and blockchain has >= 1-2 blocks.
```
curl -i --request POST --header 'Content-Type: application/json' --data '{"From":"1Hn94smEVwEd3kPvF39ozhqCQKGqce5qc","To":"15QuLUSY8m4B1GyMGBr82nEjwzGwV82Wvi","Value":100,"Data":"message2"}' http://127.0.0.1:8083/txpool 
```
