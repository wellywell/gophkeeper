package main

import (
	"fmt"

	"github.com/wellywell/gophkeeper/internal/client"
	"github.com/wellywell/gophkeeper/internal/client/prompt"
	"github.com/wellywell/gophkeeper/internal/config"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s", buildVersion, buildDate, buildCommit)

	conf, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	cli, err := client.NewClient(conf.RunAddress)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	token, pass, err := authenticate(cli)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for {
		record, err := prompt.Menu()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		switch record {
		case prompt.EXIT:
			fmt.Println("Bye!")
			return
		case prompt.ADD_RECORD:
			addRecord(token, pass, cli)
		case prompt.SEE_RECORDS:
			seeRecords(token, pass, cli)
		case prompt.EDIT_RECORD:
			editRecord(token, pass, cli)
		}
	}
}

func authenticate(cli *client.Client) (string, string, error) {
	authMethod, err := prompt.ChooseLoginOrRegister()
	if err != nil {
		fmt.Println(err.Error())
		return "", "", err
	}

	var method func(string, string) (string, error)

	switch authMethod {
	case prompt.LOGIN:
		method = cli.Login
	case prompt.REGISTER:
		method = cli.Register
	default:
		fmt.Println("Error authenticating")
		return "", "", err
	}

	return prompt.Authenticate(method)
}

func addRecord(token string, pass string, cli *client.Client) {
	action, err := prompt.ChooseDataType()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(action)
	fmt.Println("Unimplimented")

}

func seeRecords(token string, pass string, cli *client.Client) {
	fmt.Println("Unimplimented")
}

func editRecord(token string, pass string, cli *client.Client) {
	//key, err := prompt.EnterKey()
	fmt.Println("Unimplimented")
}
