// Package main - запуск клиента для хранения и получения чувствительных данных
package main

import (
	"fmt"

	"github.com/wellywell/gophkeeper/internal/client"
	"github.com/wellywell/gophkeeper/internal/client/menu"
	"github.com/wellywell/gophkeeper/internal/config"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s", buildVersion, buildDate, buildCommit)

	conf, err := config.NewClientConfig()
	if err != nil {
		panic(err)
	}

	cli, err := client.NewClient(conf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	token, pass, err := menu.Authenticate(cli)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	menu.MainMenu(token, pass, cli)
}
