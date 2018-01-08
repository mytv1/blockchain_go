# blockchain_go
My simple implement of blockchain with Golang.

Fork from https://github.com/Jeiwan/blockchain_go
(Many thanks to Jeiwan!)

This is part 1 of my articles about my blockchain's implement tutorial below :

1. [Basic prototype](https://github.com/mytv1/blockchain_go/tree/part_1)
2. [Network](https://github.com/mytv1/blockchain_go/tree/part_2)

I'm not good at English. So forgive me if something make you hard to understand.

I'm also new in Golang and Blockchain. So you can spot a lot of problems in my code, please feel free to correct it.

# Contents
- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Running](#running)
- [Program Structure](#structure)
- [References](#references)

# Introduction
In this article, i'll build a simplified blockchain and block structure.
When you run the program, the sample chain of blocks will be printed with it's hash and some information.

# Prerequisites
(My local enviroment)

+ OS : Ubuntu 16.04.2 LTS

+ Golang :
```
$ go version
go version go1.9.2 linux/amd64
```
# Running
```
go build .
./blockchain.go
```

# Structure
Basic structure :

### Block
```
type Block struct {
	Hash          []byte
	PrevBlockHash []byte
	Data          []byte
	Timestamp     int64
}
```

Block is the basic element of blockchain. In this article i just built it with very simple properties.
In practice, you can consider ethereum block structure [here](https://github.com/ethereum/go-ethereum/blob/master/core/types/block.go#L139)

In my struct, we have :
- Timestamp : timestamp when block is created. I use int64 type because it can present and consistant with Unix time.

- Data : for simple usage, it's type may be string or something else according to purpose. But with []byte, we can handle it more easily when we need encryption or json marshal.

- Hash : We need a hash to "seal" this block. Hash is special mechanism in blockchain with many purpose : hard to fake unique block, "proof of work", easy validation... . We will mention and change it that later. In this artice, hash is simple and just `sha256(Data + PrevBlockHash + Timestamp)`. About hash type, with same reason above i think []byte is the first best, second best may be string.

- PrevBlockHash : Blockchain is a linked list of blocks. So each block will link to previous (created previously) block by saving a hash to it (it's pointer in below diagram)
![linked list](https://s3-eu-west-2.amazonaws.com/dotjsonimages/2017/06/ll-4.png)

The first block have no pointer to it's previous block, and called by **Genesis block**

Expanding: If you ask me why we shouldn't add NextBlockHash (Seem like it will help we iterate blocks easily! So blockchain will transform form linkly linked list to doubly linked list), i will answer that because stablility. As a block view, previous block is stabled than next block. With some special condition, block may be changed, and potential that it's is next block higher than previous's. And i think we should save to blockchain stable block data as much as we can.

### Blockchain
```
type Blockchain struct {
	blocks []*Block
}
```

Blockchain structure is just array of Blocks, very strange. With purpose of this article: simple blockchain simulation, i think it's enough. Although we can see at least 2 problems here:
+ Memory storage : Currently bitcoin blockchain size is about 150GB, and we don't want to save something like that to memory. In fact, i think it's saved to file (database). We may implement it later.
+ PrevBlockHash wasted : With this structure, we can reference to previous block by indexing, therefore block's PrevBlockHash is wasted a bit. But i think it's fundamental, and we can see it's more useful on future articles.


Enjoy!

# References
https://jeiwan.cc/posts/building-blockchain-in-go-part-1/