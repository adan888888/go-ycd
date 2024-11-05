package controllers

import (
	"exchangeapp/utils"
	"fmt"
	"testing"
)

func TestGetUid(t *testing.T) {
	uid := utils.GetUid()
	//fmt.Printf("id : %d\n", uid)
	//
	//time.Sleep(time.Second)
	uid = utils.GetUid()
	fmt.Printf("id : %d\n", uid)

	//time.Sleep(time.Second)
	uid = utils.GetUid()
	fmt.Printf("id : %d\n", uid)

	//time.Sleep(time.Second)
	uid = utils.GetUid()
	fmt.Printf("id : %d\n", uid)

}
func TestGetUid1(t *testing.T) {
	uid := utils.GetUid()
	fmt.Printf("id : %d\n", uid)

	uid = utils.GetUid()
	fmt.Printf("id : %d\n", uid)

	uid = utils.GetUid()
	fmt.Printf("id : %d\n", uid)

	uid = utils.GetUid()
	fmt.Printf("id : %d\n", uid)

	var input struct {
		Username string
		Password string
	}
	fmt.Printf("uid : %v,%v\n", len(input.Username), input.Username == "")
}
