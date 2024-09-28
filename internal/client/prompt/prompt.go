// Package prompt пакет содержит методы для интерактивных действий с пользователем из консольного клиента
package prompt

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/wellywell/gophkeeper/internal/types"
)

const (
	REGISTER = "register"
	LOGIN    = "login"
)

const (
	ADD_RECORD  = "Add a record"
	SEE_RECORD  = "Show record data"
	SEE_RECORDS = "List all records"
	EDIT_RECORD = "Edit record"
	DOWNLOAD    = "Download binary data"
	EXIT        = "Exit"
	CANCEL      = "Back to main menu"
	NEXT        = "Next page"
)

const (
	CREDIT_CARD    = "Credit card data"
	LOGIN_PASSWORD = "Login and password"
	TEXT           = "Some text"
	BINARY_DATA    = "Binary data"
)

const (
	EDIT   = "edit"
	DELETE = "delete"
)

// EnterKey промпт для ввода названия записи для хранения на сервере
func EnterKey(key string) (string, error) {

	var entryName string
	err := survey.AskOne(&survey.Input{Message: "Enter unique name of your item: ", Default: key}, &entryName, survey.WithValidator(survey.Required))
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return entryName, nil
}

// NextBackExit промпт предлагает загрузить следующую страницу с данными, либо вернуться в главное меню или закрыть программу
func NextBackExit() (string, error) {
	var action string

	err := survey.AskOne(&survey.Select{
		Message: "",
		Options: []string{NEXT, CANCEL, EXIT},
		Default: NEXT,
	}, &action)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return action, nil
}

// EnterMetadata предлагает ввести допольнительную информацию об объекте
func EnterMetadata(info string) (string, error) {
	var metadata string
	err := survey.AskOne(&survey.Input{Message: "Enter additional info: ", Default: info}, &metadata)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return metadata, nil
}

