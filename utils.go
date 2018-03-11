package main

import (
	"strconv"
)

func intToBytes(num int) []byte {
	s := strconv.Itoa(num)

	return []byte(s)
}

func bytesToInt(data []byte) int {
	i, err := strconv.Atoi(string(data))
	if err != nil {
		Error.Panic(err)
		return 0
	}

	return i
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
