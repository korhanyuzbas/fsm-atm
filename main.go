package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/looplab/fsm"
	"os"
	"strings"
)

type ATM struct {
	To  string
	FSM *fsm.FSM
}

func NewATM(to string) *ATM {
	atm := &ATM{
		To: to,
	}

	atm.FSM = fsm.NewFSM(
		"start",
		fsm.Events{
			{Name: "insert_card", Src: []string{"start"}, Dst: "pin"},
			{Name: "remove_card", Src: []string{"pin", "main"}, Dst: "start"},
			{Name: "send_pin", Src: []string{"pin"}, Dst: "main"},
			{Name: "check_balance", Src: []string{"main"}, Dst: "balance"},
			{Name: "return_main", Src: []string{"balance"}, Dst: "main"},
		},
		fsm.Callbacks{
			"before_send_pin": func(event *fsm.Event) {
				atm.beforeSendPin(event)
			},
			"enter_main": func(event *fsm.Event) {
				atm.enterMain(event)
			},
			"enter_balance": func(event *fsm.Event) {
				atm.enterBalance(event)
			},
		})
	return atm
}

func (atm *ATM) beforeSendPin(e *fsm.Event) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter PIN: ")
	scanner.Scan()
	pin := scanner.Text()
	if pin != strings.TrimSpace("1234") {
		e.Cancel(errors.New("wrong password"))
	}
}

func (atm *ATM) enterMain(e *fsm.Event) {
	fmt.Println(e.FSM.AvailableTransitions())
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("welcome to main page. what do you want to do? ( (b)alance, (w)ithdraw, (e)xit ):")
	scanner.Scan()
	req := scanner.Text()
	switch req {
	case "b":
		e.FSM.Event("check_balance")
	case "e":
		e.FSM.Event("remove_card")
	}
}

func (atm *ATM) enterBalance(e *fsm.Event) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("welcome to balance page. what do you want to do? ( (p)rint, (c)ancel ):")
	scanner.Scan()
	req := scanner.Text()
	switch req {
	case "p":
		atm.FSM.Event("print_balance")
	case "c":
		atm.FSM.SetState("main")
	}
}

func main() {
	atm := NewATM("enpara")

	err := atm.FSM.Event("insert_card")
	if err != nil {
		fmt.Println(err)
	}

	err = atm.FSM.Event("send_pin")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(atm.FSM.Current())

}
