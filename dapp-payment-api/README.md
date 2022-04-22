# DApp payment gateway



## Payment



### Flow

- Case 1 Transfer and submited transaction success.
  Transfer -> sign -> completed

- Case 2 Transfer failed
  Transfer -> failed

- Case 3 Transfer success but submit failed.
  Transfer -> sign -> suibmit failed -> resubmit -> completed



### Transfer

| Parameters | Required | Comment    | Examples           |
| ---------- | -------- | ---------- | ------------------ |
| network    | Required |            | Eurus/Ethereum     |
| chainId    | Required |            | 1001/1984          |
| marchant   | Required |            | BTCC               |
| tag        | Required |            | 1000123#other_info |
| token      | Optional | User input | USDT               |
| amount     | Optional | User input | Integer 100000     |

http://host_name/{chainName}/{chainId}/send/{merchant}/{tag}?send_amount=100



### Sign



| Prameters | Required | Comment | Examples                                                     |
| --------- | -------- | ------- | ------------------------------------------------------------ |
| tx        | Required |         | 0x93e4e491eb6b808cdb50711c5d05741ecd9a7dc2cd7408bc95cb2f587a9c38a4 |
| network   | Required |         | Eurus/Ethereum                                               |
| chainId   | Required |         | 1001/1984                                                    |
| merchant  | Required |         | BTCC                                                         |
| tag       | Required |         | 1000123#other_info                                           |
|           |          |         |                                                              |
|           |          |         |                                                              |





### Completed
http://host_name/submitTx/

| Prameters | Required | Comment | Examples                                                     |
| --------- | -------- | ------- | ------------------------------------------------------------ |
| tx        | Required |         | 0x93e4e491eb6b808cdb50711c5d05741ecd9a7dc2cd7408bc95cb2f587a9c38a4 |
| network   | Required |         | Eurus/Ethereum                                               |
| chainId   | Required |         | 1001/1984                                                    |
| merchant  | Required |         | BTCC                                                         |
| tag       | Required |         | 1000123#other_info                                           |
| signature | Required |         |                                                              |
|           |          |         |                                                              |



## Transaction Status Enquiry

### Flow

- Case1 Tx found but not submitted
  Transaction status -> sign -> completed

- Case 2 Tx found and submitted
  Transaction status -> completed

- Case 3 Tx not found or fail or invalid to address.
  Transaction status -> failed

### Transaction status

| Prameters | Required | Comment    | Examples                                                     |
| --------- | -------- | ---------- | ------------------------------------------------------------ |
| tx        | Optional | User input | 0x93e4e491eb6b808cdb50711c5d05741ecd9a7dc2cd7408bc95cb2f587a9c38a4 |
| network   |          |            | Eurus/Ethereum                                               |
| chainId   |          |            | 1001/1984                                                    |
| Tag       | Optional |            | 1000123#other_info                                           |

http://host_name/{chainName}/{chainId}/txStatus/{tag}

### query api
| Prameters | Required | Comment    | Examples                                                     |
| --------- | -------- | ---------- | ------------------------------------------------------------ |
| tx        | Required | User input | 0x93e4e491eb6b808cdb50711c5d05741ecd9a7dc2cd7408bc95cb2f587a9c38a4 |
| network   | Required |            | Eurus/Ethereum                                               |
| chainId   | Required |            | 1001/1984                                                    |



### Sign

| Prameters | Required | Comment                           | Examples                                                     |
| --------- | -------- | --------------------------------- | ------------------------------------------------------------ |
| tx        | Required |                                   | 0x93e4e491eb6b808cdb50711c5d05741ecd9a7dc2cd7408bc95cb2f587a9c38a4 |
| network   | Required |                                   | Eurus/Ethereum                                               |
| chainId   | Required |                                   | 1001/1984                                                    |
| merchant  | Required | Retrieved by to address of the Tx | BTCC                                                         |
| tag       | Optional | User input                        | 1000123#other_info                                           |
| token     | Required | Retrieved by tx hash              | USDT                                                         |
| amount    | Required | Retrieved by tx hash              | Integer 100000                                               |



### Completed

### Merchant Enquiry
http://host_name/{merchantCode}

return list of tokens
| Properties   | Comment |      |
| ------------ | ------- | ---- |
| network      |         |      |
| chainId      |         |      |
| token        |         |      |
| walletAddress|         |      |


### Transaction status - by merchant
http://host_name/{merchantCode}/{lastSeqNO}

