# blockchain_go
My simple implement of blockchain with Golang.

Fork from https://github.com/Jeiwan/blockchain_go
(Many thanks to Jeiwan!)

This is part 5 of my articles about my blockchain's implement tutorial below :

1. [Basic prototype](https://github.com/mytv1/blockchain_go/tree/part_1)
2. [Network](https://github.com/mytv1/blockchain_go/tree/part_2)
3. [Proof of work](https://github.com/mytv1/blockchain_go/tree/part_3)
4. [Wallet + Address](https://github.com/mytv1/blockchain_go/tree/part_4)
5. [Transaction](https://github.com/mytv1/blockchain_go/tree/part_5)

I'm not good at English. So please tell me if there is something make you hard to understand.

I'm also new in Golang and Blockchain. So if you spot any problem in my code, please feel free to correct it.

Demo : https://youtu.be/X8G33BZS3WY

Describe by vietnamese:
P1 : https://kipalog.com/posts/Xay-dung-blockchain-don-gian-voi-golang--P1---Cau-truc-co-ban
P2 : https://kipalog.com/posts/Xay-dung-blockchain-don-gian-voi-golang--P2---CLI---Network
P3 : https://kipalog.com/posts/Xay-dung-blockchain-don-gian-voi-golang--P3---Persistent---Proof-Of-Work
P4 : https://kipalog.com/posts/Xay-dung-blockchain-don-gian-voi-golang--P4---Wallet--Address
P5 : https://kipalog.com/posts/Xay-dung-blockchain-don-gian-voi-golang--P5----Transaction---UTOX-Set--Phan-cuoi

# Contents
- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Running](#running)

# Introduction
In this article, we'll build a blockchain with a simple structure

When you run the program, the sample chain of blocks will be printed with its hash and its information.

# Prerequisites
(My local environment)

+ OS : Ubuntu 16.04.2 LTS

+ Golang :
```
$ go version
go version go1.9.2 linux/amd64
```
# Running
```
$ make deps
$ make build
$ ./simplebc s
```

## Running introduction for functions testing
To make it works like a network, we need to run each node independently. For example, you can prepare a network environment with 3 nodes like this :

```shell
$tree
# directory structure #
├── node_1
│   ├── config.json
│   ├── simplebc
│   └── samples (optional)
│       ├── print.json
│       └── addblock.json
├── node_2
│   ├── config.json
│   ├── simplebc
│   └── samples (optional)
│       ├── print.json
│       └── addblock.json
└── node_3
    ├── config.json
    ├── simplebc
    └── samples (optional)
        ├── print.json
        └── addblock.json

```

* simplebc : executed file. You can build it with `make build` above.
* config.json : information about your network.
* samples (optional) : contains commands to a node. You can send command to a node via tcp to request it to add a block, or print its own blockchain as you saw in part 1

### Create a wallet
```shell
# node_1,2,3
./simplebc cw
# Copy printed wallets's info to below config.json
```

`config.json` of each node will look like below :

```shell
# configuration snippets #
$ cat node_[123]/config.json
{
  "network": {
    "local_node": {
      "address": "localhost:3331"
    },
    "neighbor_nodes": [
      {"address": "localhost:3332"},
      {"address": "localhost:3333"}
    ]
  },
  "wallet": {
    "private_key": "d1ac80357c748483ac6952f09f56b1465bccfdb398262e774adfbfea7ee57331",
    "public_key": "0490dcdf24d98bba9c0da3b56e4862fd38811e05a809deab0bb300761cc51615f0b2f05f46294ae6b313f08b6d074d900efb56f80ef8c3e131119b1e800e47c0fd",
    "address": "1GiEnhrStofDdzyKpnDw7BpAcv3m9vacFv"
  }
}
...
  "network": {
    "local_node": {
      "address": "localhost:3332"
    },
    "neighbor_nodes": [
      {"address": "localhost:3331"},
      {"address": "localhost:3333"}
    ]
  },
  "wallet": {
    "private_key": "1b9205c8f84402cb32039de762e6ea541068ccbbb19e8b978b21a4e40197a92c",
    "public_key": "04040d8408016f873529f9a99102a946a4d36de80360e4b57ada4afd6e06cf8a9ea399aacea0d18605bba464d2d888adbdc0b6b9dd94a4f7e50a98b94ed907b2ac",
    "address": "1NYJwQd8YuR8kFHxs5WPRbuMtFQAyX4DEz"
  }
...
  "network": {
    "local_node": {
      "address": "localhost:3333"
    },
    "neighbor_nodes": [
      {"address": "localhost:3331"},
      {"address": "localhost:3332"}
    ]
  },
  "wallet": {
    "private_key": "291c5925c8dffae53f8f8664222cf41553786520d4a4d4a98396a911e911a056",
    "public_key": "04d1daffec2de150ab808200019d116ada5cb8942546a0e50d8a432080efa5deee3a9027267c2f5e93ef54d8a4948941c6a6cdf48aeee5b59be4f157e3f1ec2f43",
    "address": "1PQE324cGzr9GnNDpmGMybLN16WusiYyF1"
  }
```

## Start network

```shell
# node_1
./simplebc start

# node_2
./simplebc start

# node_3
./simplebc start
```

## Create transaction
```shell
# Create transaction, node1
./simplebc ct -to {address} -v {value} -f tx.json
# Send transaction to node
cat tx.json | nc localhost {port}
```

## Vardump blockchain
```shell
cat print.json | nc localhost {port}
```

## Print all address's value in blockchain
```shell
cat all_addr.json | nc localhost {port}
```

# Reference
https://jeiwan.cc/posts/building-blockchain-in-go-part-1/
https://github.com/DNAProject/DNA
