package main

import (
	"errors"
	"fmt"
	"github.com/looplab/fsm"
	"strconv"
)

type (
	ATM struct {
		FSM            *fsm.FSM
		TotalAmount    int
		UserAmount     int
		WithdrawAmount int
	}
	Pin struct {
		Code string
	}
	RequestedAmount struct {
		Amount int
	}
	Confirm struct {
		Answer string
	}
)

const (
	PIN = "1234"
	YES = "yes"
	NO  = "no"
)

func NewATM() *ATM {
	atm := &ATM{
		TotalAmount:    10000,
		UserAmount:     20000,
		WithdrawAmount: 0,
	}

	atm.FSM = fsm.NewFSM(
		"start",
		fsm.Events{
			{Name: "InsertCard", Src: []string{"start"}, Dst: "pin"},
			{Name: "RemoveCard", Src: []string{"pin", "main"}, Dst: "start"},
			{Name: "SendPin", Src: []string{"pin"}, Dst: "main"},
			{Name: "GoToBalance", Src: []string{"main"}, Dst: "balance"},
			{Name: "PrintBalance", Src: []string{"balance"}, Dst: "main"},
			{Name: "GoToWithdraw", Src: []string{"main"}, Dst: "withdraw"},
			{Name: "WithdrawMoney", Src: []string{"withdraw"}, Dst: "receipt"},
			{Name: "PrintReceipt", Src: []string{"receipt"}, Dst: "start"},
		},
		fsm.Callbacks{
			"before_SendPin": func(event *fsm.Event) {
				atm.beforeSendPin(event)
			},
			"enter_balance": func(event *fsm.Event) {
				atm.enterBalance(event)
			},
			"before_PrintReceipt": func(event *fsm.Event) {
				atm.printReceipt(event)
			},
			"before_WithdrawMoney": func(event *fsm.Event) {
				atm.beforeWithdrawMoney(event)
			},
		})
	return atm
}

func (atm *ATM) beforeSendPin(e *fsm.Event) {
	// TODO: find correct way to do it: e.Args[0].(Pin).Code
	if e.Args[0].(Pin).Code != PIN {
		e.Cancel(errors.New("wrong password"))
	}
}

func (atm *ATM) enterBalance(_ *fsm.Event) {
	fmt.Println(atm.UserAmount)
}

func (atm *ATM) printReceipt(e *fsm.Event) {
	answer := e.Args[0].(Confirm).Answer
	if answer == YES {
		fmt.Println(atm.WithdrawAmount)
	} else if answer == NO {

	}
}

func (atm *ATM) beforeWithdrawMoney(e *fsm.Event) {
	requestedAmount := e.Args[0].(RequestedAmount).Amount
	if atm.UserAmount < requestedAmount {
		e.Cancel(errors.New("not enough money on your account. you can withdraw up to: " + strconv.Itoa(atm.UserAmount)))
	} else if atm.TotalAmount < requestedAmount {
		e.Cancel(errors.New("not enough money on machine. you can withdraw up to: " + strconv.Itoa(atm.TotalAmount)))
	} else {
		atm.UserAmount -= requestedAmount
		atm.TotalAmount -= requestedAmount
		atm.WithdrawAmount += requestedAmount
	}
}

func main() {
	atm := NewATM()

	_ = atm.FSM.Event("InsertCard")
	_ = atm.FSM.Event("SendPin", Pin{Code: "1234"})
	_ = atm.FSM.Event("GoToBalance")
	_ = atm.FSM.Event("PrintBalance")
	_ = atm.FSM.Event("GoToWithdraw")
	_ = atm.FSM.Event("WithdrawMoney", RequestedAmount{Amount: 100000})
	_ = atm.FSM.Event("WithdrawMoney", RequestedAmount{Amount: 5000})
	_ = atm.FSM.Event("PrintReceipt", Confirm{Answer: YES})

}
