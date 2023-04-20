package framework

type Tool struct {
	Name        string
	Description string
	Func        func(string) string
}

func UseTool(name string, desc string, run func(string) string) Tool {
	return Tool{
		Name:        name,
		Description: desc,
		Func:        run,
	}
}
