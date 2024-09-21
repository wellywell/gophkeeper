package db

import (
	"context"
	"errors"
	"fmt"

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

func (d *Database) InsertLogoPass(ctx context.Context, userID int, data types.LoginPasswordItem) error {

	query := `
		INSERT INTO item(user_id, key, item_type, info)
		VALUES ($1, $2, $3, $4)
		RETURNING id `

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

	row := tx.QueryRow(ctx, query, userID, data.Item.Key, data.Item.Type, data.Item.Info)

	var itemID int
	if err := row.Scan(&itemID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return fmt.Errorf("%w", &KeyExistsError{Key: data.Item.Key})
		}
		return fmt.Errorf("%w", err)
	}
	query = `
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

func (d *Database) InsertCreditCard(ctx context.Context, userID int, key string, card types.CreditCardData, meta string) error {
	return nil
}

func (d *Database) InsertTextData(ctx context.Context, userID int, key string, text string, meta string) error {
	return nil
}

func (d *Database) InsertBinaryData(ctx context.Context, userID int, key string, data []byte, meta string) error {
	return nil
}

func (d *Database) DeleteItem(ctx context.Context, userID int, key string) error {
	return nil
}

func (d *Database) UpdateLogoPass(ctx context.Context, userID int, key string, logopass types.LoginPassword, meta string) error {
	return nil
}

func (d *Database) UpdateCreditCard(ctx context.Context, userID int, key string, card types.CreditCardData, meta string) error {
	return nil
}

func (d *Database) UpdateTextData(ctx context.Context, userID int, key string, text string, meta string) error {
	return nil
}

func (d *Database) UpdateBinaryData(ctx context.Context, userID int, key string, data []byte, meta string) error {
	return nil
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

func (d *Database) GetAllItems(ctx context.Context, userID int) ([]types.Item, error) {
	return nil, nil
}

// Close завершает работу базы данных
func (d *Database) Close() error {
	d.pool.Close()
	return nil
}
