# blockchain_go
My simple implement of blockchain with Golang.

This is the second part of my articles about my blockchain's implement tutorial below :

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
In this part, we'll build a simple decentralized network with a blockchain.

Here are a few things you can try with this part :
+ Run n nodes which have the same blockchain
+ Send command in json format to a node from command line via tcp. Here are two available commands : "Add block" and "Print entire blockchain"

Starting the first node will initialize its first block. Other nodes, when started, will try to connect to each other and pull blockchain from the first one it connect in the network. When a block is added into a node, it will be shared to others immediately. Therefore all blockchain in all nodes will be synchronized.



*Note : i implemented it with my basic knowledge about blockchain. So i'm not sure i did it properly. :)

[comment]: <> (Xem xét bỏ prerequisites vì là part 2 rồi. Không thỉ trỏ về part 1 để câu view :v)
# Prerequisites
(My local environment)

- OS : Ubuntu 16.04.2 LTS
- Go install : https://golang.org/doc/install

```
$ go version
go version go1.9.2 linux/amd64
```


# Running
## Prepare 

Build :
```
make deps
make build
```

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
A node's lifecycle is described by this flowchart:

![flowchart](https://i.imgur.com/F1m7SCf.jpg)

## Source structure
* `block.go, blockchain.go` : the same as part_1
* `cli.go` : contains functions that help us interact with our program using command line. I just implement only `start` command here. For more details, you can type `help` to show all available commands.
* `config.go` : contains structs and functions to store and manipulate our configuration information. In this part, only network configuration is stored.
* `config.json` : contains our configuration information.
* `log.go` : my customized log mechanism to write log easier. In this part, node will log its received messages and its state.
* `message.go` : ~~nodes communicate to each others with tcp protocol, and its message is defined here~~.
[comment]: <> (defines messages used to communicate between nodes via tcp)
* `network.go` : contains structs and functions to manipulate our network. Define how to spread and send messages.
* `server.go` : contains functions that listen and handle received messages from other nodes.

# References
* https://github.com/DNAProject/DNA : Other blockchain-based decentralized network implement by go

* https://github.com/urfave/cli : Powerful cli support package
