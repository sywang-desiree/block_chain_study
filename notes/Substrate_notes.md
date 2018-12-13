Substrate is not a blockchain. It is a tech stack to build a blockchain, para chain, sharding, etc. 
Eventually it is on ramp for polkadot.

PBFT + Aurand consensus

Latency finality trade off

Two levels of finalities
*   PBFT needs to slash large number of nodes, expensive, runs every few blocks as check point.
*   Aurand intermediate fast consensus, Aura PoA.

Designed to be light client
*   Header format: Parent hash, extrinsic root, receipt hash 
*   Receipt 
    *   digest - has events of validators alterations.
    *   Storage root
    *   Change log root - Merkle tree key that has changed data. Efficient diff.

Subtrate runtime v1
*   Could be similar to EWF chain
*   If implemented using sharding + Ethereum, can only support async calls, low tps.
*   Blitz protocol distributed collation
*   Ephemeral sub chain, similar to state channel, high throughput.
*   Polkadot, still limited app utility, async calls.

