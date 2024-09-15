package prompt

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

const (
	REGISTER    = "register"
	LOGIN       = "login"
	ADD_RECORD  = "Add a record"
	SEE_RECORDS = "See my records"
	EDIT_RECORD = "Edit record"
	EXIT        = "Exit"
)

const (
	CREDIT_CARD    = "Credit card data"
	LOGIN_PASSWORD = "Login and password"
	TEXT           = "Some text"
	BINARY_DATA    = "Binary data"
)

type CreditCardData struct {
	Number string
	Valid  string
	Name   string
	CVC    string
}

type LoginPassword struct {
	Login    string
	Password string
}

func EnterKey() (string, error) {

	var entryName string
	err := survey.AskOne(&survey.Input{Message: "Enter unique name to your item: "}, &entryName, survey.WithValidator(survey.Required))
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return entryName, nil
}

func EnterMetadata() (string, error) {
	var metadata string
	err := survey.AskOne(&survey.Input{Message: "Enter additional info: "}, &metadata)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return metadata, nil
}

func EnterLoginPassword() (*LoginPassword, error) {

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
	answers := LoginPassword{}

	err := survey.Ask(creds, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &answers, nil
}

func EnterText() (string, error) {
	var text string
	err := survey.AskOne(&survey.Input{Message: "Enter the text data: "}, &text)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return text, nil
}

func EnterFile() (string, error) {
	var file string
	err := survey.AskOne(&survey.Input{Message: "Enter path to file with binary data: "}, &file)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return file, nil
}

func EnterCreditCardData() (*CreditCardData, error) {
	creds := []*survey.Question{
		{
			Name:     "Number",
			Prompt:   &survey.Input{Message: "Card number: "},
			Validate: survey.Required,
		},
		{
			Name:     "Valid",
			Prompt:   &survey.Input{Message: "Valid through: "},
			Validate: survey.Required,
		},
		{
			Name:     "Name",
			Prompt:   &survey.Input{Message: "Owner name: "},
			Validate: survey.Required,
		},
		{
			Name:     "CVC",
			Prompt:   &survey.Password{Message: "CVC: "},
			Validate: survey.Required,
		},
	}
	answers := CreditCardData{}

	err := survey.Ask(creds, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &answers, err
}

func ChooseDataType() (string, error) {
	var dataType string

	err := survey.AskOne(&survey.Select{
		Message: "What kind of data would you like to store?",
		Options: []string{LOGIN_PASSWORD, CREDIT_CARD, TEXT, BINARY_DATA},
		Default: LOGIN_PASSWORD,
	}, &dataType)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return dataType, nil

}

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

func Menu() (string, error) {

	var action string

	err := survey.AskOne(&survey.Select{
		Message: "What do you want to do?",
		Options: []string{ADD_RECORD, SEE_RECORDS, EDIT_RECORD, EXIT},
		Default: ADD_RECORD,
	}, &action)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return action, nil
}

func Authenticate(method func(string, string) (string, error)) (string, error) {
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
	answers := LoginPassword{}

	err := survey.Ask(creds, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return method(answers.Login, answers.Password)
}
