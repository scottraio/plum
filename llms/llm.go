package llms

type LLM interface {
	Client() LLM
	Run(prompt string) string
}
