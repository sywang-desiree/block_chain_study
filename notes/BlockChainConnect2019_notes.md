**Demos**
*   Hardware based wallet to store keys.
*   B2B transaction chain, private chain, transactions stored only on transaction parties, permission only system.
*   Comic rights and PR blockchain 
*   Smart contract code audit

**Tencent blockchain use cases**
*   e-invoice
*   Un-temperable medical prescriptions

**Synchronized with a chance of partition tolerance**
*   State machine replication
*   Large scale consensus Thunderella
    *   Assume Committee, leader, ¾ votes —> notorized.
    *   Honest majority for consistency
    *   ¾ votes, honest leader for liveness
    *   Fall back to slow chain to fix voting problem in fast path
    *   Fast path + slow chain
        *   Fast path log, heartbeat to slow chain. Heartbeat served as checkpoint.
*   Math model for consensus
    *   Flaw in synchronous model
    *   Synchronous network: long delay-> fault
    *   Partially synchronous: long delay is not considered fault.

**Taxa network: layer 2 infra for dapps** 
*   Privacy preserving: tee (trusted execution environment)
*   Python smart contract. Stateless model.
*   IO: ethereum, IPFS
*   Applications: big data analytics, shared economy, game.
*   Case study: private set intersection, multiplayer dealer game.

**Token economy panel**
*   Data, content, advertisers incentives problem
*   Socially satisfactory allocation. Blockchain smart contracts regulate parties to cooperate.
*   Micro transaction could be incentive for content right monetization.
*   Exchange mobile network signals and bandwidth
*   Application scenarios with transactional values that can be transferred: e.g. tokenize resources and make them shareable.

**Gaming, Finance and blockchain integrations panel**
*   “ICO is dead” impression due to bear market.
*   Game assets and rewards. Digital identities.
*   Characteristics of potential integrations: multiple parties rather than single authority, fault tolerance.


**Permissioned and Permissionless by VMWare**
*   Enterprise: all about decentralized trust
*   Crypto currencies: public, anonymous, permission less.
*   Enterprise: parties and validators, permissioned. Need to be energy efficient.
*   BFT consensus: concord scalable BFT open source project
*   Combine votes with threshold crypto to reduce communication cost
*   Interoperability of public and private chains
    *   Path to standardization
    *   Token exchange
    *   Smart contracts 
    *   State exchange
    *   All roads lead to BFT. Quorum certificate. Generalize vote rules.

**Enterprise blockchain research by Accenture**
*   AI and blockchain: data as service. AI is mostly as good as amount of data.

**Authchain by visa research**
*   Online shopping problem: cart abandonment. Large part of reason is not trusting merchant.
*   Digital wallet problem: another entity to trust. 
*   Merchants need to deal with unknown payments.
*   Their solution: proxy re-encryption
    *   Non interactive. Non transitive
    *   Proxy after payment gateway to authenticate & re-encrypt. Use blockchain as proxy to reduce need of trust.

