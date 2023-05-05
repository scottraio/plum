package skills

import (
	skills "github.com/scottraio/plum/skills"
)

type ShellCommandSkill struct {
	skills.Skill
	skills.Shell
}

func ShellCommand() *skills.Skill {
	var shell *ShellCommandSkill
	// create the model
	shell = &ShellCommandSkill{
		// Model is the base model that you want to use
		Skill: skills.Skill{
			HowTo: `
				You are given shell commands
			`,

			Return: func(query string) string {
				return shell.Execute(query)
			},
		},
	}

	return &shell.Skill
}
