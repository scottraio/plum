package plum

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	CmdForget = "/forget"
	CmdTrain  = "/train"
	CmdSwitch = "/switch agent"
)

func Cli() {
	fmt.Print("> ")
	// Create a channel to communicate between the main function and the chat function
	msgChan := make(chan string)

	// Create a new Memory struct
	memory := &Memory{}

	//	 Continuously read user input and send it to the chat function
	reader := bufio.NewReader(os.Stdin)

	// Start the chat function in a new goroutine, passing a pointer to the Memory struct
	go chat(msgChan, memory, reader)

	for {
		msg, _ := reader.ReadString('\n')
		msgChan <- msg
	}
}

func chat(msgChan chan string, memory *Memory, reader *bufio.Reader) {
	var currentAgent string

	for {
		if currentAgent == "" {
			currentAgent = chooseAgent(reader)
		}

		input := <-msgChan

		if strings.HasPrefix(input, CmdForget) {
			memory.History = nil
			fmt.Print("\033[H\033[2J")
		} else if strings.HasPrefix(input, CmdTrain) {
			for _, job := range App.Jobs {
				job.Run()
			}
		} else if strings.HasPrefix(input, CmdSwitch) {
			currentAgent = nil
			fmt.Println("Agent switched. Choose a new Agent to chat with:")
		} else {
			answer := currentAgent.Chat(input, *memory, Version)
			history := ChatHistory{Query: input, Answer: answer}
			memory.History = append(memory.History, history)
			fmt.Println(answer)
		}

		fmt.Print("> ")
	}
}

func chooseAgent(reader *bufio.Reader) plum.Agent {
	// Print the list of available Agents
	fmt.Println("Available Agents:")
	for name := range plum.App.Agents {
		fmt.Println("- " + name)
	}

	// Prompt the user to choose an Agent
	for {
		fmt.Print("Choose an Agent to chat with: ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		// Retrieve the corresponding Agent from the App.Agents map
		agent, ok := plum.App.Agents[name]
		if ok {
			return agent
		}

		fmt.Println("Invalid Agent name. Please try again.")
	}
}
