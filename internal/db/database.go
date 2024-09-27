package db

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/wellywell/gophkeeper/internal/types"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func NewDatabase(connString string) (*Database, error) {

	err := Migrate(connString)

	if err != nil {
		return nil, fmt.Errorf("failed to migrate %w", err)
	}

	ctx := context.Background()
	p, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	return &Database{
		pool: p,
	}, nil
}

func (d *Database) CreateUser(ctx context.Context, username string, password string) error {

	query := `
		INSERT INTO auth_user (username, password)
		VALUES ($1, $2)
		`
	_, err := d.pool.Exec(ctx, query, username, password)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return fmt.Errorf("%w", &UserExistsError{Username: username})
		}
		return err
	}
	return nil
}

func (d *Database) GetUserHashedPassword(ctx context.Context, username string) (string, error) {
	query := `
		SELECT password 
		FROM auth_user 
		WHERE username = $1`

	row := d.pool.QueryRow(ctx, query, username)

	var password string

	err := row.Scan(&password)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("%w", &UserNotFoundError{Username: username})
	}
	return password, nil
}

func (d *Database) GetUserID(ctx context.Context, username string) (int, error) {
	query := `
		SELECT id 
		FROM auth_user 
		WHERE username = $1`

	row := d.pool.QueryRow(ctx, query, username)

	var id int

	err := row.Scan(&id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("%w", &UserNotFoundError{Username: username})
	}
	return id, nil

}

func (d *Database) InsertItem(ctx context.Context, tx pgx.Tx, userID int, item types.Item) (int, error) {
	query := `
		INSERT INTO item(user_id, key, item_type, info)
		VALUES ($1, $2, $3, $4)
		RETURNING id `

	row := tx.QueryRow(ctx, query, userID, item.Key, item.Type, item.Info)

	var itemID int
	if err := row.Scan(&itemID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return 0, fmt.Errorf("%w", &KeyExistsError{Key: item.Key})
		}
		return 0, fmt.Errorf("%w", err)
	}
	return itemID, nil
}

func (d *Database) UpdateItem(ctx context.Context, tx pgx.Tx, userID int, item types.Item) (int, error) {
	query := `
		UPDATE item
		SET info = $1
		WHERE key = $2 AND user_id = $3
		RETURNING id `

	row := tx.QueryRow(ctx, query, item.Info, item.Key, userID)

	var itemID int
	if err := row.Scan(&itemID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("%w", &KeyNotFoundError{Key: item.Key})
		}
		return 0, fmt.Errorf("unexpected db error %w", err)
	}
	return itemID, nil
}

