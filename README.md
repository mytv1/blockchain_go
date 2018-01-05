# blockchain_go
My simple implement of blockchain with Golang.

This is part 2 of my articles about my blockchain's implement tutorial below :

1. [Basic prototype](https://github.com/mytv1/blockchain_go/tree/part_1)
2. [Network](https://github.com/mytv1/blockchain_go/tree/part_2)

I'm not good at English. So please tell me if there is something make you hard to understand.

I'm also new in Golang and Blockchain. So if you spot any problem in my code, please feel free to correct it.

# Contents
- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Install](#install)
- [Running](#running)
- [Program Structure](#program-structure)
- [References](#references)

# Introduction
In this article, we'll built a blockchain with very simple decentralized network.

Here are some things you can try with this part :
+ Run n node with shared blockchain
+ First started node will initialize its first block. Other nodes, when started, will pull blockchain from first connected node in network
+ Send command in json format to a node from command line by tcp. Here are two avaiable commands : "Add block" and "Print entire blockchain"
+ When a block is added in a node, it will be shared to others immediately. Therefore all blockchain in all node will be synchronized

*Note : i implemented it with my basic knowledge about blockchain. So i'm not sure i did it properly. :)

# Prerequisites
(My local environment)

- OS : Ubuntu 16.04.2 LTS
- Go install : https://golang.org/doc/install

```
$ go version
go version go1.9.2 linux/amd64
```

# Install
```
make deps
make build
```

# Running
## Prepare 
To make it works like a network, we need run each node independently. In example, you can prepare 3 nodes network's environment like this :

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

* simplebc : executed file. May be you've built it with `make build` on install section above
* config.json : information about your network.
* samples (optional) : contains commands to a node. You can send command to a node by tcp to request it add a block, or print its own blockchain as you saw in part 1

config.json on each node will look like below :

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
  }
}
...
{
  "network": {
    "local_node": {
      "address": "localhost:3332"
    },
    "neighbor_nodes": [
      {"address": "localhost:3331"},
      {"address": "localhost:3333"}
    ]
  }
}
...
{
  "network": {
    "local_node": {
      "address": "localhost:3333"
    },
    "neighbor_nodes": [
      {"address": "localhost:3331"},
      {"address": "localhost:3332"}
    ]
  }
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

## Send command to a node
```shell
# print its own blockchain
cat samples/print.json | nc localhost {node_port}

# add block with your own data (bytes type) to it blockchain
cat samples/addblock.json | nc localhost {node_port}
```

# Program Structure

## Node lifecycle
Node's lifecycle is described by flowchart below:

![flowchart](https://i.imgur.com/F1m7SCf.jpg)

## Source structure
* `block.go, blockchain.go` : no changes compare to part_1
* `cli.go` : contains funtions that help us integrate with program by command line. I just implement only `start` command here. In detail, you can type `help` to show all avaiable commands provided by program.
* `config.go` : contains structs and functions to store and manipulate our configuration information. In this part, only network's configuration is stored.
* `config.json` : contains our configuration information.
* `log.go` : contains my customized log mechanism. With this, we can write log easier. In this part, node will log its received messages and its state.
* `message.go` : nodes communicate to each others with tcp protocol, and its message is defined here.
* `network.go` : contains structs and functions to manipulate our network. Message sprearding and sending is defined here.
* `server.go` : contains functions that listen and handle received messages from other nodes.

# References
https://github.com/DNAProject/DNA
https://github.com/urfave/cli
