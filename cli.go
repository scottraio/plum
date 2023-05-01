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
	CmdAgent  = "/agent"
	CmdModel  = "/model"
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
	var currentModel string
	var currentContext string

	for {
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

			if currentContext == "model" {
				trainConfig.Train(currentModel)
			} else {
				trainConfig.TrainAll()
			}

			// clear screen
			cursor(currentContext, currentAgent, currentModel)
			continue

		} else if strings.HasPrefix(input, CmdPurge) {
			// purge
			config.Purge()

			// clear screen
			cursor(currentContext, currentAgent, currentModel)
			continue

		} else if strings.HasPrefix(input, CmdMemory) {
			// purge
			fmt.Println("%v", memory)

			// clear screen
			cursor(currentContext, currentAgent, currentModel)
			continue

		} else if strings.HasPrefix(input, CmdAgent) {
			// reset current agent
			agent := chooseAgent(msgChan)
			currentAgent = agent
			currentContext = "agent"

			// clear screen
			fmt.Println("Agent switched. Choose a new Agent to chat with:")
			cursor(currentContext, currentAgent, currentModel)
			continue
		} else if strings.HasPrefix(input, CmdModel) {
			// reset current agent
			model := chooseModel(msgChan)
			currentModel = model
			currentContext = "model"

			// clear screen
			fmt.Println("Model switched. Choose a new Model to query:")
			cursor(currentContext, currentAgent, currentModel)
			continue
		}

		if currentContext == "agent" {
			agent := App.Agents[currentAgent]
			answer := agent.Run(input, memory)
			history := ChatHistory{Query: input, Answer: answer}
			memory.History = append(memory.History, history)

			fmt.Printf("[%s] > ", currentAgent)
		} else if currentContext == "model" {

			model := App.Models[currentModel]
			answer := model.Find(input, map[string]string{}, map[string]interface{}{})
			App.Log("[Answer]", answer, "gray")

			cursor(currentContext, currentAgent, currentModel)
		} else {
			cursor(currentContext, currentAgent, currentModel)
		}

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

func chooseModel(msgChan <-chan string) string {
	// Create a list of available Models
	var models []string
	for name := range App.Models {
		models = append(models, name)
	}

	// Print the list of available Models
	fmt.Println("Available Models:")
	for i, name := range models {
		fmt.Printf("\t%d. %s\n", i+1, name)
	}

	// Prompt the user to choose a Model
	for {
		fmt.Print("Choose a Model to use (enter a number): ")

		input := <-msgChan

		// Parse the user input as an integer
		index, err := strconv.Atoi(input)
		if err != nil || index < 1 || index > len(models) {
			fmt.Println("Invalid input. Please choose an available Model.")
			continue
		}

		return models[index-1]
	}
}

func cursor(context string, agent string, model string) {
	if context == "agent" {
		fmt.Printf("[%s] > ", agent)
	} else if context == "model" {
		fmt.Printf("[%s] > ", model)
	} else {
		fmt.Print("> ")
	}
}
