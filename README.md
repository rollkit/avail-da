
# Avail DA

## Abstract

This package implements the generic DA interface defined in [go-da](https://github.com/rollkit/go-da)

## Details

The generic DA interface defines how DA implementations can submit, retrieve and validate blobs.

The Avail DA implementation connects to a local [Avail-light-node](https://github.com/availproject/avail-light) instance using the given config  and allows using Avail as the DA layer.

## Implementation

The implementation calls the corresponding Avail [node api docs](https://github.com/availproject/avail-light/blob/main/src/api/v2/README.md) methods.

### Get

Get retrieves blobs referred to by their ids.

The implementation calls [Get](https://github.com/availproject/avail-light/blob/main/src/api/v2/README.md#get-v2blocksblock_numberdatafieldsdataextrinsic) endpoint on the Avail-light-node API.

### Submit

Submit submits blobs and returns their ids and proofs.

The implementation calls [Submit](https://github.com/availproject/avail-light/blob/main/src/api/v2/README.md#post-v2submit) endpoint on the Avail-light-node API.

## Installation & Setup

### Dependencies

* Operating systems: GNU/Linux or macOS
* [Golang 1.21+](https://go.dev/)
* [Ignite CLI v28.1.0](https://github.com/ignite/cli)
* [Homebrew](https://brew.sh/)
* [wget](https://www.gnu.org/software/wget/)
* [Avail Light Node](https://github.com/availproject/avail-light)
* [Rust](https://www.rust-lang.org/tools/install)

## Avail-Light installation for local development Environment

### 1. Data availability node

* clone the repo

    ``` https://github.com/availproject/avail.git ```

* go to root folder

    ``` cd avail ```

* checkout to the following branch

    ``` git checkout v1.9.0.3 ```

* run node

    ``` cargo run --locked --release -- --dev ```
  
#### 2. Avail light  bootstrap node

* clone the repo

    ``` https://github.com/availproject/avail-light-bootstrap.git ```

* go to root folder

    ``` cd avail-light-bootstrap ```

* run node

    ``` cargo run --release ```

#### 3. Avail light node

* clone the repo

    ``` https://github.com/availproject/avail-light.git ```

* go to root folder

    ``` cd avail-light ```

* create a configuration file ```touch config.yaml``` in the root directory & put following content.

    ```yaml
    http_server_host = '127.0.0.1'
    http_server_port = 8000
    port = 38000
    tcp_port_reuse = false
    autonat_only_global_ips = false
    autonat_throttle = 1
    autonat_retry_interval = 10
    autonat_refresh_interval = 30
    autonat_boot_delay = 5
    identify_protocol = '/avail_kad/id/1.0.0'
    identify_agent = 'avail-light-client/rust-client'
    bootstraps = [["12D3KooWStAKPADXqJ7cngPYXd2mSANpdgh1xQ34aouufHA2xShz", "/ip4/127.0.0.1/udp/39000"]]
    bootstrap_period = 300
    relays = []
    full_node_ws = ['ws://127.0.0.1:9944']
    confidence = 92.0
    avail_path = 'avail_path'
    log_level = 'INFO'
    log_format_json = false
    ot_collector_endpoint = 'http://otelcollector.avail.tools:4317'
    disable_rpc = false
    disable_proof_verification = false
    dht_parallelization_limit = 20
    put_batch_size = 1000
    query_proof_rpc_parallel_tasks = 8
    max_cells_per_rpc = 30
    threshold = 5000
    kad_record_ttl = 86400
    publication_interval = 43200
    replication_interval = 10800
    replication_factor = 20
    connection_idle_timeout = 30
    query_timeout = 60
    query_parallelism = 3
    caching_max_peers = 1
    disjoint_query_paths = false
    max_kad_record_number = 2400000
    max_kad_record_size = 8192
    max_kad_provided_keys = 1024
    app_id=1

    ```

* run node

    ``` cargo run --release -- --network local -c config.yaml --clean ```

* generate avail key pair if not configured by adding the following seed phrase in ```identity.toml``` in root directory

    ``` avail_secret_seed_phrase = 'bottom drive obey lake curtain smoke basket hold race lonely fit walk//Alice' ```

## Building your soverign rollup

Now that you have a da node and light client running, we are ready to build and run our Cosmos-SDK blockchain (here we have taken [gm application](https://rollkit.dev/tutorials/gm-world))

* go to the root directory and install rollkit by adding the following lines to go.mod

    ```text

    replace github.com/cosmos/cosmos-sdk => github.com/rollkit/cosmos-sdk v0.50.1-rollkit-v0.11.9-no-fraud-proofs

    ```

  and run
  
    ```text

    go mod tidy

    ```

* start your rollup

  create one script file (init-local.sh) in root folder

  ```text
  touch init-local.sh

  ```

  add the following script to the script file (init-local.sh) or you can get the script from [here](https://gist.github.com/chandiniv1/27397b93e08e2c40e7e1b746f13e5d7b)

* to make use of avail as a da layer,

  * clone the repo

    ``` git clone https://github.com/rollkit/avail-da.git ```
  
  * Go to the root dir
  
    ``` cd avail-da ```

  * run the server

    ``` go run ./cmd/avail-da/main.go ```

  note: make sure that port address(serving avail-da) must be the ```rollkit.da_address``` in the [script file](https://gist.github.com/chandiniv1/27397b93e08e2c40e7e1b746f13e5d7b)

* run the rollup chain

    go to root of the gm repo and run

    ```text
    bash init-local.sh

    ```

* list your keys
  
  ``` gmd keys list --keyring-backend test ```
  
  You should see an output like the following

  ```text
  - address: gm1ffdft5ku0qw67eypavgyltjqj54yraaa4uj8pl
    name: gm-key-2
    pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AiRf1Vxwcuunu1yjPpSUPxW85Q9NhM1y0Dg3ozOEGUto"}'
    type: local
  - address: gm1r4g9lleykkw6mjdpmp6e0tgusymh4sa2swcw69
    name: gm-key
    pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"A9hnpPl0G6VOBijJNheFXvJDg1O7NszzA8innoXF81N5"}'
    type: local
  ```

* Now we can test by sending the transaction by sending amount from one account to another
  
  ```text

  gmd tx bank send [from_key_or_address] [to_address] [amount] [flags]

  ```

  ex:

  ```text

   gmd tx bank send gm1ffdft5ku0qw67eypavgyltjqj54yraaa4uj8pl gm1r4g9lleykkw6mjdpmp6e0tgusymh4sa2swcw69 11stake --chain-id gm --keyring-backend test

  ```

  You'll be prompted to accept the transaction:

  ```text
  auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
  body:
    extension_options: []
    memo: ""
    messages:
    - '@type': /cosmos.bank.v1beta1.MsgSend
      amount:
      - amount: "11"
        denom: stake
      from_address: gm1ffdft5ku0qw67eypavgyltjqj54yraaa4uj8pl
      to_address: gm1r4g9lleykkw6mjdpmp6e0tgusymh4sa2swcw69
    non_critical_extension_options: []
    timeout_height: "0"
  signatures: []
  confirm transaction before signing and broadcasting [y/N]: 

  ```

  Type y if you'd like to confirm and sign the transaction. Then, you'll see the confirmation:

  ```text
  confirm transaction before signing and broadcasting [y/N]: y
  code: 0
  codespace: ""
  data: ""
  events: []
  gas_used: "0"
  gas_wanted: "0"
  height: "0"
  info: ""
  logs: []
  raw_log: ""
  timestamp: ""
  tx: null
  txhash: 130EA420F2373C88F6191E1D203CEF272B666BE283316A17BC8B02FBABCBA1C9

  ```

  you can query the tx using using

  ```text

  gmd q tx <hash>

  ```

  ex:

  ```text

  gmd q tx 130EA420F2373C88F6191E1D203CEF272B666BE283316A17BC8B02FBABCBA1C9

  ```

  then you'll see

  ```text

  code: 0
  codespace: ""
  data:   12260A242F636F736D6F732E62616E6B2E763162657461312E4D736753656E64526573706F6E736  5
  events:
  - attributes:
    - index: true
      key: fee
      value: ""
    - index: true
      key: fee_payer
      value: gm1ffdft5ku0qw67eypavgyltjqj54yraaa4uj8pl
    type: tx
  - attributes:
    - index: true
      key: acc_seq
      value: gm1ffdft5ku0qw67eypavgyltjqj54yraaa4uj8pl/0
    type: tx
  - attributes:
    - index: true
      key: signature
      value: y2ZU5c/CLcJpKbT8e2uZZz3T5buO/  efa9yNuODgLBd90xt4c6ErNf2H1OQnZerdpdez1H3kWdLyk769Pb2Sisg==
    type: tx
  - attributes:
    - index: true
      key: action
      value: /cosmos.bank.v1beta1.MsgSend
    - index: true
      key: sender
      value: gm1ffdft5ku0qw67eypavgyltjqj54yraaa4uj8pl
    - index: true
      key: module
      value: bank
    - index: true
      key: msg_index
      value: "0"
    type: message
  - attributes:
    - index: true
      key: spender
      value: gm1ffdft5ku0qw67eypavgyltjqj54yraaa4uj8pl
    - index: true
      key: amount
      value: 11stake
    - index: true
      key: msg_index
      value: "0"
    type: coin_spent
  - attributes:
    - index: true
      key: receiver
      value: gm1r4g9lleykkw6mjdpmp6e0tgusymh4sa2swcw69
    - index: true
      key: amount
      value: 11stake
    - index: true
      key: msg_index
      value: "0"
    type: coin_received
  - attributes:
    - index: true
      key: recipient
      value: gm1r4g9lleykkw6mjdpmp6e0tgusymh4sa2swcw69
    - index: true
      key: sender
      value: gm1ffdft5ku0qw67eypavgyltjqj54yraaa4uj8pl
    - index: true
      key: amount
      value: 11stake
    - index: true
      key: msg_index
      value: "0"
    type: transfer
  - attributes:
    - index: true
      key: sender
      value: gm1ffdft5ku0qw67eypavgyltjqj54yraaa4uj8pl
    - index: true
      key: msg_index
      value: "0"
    type: message
  gas_used: "61687"
  gas_wanted: "200000"
  height: "103"
  info: ""
  logs: []
  raw_log: ""
  timestamp: "2024-01-08T07:14:14Z"
  tx:
    '@type': /cosmos.tx.v1beta1.Tx
    auth_info:
      fee:
        amount: []
        gas_limit: "200000"
        granter: ""
        payer: ""
      signer_infos:
      - mode_info:
          single:
            mode: SIGN_MODE_DIRECT
        public_key:
          '@type': /cosmos.crypto.secp256k1.PubKey
          key: AiRf1Vxwcuunu1yjPpSUPxW85Q9NhM1y0Dg3ozOEGUto
        sequence: "0"
      tip: null
    body:
      extension_options: []
      memo: ""
      messages:
      - '@type': /cosmos.bank.v1beta1.MsgSend
        amount:
        - amount: "11"
          denom: stake
        from_address: gm1ffdft5ku0qw67eypavgyltjqj54yraaa4uj8pl
        to_address: gm1r4g9lleykkw6mjdpmp6e0tgusymh4sa2swcw69
      non_critical_extension_options: []
      timeout_height: "0"
    signatures:
    - y2ZU5c/CLcJpKbT8e2uZZz3T5buO/  efa9yNuODgLBd90xt4c6ErNf2H1OQnZerdpdez1H3kWdLyk769Pb2Sisg==
  txhash: 130EA420F2373C88F6191E1D203CEF272B666BE283316A17BC8B02FBABCBA1C9


  ```

  then query the bank balances

  ```text

  gmd query bank balances gm1ffdft5ku0qw67eypavgyltjqj54yraaa4uj8pl

  ```

  you can see

  ```text

  balances:
  - amount: "9999999999999999999999989"
    denom: stake
  pagination:
    total: "1"

  ```

  query the balance of other key

  ```text

  gmd query bank balances gm1r4g9lleykkw6mjdpmp6e0tgusymh4sa2swcw69

  ```

  you can see

  ```text
  balances:
  - amount: "10000000000000000000000011"
    denom: stake
  pagination:
    total: "1"
  
  ```

With this You've built a local rollup that posts to a local avail light node
