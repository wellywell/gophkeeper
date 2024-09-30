//go:build integration_tests
// +build integration_tests

/* В связи с санкциями, нужен VPN, чтобы докерхаб работал */

package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wellywell/gophkeeper/internal/testutils"
	"github.com/wellywell/gophkeeper/internal/types"
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
	assert.NoError(t, err)

	err = d.Close()
	assert.NoError(t, err)
}

func TestUserMethods(t *testing.T) {

	d, err := NewDatabase(DBDSN)
	assert.NoError(t, err)

	err = d.CreateUser(context.Background(), "myUser", "pass")
	assert.NoError(t, err)

	err = d.CreateUser(context.Background(), "myUser", "pass")
	assert.Error(t, err)

	pass, err := d.GetUserHashedPassword(context.Background(), "myUser")
	assert.NoError(t, err)
	assert.Equal(t, "pass", pass)

	id, err := d.GetUserID(context.Background(), "myUser")
	assert.NoError(t, err)
	assert.Greater(t, id, 0)
}

func TestItemMethods(t *testing.T) {

	d, _ := NewDatabase(DBDSN)

	_ = d.CreateUser(context.Background(), "myUser", "pass")

	userID, err := d.GetUserID(context.Background(), "myUser")
	assert.NoError(t, err)

	tx, err := d.pool.Begin(context.Background())
	defer func() {
		err = tx.Rollback(context.Background())
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	_, err = d.InsertItem(context.Background(), tx, userID, types.Item{Key: "1", Type: "text"})
	assert.NoError(t, err)

	_ = tx.Commit(context.Background())

	tx, err = d.pool.Begin(context.Background())

	_, err = d.UpdateItem(context.Background(), tx, userID, types.Item{Key: "1", Type: "text", Info: "new info"})
	assert.NoError(t, err)

	_ = tx.Commit(context.Background())

	i, err := d.GetItem(context.Background(), userID, "1")
	assert.NoError(t, err)

	assert.Equal(t, "new info", i.Info)

	defer func() {
		err = tx.Rollback(context.Background())
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	_, err = d.InsertItem(context.Background(), tx, userID, types.Item{Key: "1", Type: "text"})
	assert.Error(t, err)
}

func TestStoreDataMethods(t *testing.T) {

	d, _ := NewDatabase(DBDSN)

	_ = d.CreateUser(context.Background(), "myUser", "pass")

	userID, err := d.GetUserID(context.Background(), "myUser")
	assert.NoError(t, err)

	ctx := context.Background()

	err = d.InsertCreditCard(ctx, userID, types.CreditCardItem{Item: types.Item{Type: types.TypeCreditCard, Key: "2"}, Data: &types.CreditCardData{ValidMonth: "1", ValidYear: "2000"}})
	assert.NoError(t, err)

	err = d.InsertCreditCard(ctx, userID, types.CreditCardItem{Item: types.Item{Type: types.TypeCreditCard, Key: "2"}, Data: &types.CreditCardData{ValidMonth: "1", ValidYear: "2000"}})
	assert.Error(t, err)

	err = d.UpdateCreditCard(ctx, userID, types.CreditCardItem{Item: types.Item{Type: types.TypeCreditCard, Key: "2"}, Data: &types.CreditCardData{ValidMonth: "1", ValidYear: "2005"}})
	assert.NoError(t, err)

	i, err := d.GetItem(ctx, userID, "2")
	assert.NoError(t, err)
	card, err := d.GetCreditCard(ctx, i.Id)
	assert.NoError(t, err)

	assert.Equal(t, "2005", card.ValidYear)

	err = d.InsertLogoPass(ctx, userID, types.LoginPasswordItem{Item: types.Item{Type: types.TypeLogoPass, Key: "3"}, Data: &types.LoginPassword{}})
	assert.NoError(t, err)

	err = d.InsertLogoPass(ctx, userID, types.LoginPasswordItem{Item: types.Item{Type: types.TypeLogoPass, Key: "3"}, Data: &types.LoginPassword{}})
	assert.Error(t, err)

	err = d.UpdateLogoPass(ctx, userID, types.LoginPasswordItem{Item: types.Item{Type: types.TypeLogoPass, Key: "3"}, Data: &types.LoginPassword{Login: "a"}})
	assert.NoError(t, err)

	i, err = d.GetItem(ctx, userID, "3")
	assert.NoError(t, err)
	logopass, err := d.GetLogoPass(ctx, i.Id)
	assert.NoError(t, err)
	assert.Equal(t, "a", logopass.Login)

	err = d.InsertText(ctx, userID, types.TextItem{Item: types.Item{Type: types.TypeText, Key: "4"}})
	assert.NoError(t, err)

	err = d.InsertText(ctx, userID, types.TextItem{Item: types.Item{Type: types.TypeText, Key: "4"}})
	assert.Error(t, err)

	err = d.UpdateText(ctx, userID, types.TextItem{Item: types.Item{Type: types.TypeText, Key: "4"}, Data: "sss"})
	assert.NoError(t, err)

	i, err = d.GetItem(ctx, userID, "4")
	assert.NoError(t, err)
	text, err := d.GetText(ctx, i.Id)
	assert.NoError(t, err)

	assert.Equal(t, "sss", string(*text))

	err = d.InsertBinaryData(ctx, userID, types.BinaryItem{Item: types.Item{Type: types.TypeText, Key: "5"}})
	assert.NoError(t, err)

	err = d.InsertBinaryData(ctx, userID, types.BinaryItem{Item: types.Item{Type: types.TypeText, Key: "5"}})
	assert.Error(t, err)

	err = d.UpdateBinaryData(ctx, userID, types.BinaryItem{Item: types.Item{Type: types.TypeText, Key: "5"}, Data: []byte("www")})
	assert.NoError(t, err)

	data, err := d.GetBinaryData(ctx, userID, "5")
	assert.NoError(t, err)

	assert.Equal(t, "www", string(data))

}
