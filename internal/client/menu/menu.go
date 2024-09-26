package menu

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wellywell/gophkeeper/internal/client"
	"github.com/wellywell/gophkeeper/internal/client/prompt"
	"github.com/wellywell/gophkeeper/internal/types"
)

func MainMenu(token string, pass string, cli *client.Client) {

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

func Authenticate(cli *client.Client) (string, string, error) {
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
		// empty object is passed for new item
		card, err := prompt.EnterCreditCardData(types.CreditCardData{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = client.CreateItem[*types.CreditCardData](token, pass, types.GenericItem[*types.CreditCardData]{Item: *item, Data: card}, cli.CreateCreditCardItem)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	case prompt.LOGIN_PASSWORD:
		item.Type = types.TypeLogoPass
		// empty object is passed for new item
		logopass, err := prompt.EnterLoginPassword(types.LoginPassword{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = client.CreateItem(token, pass, types.GenericItem[*types.LoginPassword]{Item: *item, Data: logopass}, cli.CreateLoginPasswordItem)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	case prompt.TEXT:
		item.Type = types.TypeText
		text, err := prompt.EnterText("")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = client.CreateItem(token, pass, types.GenericItem[*types.TextData]{Item: *item, Data: &text}, cli.CreateTextItem)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	case prompt.BINARY_DATA:
		item.Type = types.TypeBinary
		filename, err := prompt.EnterFileName()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		dat, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := types.BinaryData(dat)
		err = client.CreateItem(token, pass, types.GenericItem[*types.BinaryData]{Item: *item, Data: &data}, cli.CreateBinaryItem)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	fmt.Println("saved")
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
	case types.TypeCreditCard:
		card, err := types.ParseItem[*types.CreditCardData](data, pass)
		if err != nil {
			return err
		}
		fmt.Println(card.Data.String())
	case types.TypeText:
		text, err := types.ParseItem[*types.TextData](data, pass)
		if err != nil {
			return err
		}
		fmt.Println(text.Data.String())
	case types.TypeBinary:
		fmt.Println("to download binary content use download menu")
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
			}
		case types.TypeCreditCard:
			card, err := types.ParseItem[*types.CreditCardData](data, pass)
			if err != nil {
				return err
			}
			err = updateCreditCardData(token, pass, card, cli)
			if err != nil {
				fmt.Println(err.Error())
			}
		case types.TypeText:
			text, err := types.ParseItem[*types.TextData](data, pass)
			if err != nil {
				return err
			}
			err = updateTextData(token, pass, text, cli)
			if err != nil {
				fmt.Println(err.Error())
			}
		case types.TypeBinary:
			data, err := types.ParseItem[*types.BinaryData](data, pass)
			if err != nil {
				return err
			}
			err = updateBinaryData(token, pass, data, cli)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		fmt.Println("success")
	}
	return nil
}

func updateBinaryData(token string, pass string, data *types.GenericItem[*types.BinaryData], cli *client.Client) error {
	meta, err := prompt.EnterMetadata(data.Item.Info)
	if err != nil {
		return err
	}
	filename, err := prompt.EnterFileName()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	dat, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	d := types.BinaryData(dat)

	newItem := types.GenericItem[*types.BinaryData]{Item: types.Item{Key: data.Item.Key, Info: meta, Type: data.Item.Type}, Data: &d}

	return client.UpdateItem(token, pass, newItem, cli.UpdateBinaryItem)
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
	newItem := types.GenericItem[*types.LoginPassword]{Item: types.Item{Key: logopass.Item.Key, Info: meta}, Data: newLogoPass}

	return client.UpdateItem(token, pass, newItem, cli.UpdateLogoPassData)
}

func updateCreditCardData(token string, pass string, card *types.GenericItem[*types.CreditCardData], cli *client.Client) error {

	meta, err := prompt.EnterMetadata(card.Item.Info)
	if err != nil {
		return err
	}
	newData, err := prompt.EnterCreditCardData(*card.Data)
	if err != nil {
		return err
	}
	if meta == card.Item.Info && *newData == *card.Data {
		return fmt.Errorf("nothing changed")
	}
	newItem := types.GenericItem[*types.CreditCardData]{Item: types.Item{Key: card.Item.Key, Info: meta}, Data: newData}

	return client.UpdateItem(token, pass, newItem, cli.UpdateCreditCardData)
}

func updateTextData(token string, pass string, text *types.GenericItem[*types.TextData], cli *client.Client) error {

	meta, err := prompt.EnterMetadata(text.Item.Info)
	if err != nil {
		return err
	}
	newData, err := prompt.EnterText(string(*text.Data))
	if err != nil {
		return err
	}
	if meta == text.Item.Info && newData == *text.Data {
		return fmt.Errorf("nothing changed")
	}
	newItem := types.GenericItem[*types.TextData]{Item: types.Item{Key: text.Item.Key, Info: meta}, Data: &newData}

	return client.UpdateItem(token, pass, newItem, cli.UpdateTextData)
}