func (d *Database) InsertText(ctx context.Context, userID int, item types.TextItem) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	itemID, err := d.InsertItem(ctx, tx, userID, types.Item{Key: item.Item.Key, Type: types.TypeText, Info: item.Item.Info})
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	query := `
		INSERT INTO text_data (item_id, data)
		VALUES ($1, $2)
	`
	_, err = tx.Exec(ctx, query, itemID, item.Data)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (d *Database) InsertCreditCard(ctx context.Context, userID int, item types.CreditCardItem) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	itemID, err := d.InsertItem(ctx, tx, userID, types.Item{Key: item.Item.Key, Type: types.TypeCreditCard, Info: item.Item.Info})
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	query := `
		INSERT INTO credit_card (item_id, number, owner_name, valid_till, cvc)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, query, itemID, item.Data.Number, item.Data.Name, fmt.Sprintf("%s-%s-01", item.Data.ValidYear, item.Data.ValidMonth), item.Data.CVC)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (d *Database) InsertBinaryData(ctx context.Context, userID int, item types.BinaryItem) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	itemID, err := d.InsertItem(ctx, tx, userID, types.Item{Key: item.Item.Key, Type: types.TypeBinary, Info: item.Item.Info})
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	query := `
		INSERT INTO binary_data (item_id, data)
		VALUES ($1, $2)
	`
	_, err = tx.Exec(ctx, query, itemID, item.Data)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (d *Database) InsertLogoPass(ctx context.Context, userID int, data types.LoginPasswordItem) error {

	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	itemID, err := d.InsertItem(ctx, tx, userID, data.Item)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	query := `
		INSERT INTO logopass (item_id, login, password)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(ctx, query, itemID, data.Data.Login, data.Data.Password)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (d *Database) UpdateLogoPass(ctx context.Context, userID int, data types.LoginPasswordItem) error {

	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	itemID, err := d.UpdateItem(ctx, tx, userID, data.Item)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	query := `
		UPDATE logopass
		SET login = $1, password = $2
		WHERE item_id = $3
	`
	_, err = tx.Exec(ctx, query, data.Data.Login, data.Data.Password, itemID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (d *Database) UpdateCreditCard(ctx context.Context, userID int, data types.CreditCardItem) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	itemID, err := d.UpdateItem(ctx, tx, userID, data.Item)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	query := `
		UPDATE credit_card
		SET number = $1, owner_name = $2, valid_till = $3, cvc = $4
		WHERE item_id = $5
	`
	_, err = tx.Exec(ctx, query, data.Data.Number, data.Data.Name, fmt.Sprintf("%s-%s-01", data.Data.ValidYear, data.Data.ValidMonth), data.Data.CVC, itemID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (d *Database) UpdateBinaryData(ctx context.Context, userID int, data types.BinaryItem) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	itemID, err := d.UpdateItem(ctx, tx, userID, data.Item)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	query := `
		UPDATE binary_data
		SET data = $1
		WHERE item_id = $2
	`
	_, err = tx.Exec(ctx, query, data.Data, itemID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (d *Database) UpdateText(ctx context.Context, userID int, data types.TextItem) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	itemID, err := d.UpdateItem(ctx, tx, userID, data.Item)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	query := `
		UPDATE text_data
		SET data = $1
		WHERE item_id = $2
	`
	_, err = tx.Exec(ctx, query, data.Data, itemID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (d *Database) DeleteItem(ctx context.Context, userID int, key string) error {
	query := `
		DELETE FROM item
		WHERE user_id = $1 and key = $2
	`
	_, err := d.pool.Exec(ctx, query, userID, key)

	if err != nil {
		return err
	}
	return nil
}

func (d *Database) GetBinaryData(ctx context.Context, userID int, key string) ([]byte, error) {

	item, err := d.GetItem(ctx, userID, key)
	if err != nil {
		return nil, err
	}

	if item.Type != types.TypeBinary {
		return nil, &KeyNotFoundError{Key: key}
	}

	query := `
		SELECT data
		FROM binary_data
		WHERE item_id = $1
	`

	rows, err := d.pool.Query(ctx, query, item.Id)
	if err != nil {
		return nil, fmt.Errorf("failed collecting rows %w", err)
	}

	data, err := pgx.CollectOneRow(rows, pgx.RowTo[[]byte])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &KeyNotFoundError{Key: key}
		}
		return nil, fmt.Errorf("failed unpacking rows %w", err)
	}
	return data, nil
}

func (d *Database) GetItem(ctx context.Context, userID int, key string) (*types.Item, error) {
	query := `
		SELECT id, item_type, info, key
		FROM item
		WHERE user_id = $1 AND key = $2
	`

	rows, err := d.pool.Query(ctx, query, userID, key)
	if err != nil {
		return nil, fmt.Errorf("failed collecting rows %w", err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[types.Item])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &KeyNotFoundError{Key: key}
		}
		return nil, fmt.Errorf("failed unpacking rows %w", err)
	}
	return &item, nil
}

func (d *Database) GetLogoPass(ctx context.Context, itemID int) (*types.LoginPassword, error) {
	query := `
		SELECT login, password
		FROM logopass
		WHERE item_id = $1
	`

	rows, err := d.pool.Query(ctx, query, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed collecting rows %w", err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[types.LoginPassword])
	if err != nil {
		return nil, fmt.Errorf("failed unpacking rows %w", err)
	}
	return &item, nil
}

func (d *Database) GetCreditCard(ctx context.Context, itemID int) (*types.CreditCardData, error) {
	query := `
		SELECT number, owner_name, valid_till, cvc
		FROM credit_card
		WHERE item_id = $1
	`

	rows, err := d.pool.Query(ctx, query, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed collecting rows %w", err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[types.CreditCardData])
	if err != nil {
		return nil, fmt.Errorf("failed unpacking rows %w", err)
	}
	item.ValidMonth = strconv.Itoa(int(item.ValidDate.Month()))
	item.ValidYear = strconv.Itoa(item.ValidDate.Year())
	return &item, nil
}

func (d *Database) GetText(ctx context.Context, itemID int) (*types.TextData, error) {
	query := `
		SELECT data
		FROM text_data
		WHERE item_id = $1
	`

	row := d.pool.QueryRow(ctx, query, itemID)

	var text types.TextData

	err := row.Scan(&text)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &text, nil
}

func (d *Database) GetItems(ctx context.Context, userID int, limit int, offset int) ([]types.Item, error) {

	query := `
		SELECT id, item_type, info, key
		FROM item
		WHERE user_id = $1
		ORDER BY id
		LIMIT $2
		OFFSET $3
	`

	rows, err := d.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed collecting rows %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.Item])
	if err != nil {
		return nil, fmt.Errorf("failed unpacking rows %w", err)
	}
	return items, nil
}

// Close завершает работу базы данных
func (d *Database) Close() error {
	d.pool.Close()
	return nil
}
