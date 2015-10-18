package main

import
(
	_ "fmt"
)



func encrypt(data []byte, key []byte) []byte{
	cypher := make([]byte, len(data))

	keyLen := len(key)

	var IV byte

	for i:=0;i<len(data) ;i++ {
		cypher[i] = data[i] ^ key[i%keyLen] ^ IV
		IV = cypher[i]
	}

	return cypher
}

func decrypt(cypher []byte, key []byte) []byte{
	data := make([]byte, len(cypher))
	
	keyLen := len(key)

	var IV byte

	for i:=0;i<len(cypher);i++ {
		data[i] = cypher[i] ^ key[i%keyLen] ^ IV
		IV = cypher[i]
	}
	return data
}