return list of confirmed transaction
| Properties      | Comment | Examples |
| --------------- | ------- | -------- |
| network         |         |          |
| chainId         |         |          |
| token           |         |          |
| tag             |         |          |
| referenceNo     |         |          |
| amount          |         |          |
| merchantSeqNo   |         |          |
| txSignature     |         |          |
| txHash          |         |          |

## Appendix

### Terminology

#### id, xxx_id
Unique number, use id as first column in every tables and set it to primary key

#### xxx_code
Unique short string, only use a-z, A-Z, 0-9, underscore, hyphen

#### xxx_name
Longer string for display purpose, the text may be changed in future so should not use it as key

#### xxx_hash
If stored in DB as hex string, always start with 0x and use lower case

### Abbreviation

| Word        | Abbreviation |
|-------------|--------------|
| Transaction | tx           |

### Network List

#### t_networks
|  Properties  |     Type     | Examples                        |
|--------------|--------------|---------------------------------|
| id           | bigint       |                                 |
| network_code | varchar(32)  | ETH, ETH_Rinkey                 |
| network_name | varchar(128) | Ethereum, Ethereum Test Network |
| chain_id     | int          | 1, 4                            |
| rpc_url      | text         |                                 |

### Token List

#### t_tokens

| Properties |     Type    |    Examples    |
|------------|-------------|----------------|
| id         | bigint      |                |
| network_id | bigint      |                |
| address    | text        |                |
| symbol     | varchar(16) | ETH, BTC       |
| name       | varchar(64) | Ether, Bitcoin |
| decimals   | int         |                |

### Merchant List

#### t_merchants

|    Properties     |     Type     |               Comment               |             Examples             |
|-------------------|--------------|-------------------------------------|----------------------------------|
| id                | bigint       |                                     |                                  |
| merchant_code     | varchar(32)  |                                     | BTCC                             |
| merchant_name     | varchar(128) |                                     | Bitcoin China                    |
| tag_display_name  | text         | Display name of tag display in DApp | {"en":"Account","zh_tw":"帳號"}   |
| tag_description   | varchar(128) |                                     | Account No. / Address            |
| merchant_last_seq | bigint       |                                     |                                  |

### Merchant - API Key

#### t_merchant_api_keys

| Properties  |  Type  |
|-------------|--------|
| id          | bigint |
| merchant_id | bigint |
| api_key     | text   |

### Merchant - wallet

#### t_merchant_wallets
| Properties  |  Type  | Examples |
|-------------|--------|----------|
| id          | bigint |          |
| merchant_id | bigint |          |
| token_id    | bigint |          |
| address     | text   |          |

### Submission
#### t_submissions
All entries in this table are signature verified

However this does not mean the content must match on chain transaction

`tx_status` and `payment_status` will be updated after transaction verification

|   Properties    |     Type     |                              Examples                                   |
|-----------------|--------------|-------------------------------------------------------------------------|
| id              | bigint       |                                                                         |
| submit_time     | timestamptz  |                                                                         |
| token_id        | bigint       |                                                                         |
| from_address    | text         |                                                                         |
| amount          | numeric      |                                                                         |
| merchant_id     | bigint       |                                                                         |
| tag             | text         |                                                                         |
| tx_hash         | text         |                                                                         |
| tx_status       | int          | -1: default 0: fail 1: success                                          |
| payment_status  | int          | 0: default 1: confirming 2: confirmed 3: rejected 4: not found(timeout) |
| signature       | text         |                                                                         |
| message_body    | text         |                                                                         |

### Transaction
#### t_transactions

All rows in this table should not be changed again, that means all transactions must in one of the following conditions:

- Successful transaction, signature verified, payment information also correct
- A failed transaction
- Transaction is not a valid payment (neither ETH transfer nor ERC20 transfer)

|   Properties    |     Type     |                              Examples                              |
|-----------------|--------------|--------------------------------------------------------------------|
| id              | bigint       |                                                                    |
| submit_time     | timestamptz  |                                                                    |
| confirmed_time  | timestamptz  |                                                                    |
| token_id        | bigint       |                                                                    |
| from_address    | text         |                                                                    |
| amount          | numeric      |                                                                    |
| merchant_id     | bigint       |                                                                    |
| tag             | text         |                                                                    |
| merchant_seq_no | bigint       |                                                                    |
| submission_id   | bigint       |                                                                    |
| onchain_status  | int          | 0:default 1:confirming 2:confirmed 3:rejected 4:not found(timeout) |
| confirm_status  | int          | 0:default 1:confirmed 2:rejected                                   |
| signature       | text         |                                                                    |
| tx_hash         | text         |                                                                    |
| block_hash      | text         |                                                                    |
| block_number    | bigint       |                                                                    |
