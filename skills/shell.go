package skills

type Shell struct {
}

func (shell *Shell) Execute(cmd string) string {
	return cmd
	// exec := exec.Command(cmd)
	// output, err := exec.Output()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return ""
	// }

	// return string(output)
}