// EnterLoginPassword предлагает ввести логин и пароль для сохранения на сервере
func EnterLoginPassword(item types.LoginPassword) (*types.LoginPassword, error) {

	creds := []*survey.Question{
		{
			Name:     "login",
			Prompt:   &survey.Input{Message: "Login: ", Default: item.Login},
			Validate: survey.Required,
		},
		{
			Name:     "password",
			Prompt:   &survey.Password{Message: "Password: "},
			Validate: survey.Required,
		}}
	answers := types.LoginPassword{}

	err := survey.Ask(creds, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &answers, nil
}

// EnterText предлагает ввести текст для сохранения на сервере
func EnterText(defaultText string) (types.TextData, error) {
	var text string
	err := survey.AskOne(&survey.Input{Message: "Enter the text data: ", Default: defaultText}, &text)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return types.TextData(text), nil
}

// EnterFileName предлагает ввести имя (существующего) файла для последующего открытия этого файла
func EnterFileName() (string, error) {
	file := ""
	prompt := &survey.Input{
		Message: "File to read binary data from...",
		Suggest: func(toComplete string) []string {
			files, _ := filepath.Glob(toComplete + "*")
			return files
		},
	}
	err := survey.AskOne(prompt, &file)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return file, nil
}

// EnterFile предлагает ввести имя нового файла для сохранения данных в нём
func EnterFile() (string, error) {
	var file string
	err := survey.AskOne(&survey.Input{Message: "Enter filename to write data to: "}, &file)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return file, nil
}

// EnterCreditCardData предлагает ввести данные кредитной карты
func EnterCreditCardData(card types.CreditCardData) (*types.CreditCardData, error) {

	creds := []*survey.Question{
		{
			Name:   "Number",
			Prompt: &survey.Input{Message: "Card number: ", Default: card.Number},
			Validate: func(val interface{}) error {
				_, err := strconv.Atoi(val.(string))
				if err != nil {
					return errors.New("digits only")
				}
				err = survey.MaxLength(16)(val)
				if err != nil {
					return err
				}
				return survey.MinLength(15)(val)
			},
		},
		{
			Name:   "ValidYear",
			Prompt: &survey.Input{Message: "Valid through (year): ", Default: card.ValidYear},
			Validate: func(val interface{}) error {
				num, err := strconv.Atoi(val.(string))
				if err != nil {
					return errors.New("invalid year (numbers only)")
				}
				if num < 1000 || num > 9999 {
					return errors.New("4 digits only")
				}
				return nil
			},
		},
		{
			Name:   "ValidMonth",
			Prompt: &survey.Input{Message: "Valid through (month): ", Default: card.ValidMonth},
			Validate: func(val interface{}) error {
				num, err := strconv.Atoi(val.(string))
				if err != nil {
					return errors.New("invalid month (numbers only)")
				}
				if num < 1 || num > 12 {
					return errors.New("month is from 1 to 12")
				}
				return nil
			},
		},
		{
			Name:     "Name",
			Prompt:   &survey.Input{Message: "Owner name: ", Default: card.Name},
			Validate: survey.Required,
		},
		{
			Name:   "CVC",
			Prompt: &survey.Password{Message: "CVC: "},
			Validate: func(val interface{}) error {
				num, err := strconv.Atoi(val.(string))
				if err != nil {
					return errors.New("invalid year (numbers only)")
				}
				if num < 100 || num > 999 {
					return errors.New("3 digits only")
				}
				return nil
			},
		},
	}
	answers := types.CreditCardData{}

	err := survey.Ask(creds, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &answers, err
}

// ChooseDataType предлагает выбрать, какой тип данных хочет сохранить пользователь
func ChooseDataType() (string, error) {
	var dataType string

	err := survey.AskOne(&survey.Select{
		Message: "What kind of data would you like to store?",
		Options: []string{LOGIN_PASSWORD, CREDIT_CARD, TEXT, BINARY_DATA, CANCEL},
		Default: LOGIN_PASSWORD,
	}, &dataType)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return dataType, nil

}

// ChooseLoginOrRegister предлагает залогиниться или зарегистрироваться как новый пользователь
func ChooseLoginOrRegister() (string, error) {

	var authType string

	err := survey.AskOne(&survey.Select{
		Message: "Register or Login",
		Options: []string{LOGIN, REGISTER},
		Default: LOGIN,
	}, &authType)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return authType, nil
}

// ChooseEditOrDelete предлагает выбрать, хочет пользователь удалить данные, либо отредактировать
func ChooseEditOrDelete() (string, error) {

	var action string

	err := survey.AskOne(&survey.Select{
		Message: "Would you like to edit or delete item?",
		Options: []string{EDIT, DELETE, CANCEL},
		Default: EDIT,
	}, &action)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return action, nil
}

// Menu промпт корневого меню - предлагает набор действий пользователю - просмотреть записи,
// отредактировать запись, получить запись по ключу, загрузить бинарные данные с сервера в файл
func Menu() (string, error) {

	var action string

	err := survey.AskOne(&survey.Select{
		Message: "What do you want to do?",
		Options: []string{ADD_RECORD, SEE_RECORDS, SEE_RECORD, EDIT_RECORD, DOWNLOAD, EXIT},
		Default: ADD_RECORD,
	}, &action)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return action, nil
}

// Authenticate аутентификация пользователя
func Authenticate(method func(string, string) (string, error)) (string, string, error) {
	creds := []*survey.Question{
		{
			Name:     "login",
			Prompt:   &survey.Input{Message: "Login: "},
			Validate: survey.Required,
		},
		{
			Name:     "password",
			Prompt:   &survey.Password{Message: "Password: "},
			Validate: survey.Required,
		}}
	answers := types.LoginPassword{}

	err := survey.Ask(creds, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return "", "", err
	}
	token, err := method(answers.Login, answers.Password)
	return token, answers.Password, err

}

// CreateBasicItem создаёт метаданные для любого типа данных
func CreateBasicItem() (*types.Item, error) {
	key, err := EnterKey("")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	meta, err := EnterMetadata("")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &types.Item{
		Key:  key,
		Info: meta,
	}, nil
}
