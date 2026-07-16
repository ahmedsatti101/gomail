package main

import "fmt"

func check(err error) {
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}
}
