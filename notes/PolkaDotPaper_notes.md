## Overview

Relay chain to connect heterogeneous chains (para chains).
No one size fits all.
Decouple consensus from state transition.
Two points: pooled security, trust free interchain transactions

On related work, side chains address extensibility rather than scalability.
For cosmos, validator sets for the zoned chains, especially incentivizing then, like side chains, remain unsolved.

## Tech Details

4 roles in Polkadot: nominator (nominates) → validator <-- collator (full node, helps produce blocks), fisherman (timely proof for spotting misbehaving validator)
These roles and related governance integrate ideas from Ethereum, Tendermint, and TrueBit.

Low level consensus: BFT
Consensus for determining roles: PoS
Checkpoint latch to prevent long range attack

The parachains’ headers are sealed within the relay-chain block.
Validators are randomly segmented into subsets; one subset per parachain, the subsets potentially differing per block.
Polkadot’s relay-chain itself will probably exist as an Ethereum-like accounts and state chain.

Interchain transactions

*   payment managed through negotiation logic on the source and destination parachains. 
*   queuing based around a Merkle tree
*   move transactions on the output queue of one parachain into the input queue of the destination parachain.
*   Polkadot → Ethereum
    *   Use majority of validators signatures, or threshold signatures.
*   Ethereum → Polkadot
    *   Break out contract, emit logs, verify block header’s validity and canonicality.
*   Bitcoin <-> Polkadot
    *   More difficult. How the deposits can be se-curely controlled from a rotating validator set?
