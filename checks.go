package main

import (
	"fmt"
	"os"
	"strconv"
)

func checkErr(err error) {
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
}

func fatalErr(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(2)
	}
}

func IsInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func IsBool(s string) bool {
	_, err := strconv.ParseBool(s)
	return err == nil
}
