package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

func main() {
	var authType string

	err := survey.Ask([]*survey.Question{{
        Name: "color",
        Prompt: &survey.Select{
            Message: "Register or Login",
            Options: []string{"login", "register"},
            Default: "login",
        },
    }}, &authType)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Hello, %s!\n", authType)
}