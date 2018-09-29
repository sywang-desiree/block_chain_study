Modification from github.com/t-bast/ring-signatures/main.go for testing the performance of ring-signatures with different number of decoys.

I basically only modified the generate command.

Example command:
$ go run block_chain_study/ring-signatures-test/main.go generate --decoy 5 --message hello

