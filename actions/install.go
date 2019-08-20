package actions

import (
	"fmt"

	"gopkg.in/AlecAivazis/survey.v1"
)

func StartInstall() {
	fmt.Println("Starting wizard")
	askFirstQuestions()
}

type answers struct {
	RootDomain string
}

func askFirstQuestions() {
	var questions = []*survey.Question{
		{
			Name:     "RootDomain",
			Prompt:   &survey.Input{Message: "Base URL:"},
			Validate: survey.Required,
		},
	}

	a := &answers{}

	if err := survey.Ask(questions, a); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("it worked!")
}
