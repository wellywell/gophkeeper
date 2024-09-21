package main

import (
	"encoding/json"
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
			key, err := prompt.EnterKey("")
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			err = seeRecord(token, pass, key, cli)
			if err != nil {
				fmt.Println(err.Error())
			}
		case prompt.EDIT_RECORD:
			err = editRecord(token, pass, cli)
			if err != nil {
				fmt.Println(err.Error())
			}
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

	item, err := prompt.CreateBasicItem()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	switch action {
	case prompt.CREDIT_CARD:
		card, err := prompt.EnterCreditCardData()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(item, card)
	case prompt.LOGIN_PASSWORD:
		item.Type = types.TypeLogoPass
		logopass, err := prompt.EnterLoginPassword(types.LoginPassword{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = cli.CreateLoginPasswordItem(token, types.LoginPasswordItem{Item: *item, Data: logopass}, pass)
		if err != nil {
			fmt.Println(err.Error())
			return
		} else {
			fmt.Println("saved")
		}
	}

}

func seeRecord(token string, pass string, key string, cli *client.Client) error {
	data, err := cli.GetItem(token, key)
	if err != nil {
		return err
	}

	var i types.AnyItem
	err = json.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	fmt.Println(i.Item.String())
	switch i.Item.Type {
	case types.TypeLogoPass:
		logopassItem, err := types.ParseItem[*types.LoginPassword](data, pass)
		if err != nil {
			return err
		}
		fmt.Println(logopassItem.Data.String())
	}
	return nil
}

func seeRecords(token string, pass string, cli *client.Client) {
	fmt.Println("Unimplimented")
}

func editRecord(token string, pass string, cli *client.Client) error {
	key, err := prompt.EnterKey("")
	if err != nil {
		return err
	}
	data, err := cli.GetItem(token, key)
	if err != nil {
		return err
	}

	var i types.AnyItem
	err = json.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	fmt.Println(i.Item.String())

	result, err := prompt.ChooseEditOrDelete()
	if err != nil {
		return err
	}

	switch result {
	case prompt.DELETE:
		err = cli.DeleteItem(token, key)
		if err != nil {
			return err
		}
		fmt.Println("Deleted")
		return nil
	case prompt.EDIT:
		switch i.Item.Type {
		case types.TypeLogoPass:
			logopassItem, err := types.ParseItem[*types.LoginPassword](data, pass)
			if err != nil {
				return err
			}
			err = updateLogoPassData(token, pass, logopassItem, cli)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("success")
			}
		}
	}
	return nil
}

func updateLogoPassData(token string, pass string, logopass *types.GenericItem[*types.LoginPassword], cli *client.Client) error {

	meta, err := prompt.EnterMetadata(logopass.Item.Info)
	if err != nil {
		return err
	}
	newLogoPass, err := prompt.EnterLoginPassword(*logopass.Data)
	if err != nil {
		return err
	}
	if meta == logopass.Item.Info && *newLogoPass == *logopass.Data {
		return fmt.Errorf("nothing changed")
	}
	newItem := &types.LoginPasswordItem{Item: types.Item{Key: logopass.Item.Key, Info: meta}, Data: newLogoPass}

	return cli.UpdateLogoPassData(token, pass, *newItem)
}
