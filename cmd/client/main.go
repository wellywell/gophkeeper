package main

import (
	"fmt"

	"github.com/wellywell/gophkeeper/internal/client"
	"github.com/wellywell/gophkeeper/internal/client/prompt"
	"github.com/wellywell/gophkeeper/internal/config"
	"github.com/wellywell/gophkeeper/internal/types"
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
		case prompt.SEE_RECORD:
			err := seeRecord(token, pass, cli)
			if err != nil {
				fmt.Println(err.Error())
			}
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
	switch action {
	case prompt.CREDIT_CARD:
		item, err := CreateBasicItem()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		card, err := prompt.EnterCreditCardData()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(item, card)
	case prompt.LOGIN_PASSWORD:
		item, err := CreateBasicItem()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		item.Type = types.TypeLogoPass
		logopass, err := prompt.EnterLoginPassword()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = cli.CreateLoginPasswordItem(token, types.LoginPasswordItem{Item: *item, Data: *logopass}, pass)
		if err != nil {
			fmt.Println(err.Error())
			return
		} else {
			fmt.Println("saved")
		}
	}

}

func CreateBasicItem() (*types.Item, error) {
	key, err := prompt.EnterKey()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	meta, err := prompt.EnterMetadata()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &types.Item{
		Key:  key,
		Info: meta,
	}, nil
}

func seeRecord(token string, pass string, cli *client.Client) error {
	key, err := prompt.EnterKey()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	result, err := cli.SeeRecord(token, pass, key)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}

func seeRecords(token string, pass string, cli *client.Client) {
	fmt.Println("Unimplimented")
}

func editRecord(token string, pass string, cli *client.Client) {
	//key, err := prompt.EnterKey()
	fmt.Println("Unimplimented")
}
