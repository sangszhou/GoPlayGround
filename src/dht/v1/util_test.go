package v1

import (
	"testing"
	"log"
	"math/big"
)

func TestGetAddress(t *testing.T) {
	address := GetAddress()
	log.Print("[" + address + "]")
}

func TestHash(t *testing.T) {
	hash := Hash("hello world")
	// 第二个参数是 base
	expected, _ := big.NewInt(0).SetString("243667368468580896692010249115860146898325751533", 10)
	if hash.Cmp(expected) != 0 {
		t.Error("not equal")
	}
}


