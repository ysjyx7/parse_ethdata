package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
)

func main() {
	dbPath := os.Args[1]
	ancientPath := dbPath + "/ancient"
	upNum, err := strconv.Atoi(os.Args[2])
	endNum, err := strconv.Atoi(os.Args[3])

	ancientDb, err := rawdb.NewLevelDBDatabaseWithFreezer(dbPath, 16, 1, ancientPath, "", true)
	if err != nil {
		panic(err)
	}

	// ReadHeadHeaderHash retrieves the hash of the current canonical head header.
	currHeader := rawdb.ReadHeadHeaderHash(ancientDb)
	fmt.Printf("currHeader: %x\n", currHeader)

	// ReadHeaderNumber returns the header number assigned to a hash.
	currHiehgt := rawdb.ReadHeaderNumber(ancientDb, currHeader)
	fmt.Printf("currHiehgt: %d\n", currHiehgt)

	fmt.Println("----------------------------------------------------------------")

	for i := upNum; i <= endNum; i++ {
		// ReadCanonicalHash retrieves the hash assigned to a canonical block number.
		blkHash := rawdb.ReadCanonicalHash(ancientDb, uint64(i))

		// hash := rawdb.ReadAllHashes(ancientDb, uint64(i))

		fmt.Printf("etherscan url: https://etherscan.io/block/%v\n", i)

		if blkHash == (common.Hash{}) {
			fmt.Printf("i: %v\n", i)
		} else {
			fmt.Printf("blkHash: %x\n", blkHash)
		}

		// ReadBody retrieves the block body corresponding to the hash.
		blkHeader := rawdb.ReadHeader(ancientDb, blkHash, uint64(i))
		fmt.Printf("blkHeader Coinbase: 0x%x\n", blkHeader.Coinbase)
		fmt.Printf("blkHeader Time: %d\n", blkHeader.Time)

		// ReadBody retrieves the block body corresponding to the hash.
		blkBody := rawdb.ReadBody(ancientDb, blkHash, uint64(i))
		fmt.Printf("blkBody: %v\n", blkBody)
		fmt.Printf("blkBody Uncles size: %x\n", len(blkBody.Uncles))
		for _, uncle := range blkBody.Uncles {
			fmt.Printf("uncle Hash: 0x%x\n", uncle.Hash())
		}

		fmt.Printf("blkBody Tx size: %x\n", len(blkBody.Transactions))
		for _, tx := range blkBody.Transactions {
			fmt.Println("tx Hash: 0x%x\n", tx.Hash())
			fmt.Println("tx from addr: 0x%x\n", getFromAddr(tx))
			fmt.Println("tx To: 0x%x\n", tx.To())
		}

		// ReadBlock retrieves an entire block corresponding to the hash
		block := rawdb.ReadBlock(ancientDb, blkHash, uint64(i))
		fmt.Printf("block hash: 0x%x\n", block.Hash())

		fmt.Println("----------------------------------------------------------------")
	}

}

func getFromAddr(tx *types.Transaction) common.Address {
	var signer types.Signer = types.FrontierSigner{} // 这是 Frontier阶段的签名

	from, err := types.Sender(signer, tx)
	if err != nil {
		panic(err)
	}

	return from
}
