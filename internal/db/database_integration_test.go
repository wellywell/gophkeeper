//go:build integration_tests
// +build integration_tests

/* В связи с санкциями, нужен VPN, чтобы докерхаб работал */

package db

import (
	"log"
	"os"
	"testing"

	"gotest.tools/assert"

	"github.com/wellywell/gophkeeper/internal/testutils"
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

func TestNewDatabase(t *testing.T) {

	d, err := NewDatabase(DBDSN)
	assert.NilError(t, err)

	err = d.Close()
	assert.NilError(t, err)
}
