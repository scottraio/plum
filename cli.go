package plum

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/scottraio/plum/logger"
	memory "github.com/scottraio/plum/memory"
)

const (
	CmdForget = "/forget"
	CmdTrain  = "/train"
	CmdAgent  = "/agent"
	CmdModel  = "/model"
	CmdPurge  = "/purge"
	CmdMemory = "/memory"
	CmdExit   = "exit"
)

type CliConfig struct {
	DefaultAgent string
}

func Cli(config CliConfig) {
	logger.Print(`
Welcome to the Plum CLI!
------------------------

The Plum CLI Tool allows you to chat with agents and models in the terminal. Useful during development and testing. 

Commands:
  /agent - choose an agent
  /model - choose a model
  /forget - forget everything
  /train - train the current model
  /purge - purge the current model
  /memory - show the current memory
  exit - exit the CLI

To begin type /agent to select an agent or chat with a model 
by typing /model. 

	`, "purple")

	// Create a channel to communicate between the main function and the chat function
	msgChan := make(chan string)

	// Create a new Memory struct
	mem := memory.Memory{}

	// Continuously read user input and send it to the chat function
	reader := bufio.NewReader(os.Stdin)

	// Start the chat function in a new goroutine, passing a pointer to the Memory struct and the msgChan
	go startCli(mem, reader, msgChan, config)

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

func startCli(mem memory.Memory, reader *bufio.Reader, msgChan <-chan string, config CliConfig) {
	var currentAgent string
	var currentModel string
	var currentContext string

	if config.DefaultAgent != "" {
		currentContext = "agent"
		currentAgent = config.DefaultAgent
	}

	cursor(currentContext, currentAgent, currentModel)

	// ctx := context.Background()

	for {
		input := <-msgChan

		if strings.HasPrefix(input, CmdForget) {
			// reset memory
			mem.History = nil

			// clear screen
			fmt.Print("\033[H\033[2J")
			fmt.Print("> ")
			continue

		} else if strings.HasPrefix(input, CmdExit) {
			os.Exit(0)

		} else if strings.HasPrefix(input, CmdTrain) {
			if currentContext == "model" {
				App.Models[currentModel].TrainModel(map[string]string{})
			} else {
				for key, model := range App.Models {
					logger.Log("Training model", key, "orange")
					model.TrainModel(map[string]string{})
				}
			}
			// clear screen
			cursor(currentContext, currentAgent, currentModel)
			continue

		} else if strings.HasPrefix(input, CmdPurge) {
			if currentContext == "model" {
				App.Models[currentModel].Purge()
			} else {
				for key, model := range App.Models {
					logger.Log("Purging model", key, "yellow")
					model.Purge()
				}
			}
			// clear screen
			cursor(currentContext, currentAgent, currentModel)
			continue

		} else if strings.HasPrefix(input, CmdMemory) {
			// purge
			logger.Log("Memory", mem.Format(), "cyan")

			// clear screen
			cursor(currentContext, currentAgent, currentModel)
			continue

		} else if strings.HasPrefix(input, CmdAgent) {
			// reset current agent
			agent := chooseAgent(msgChan)
			currentAgent = agent
			currentContext = "agent"

			// clear screen
			color.Cyan("\n- Agent switched to " + currentAgent + ".\n\n")
			cursor(currentContext, currentAgent, currentModel)
			continue
		} else if strings.HasPrefix(input, CmdModel) {
			// reset current agent
			model := chooseModel(msgChan)
			currentModel = model
			currentContext = "model"

			// clear screen
			color.Yellow("Model switched to " + currentModel + ".\n\n")
			cursor(currentContext, currentAgent, currentModel)
			continue
		}

		if currentContext == "agent" {
			// Run the agent
			agent := App.Agents[currentAgent]

			agent.Remember(mem).Answer(input)

			cursor(currentContext, currentAgent, currentModel)
		} else if currentContext == "model" {

			model := App.Models[currentModel]
			answer := model.Find(input, map[string]string{}, map[string]interface{}{
				"TopK": float64(3),
			})
			logger.Log("[Answer]", answer, "gray")

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
	color.Cyan("\n- Available Agents:")
	for i, name := range agents {
		agentText := fmt.Sprintf("\t%d. %s\n", i+1, name)
		color.Cyan(agentText)
	}

	// Prompt the user to choose an Agent
	for {
		fmt.Print(color.CyanString("\n- Choose an Agent to chat with (enter a number): "))

		input := <-msgChan

		// Parse the user input as an integer
		index, err := strconv.Atoi(input)
		if err != nil || index < 1 || index > len(agents) {
			color.Red("- Invalid input. Please choose an available Agent.")
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
	color.Yellow("\n- Available Models:")
	for i, name := range models {
		color.Yellow(fmt.Sprintf("\t%d. %s\n", i+1, name))
	}

	// Prompt the user to choose a Model
	for {
		fmt.Print(color.YellowString("\n- Choose a Model to use (enter a number): "))

		input := <-msgChan

		// Parse the user input as an integer
		index, err := strconv.Atoi(input)
		if err != nil || index < 1 || index > len(models) {
			color.Red("Invalid input. Please choose an available Model.")
			continue
		}

		return models[index-1]
	}
}

func cursor(context string, agent string, model string) {
	fmt.Println("")

	if context == "agent" {
		cur := color.CyanString(fmt.Sprintf("[%s] > ", agent))
		fmt.Print(cur)
	} else if context == "model" {
		cur := color.YellowString(fmt.Sprintf("[%s] > ", model))
		fmt.Print(cur)
	} else {
		fmt.Print("> ")
	}
}
