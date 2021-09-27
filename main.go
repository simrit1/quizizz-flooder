package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/manifoldco/promptui"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Println("Loading bots...")
	bots, err := GetBots(10)

	if err != nil {
		log.Fatal("Error while getting bots")
	}

	fmt.Printf("Succesfully loaded %v bots \n", len(bots.Data.Quizzes))

	pinPrompt := promptui.Prompt{
		Label: "Pin",
		Validate: func(input string) error {
			if len(input) < 5 {
				return errors.New("pin must have more than 5 characters")
			}
			if _, err := strconv.Atoi(input); err != nil {
				return errors.New("pin must be a number")
			}
			return nil
		},
	}
	pin, err := pinPrompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	hash, err := GetRoomHash(pin)
	if err != nil {
		log.Fatal("invalid pin")
	}

	fmt.Println("Your room hash " + hash)

	prompt := promptui.Select{
		Label: "Select Mode",
		Items: []string{"Use bots real names", "Use custom names"},
	}

	result, _, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	var customName string
	if result == 1 {
		customNamePrompt := promptui.Prompt{
			Label: "Your custom name",
		}

		customName, err = customNamePrompt.Run()
		if err != nil {
			log.Fatal(err)
		}

	}

	cfg := Config{
		RoomHash:   hash,
		Delay:      1000,
		Mode:       result,
		CustomName: customName,
	}

	go Spam(bots, cfg)

	<-interrupt
}
