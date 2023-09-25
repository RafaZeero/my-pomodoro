package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	for {

		Welcome()
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			Choices()
			text := strings.ToUpper(scanner.Text())
			switch text {
			case "1":
				fmt.Println("Starting 30 min ...")
				notification, err := Notify("30 min pomodoro", "Starting 30 min !! 25 minutes focused + 5 minutes break")
				if err != nil {
					panic(err)
				}
				fmt.Println(notification)
			case "2":
				fmt.Println("Starting 45 min ...")
			case "3":
				fmt.Println("Starting 60 min ...")
			case "Q":
				fmt.Println("Quitting zeero pomodoro ...")
				time.Sleep(time.Second)
				os.Exit(0)

			}
		}
	}
}

func Welcome() {
	fmt.Println(`
============================================
	WELCOME TO ZEERO POMODORO !!
============================================

Press enter to start
	`)
}

func Choices() {
	fmt.Println("Please choose one of the following:")
	fmt.Println("(1) - Start 30 min pomodoro")
	fmt.Println("(2) - Start 45 min pomodoro")
	fmt.Println("(3) - Start 60 min pomodoro")
	fmt.Println("(q/Q) - Quit")
}

func Notify(summary, body string) (string, error) {
	cmd := exec.Command("notify-send", "-t", "3000", summary, body)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
