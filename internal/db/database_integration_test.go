//go:build integration_tests
// +build integration_tests

/* В связи с санкциями, нужен VPN, чтобы докерхаб работал */

package db

import (
	"log"
	"os"
	"testing"

	"github.com/wellywell/bonusy/internal/testutils"
)

var DBDSN string

func TestMain(m *testing.M) {
	code, err := runMain(m)

	if err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

func runMain(m *testing.M) (int, error) {

	databaseDSN, cleanUp, err := testutils.RunTestDatabase()
	defer cleanUp()

	if err != nil {
		return 1, err
	}
	DBDSN = databaseDSN

	exitCode := m.Run()

	return exitCode, nil

}
