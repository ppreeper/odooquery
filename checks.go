package main

import (
	"fmt"
	"strconv"
)

func checkErr(err error) {
	if err != nil {
		fmt.Printf("%v\n", err.Error())
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
