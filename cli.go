package plum

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	CmdForget = "/forget"
	CmdTrain  = "/train"
	CmdSwitch = "/switch agent"
	CmdPurge  = "/purge"
	CmdMemory = "/memory"
)

type CliConfig struct {
	BeforeTrain     func(config *TrainConfig)
	Purge           func()
	ModelAttributes map[string]string
	Train           func(config *TrainConfig)
}

func Cli(config CliConfig) {
	fmt.Print("> ")
	// Create a channel to communicate between the main function and the chat function
	msgChan := make(chan string)

	// Create a new Memory struct
	memory := &Memory{}

	// Continuously read user input and send it to the chat function
	reader := bufio.NewReader(os.Stdin)

	// Start the chat function in a new goroutine, passing a pointer to the Memory struct and the msgChan
	go chat(memory, reader, msgChan, config)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}
		msg = strings.TrimSuffix(msg, "\n")
		msgChan <- msg
	}
}

func chat(memory *Memory, reader *bufio.Reader, msgChan <-chan string, config CliConfig) {
	var currentAgent string
	var agent Agent

	for {
		if currentAgent == "" {
			currentAgent = chooseAgent(msgChan)
			fmt.Println("You are now chatting with", currentAgent)
			fmt.Print("> ")
		}

		input := <-msgChan

		if strings.HasPrefix(input, CmdForget) {
			// reset memory
			memory.History = nil

			// clear screen
			fmt.Print("\033[H\033[2J")
			fmt.Print("> ")
			continue

		} else if strings.HasPrefix(input, CmdTrain) {
			trainConfig := &TrainConfig{
				ModelAttributes: config.ModelAttributes,
			}

			config.Train(trainConfig)

			// clear screen
			fmt.Print("> ")
			continue

		} else if strings.HasPrefix(input, CmdPurge) {
			// purge
			config.Purge()

			// clear screen
			fmt.Print("> ")
			continue

		} else if strings.HasPrefix(input, CmdMemory) {
			// purge
			fmt.Println("%v", memory)

			// clear screen
			fmt.Print("> ")
			continue

		} else if strings.HasPrefix(input, CmdSwitch) {
			// reset current agent
			currentAgent = ""

			// clear screen
			fmt.Println("Agent switched. Choose a new Agent to chat with:")
			fmt.Print("> ")
			continue
		}

		agent = App.Agents[currentAgent]
		answer := agent.Run(input, memory)
		history := ChatHistory{Query: input, Answer: answer}
		memory.History = append(memory.History, history)

		fmt.Print("> ")
	}
}

func chooseAgent(msgChan <-chan string) string {
	// Create a list of available Agents
	var agents []string
	for name := range App.Agents {
		agents = append(agents, name)
	}

	// Print the list of available Agents
	fmt.Println("Available Agents:")
	for i, name := range agents {
		fmt.Printf("\t%d. %s\n", i+1, name)
	}

	// Prompt the user to choose an Agent
	for {
		fmt.Print("Choose an Agent to chat with (enter a number): ")

		input := <-msgChan

		// Parse the user input as an integer
		index, err := strconv.Atoi(input)
		if err != nil || index < 1 || index > len(agents) {
			fmt.Println("Invalid input. Please choose an available Agent.")
			continue
		}

		return agents[index-1]
	}
}
