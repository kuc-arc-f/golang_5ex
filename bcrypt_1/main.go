package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "test1234"

	// hash化
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("hash=" + string(hash))

	var chk_password = password;
	//var chk_password = "test999";
	// 照合
	err = bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(chk_password),
	)

	if err != nil {
		fmt.Println("NG")
		return
	}

	fmt.Println("OK")	
}