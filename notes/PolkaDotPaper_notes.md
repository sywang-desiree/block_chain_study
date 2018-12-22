## Overview

Relay chain to connect heterogeneous chains (para chains).

No one size fits all.

Decouple consensus from state transition.

Two points: pooled security, trust free interchain transactions

On related work, side chains address extensibility rather than scalability.

For cosmos, validator sets for the zoned chains, especially incentivizing then, like side chains, remain unsolved.

## Tech Details

4 **roles** in Polkadot: nominator (nominates) → validator <-- collator (full node, helps produce blocks), fisherman (timely proof for spotting misbehaving validator)

These roles and related governance integrate ideas from Ethereum, Tendermint, and TrueBit.

Low level **consensus**: BFT

Consensus for determining roles: PoS

Checkpoint latch to prevent long range attack

The parachains’ headers are sealed within the relay-chain block.

Validators are randomly segmented into subsets; one subset per parachain, the subsets potentially differing per block.

Polkadot’s relay-chain itself will probably exist as an Ethereum-like accounts and state chain.

**Interchain transactions**
*   payment managed through negotiation logic on the source and destination parachains. 
*   queuing based around a Merkle tree
*   move transactions on the output queue of one parachain into the input queue of the destination parachain.
*   Polkadot → Ethereum
    *   Use majority of validators signatures, or threshold signatures.
*   Ethereum → Polkadot
    *   Break out contract, emit logs, verify block header’s validity and canonicality.
*   Bitcoin <-> Polkadot
    *   More difficult. How the deposits can be se-curely controlled from a rotating validator set?

**Relay chain ops**
*   If implemented with EVM, need to be implemented as built in contracts.
*   If implemented with WASM, no need to do that.

**Staking contract**
*   Used for managing validators.
*   To remain stake token liquidity, not all tokens are staked.
*   Nominating validators via approval voting. 
*   The reward needs to large enough to make validation process worthwhile, but not so large to be subject to attacks that force validators misbehaving.
*   In cases that validators cannot validate each other, e.g. multi forks, fishermen come to support.

**Parachain registry**
*   New parachain addition needs a hard fork now for validation, and full referendum (e.g. 2/3) voting
*   Suspension of parachain via dynamic validator voting.
*   Removal of parachain via full referendum voting.

**Sealing relay chain blocks**
*   Each participant has signed votes info of other participants re. availability and validity of blocks.
*   ⅔ validators vote for validity. ⅓ validators vote for availability.
*   For performance, chains can grow (along with the original chain validator set), whereas the participants can remain at the least sub-linear.
*   Weight parachain blocks. Avoid over-weight blocks via 1) validators publishing performance data on blocks, and 2) collator insurance with funds or block proposing success history.

