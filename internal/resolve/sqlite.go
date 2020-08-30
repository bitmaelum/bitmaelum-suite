package resolve

import (
	"database/sql"
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"strings"
	"sync"
)

const (
	tableName = "keyresolve"
)

type sqliteRepo struct {
	dsn  string
	conn *sql.DB
	mu   *sync.Mutex
}

// NewSqliteRepository creates new local repository where keys are stored in an SQLite database
func NewSqliteRepository(dsn string) (Repository, error) {
	// Work around some bugs/issues
	if !strings.HasPrefix(dsn, "file:") {
		if dsn == ":memory:" {
			dsn = "file::memory:?mode=memory&cache=shared"
		} else {
			dsn = fmt.Sprintf("file:%s?cache=shared&mode=rwc", dsn)
		}
	}

	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	db := &sqliteRepo{
		dsn:  dsn,
		conn: conn,
		mu:   new(sync.Mutex),
	}

	createTableIfNotExist(db)

	return db, nil
}

// createTableIfNotExist creates the key table if it doesn't exist already in the database
func createTableIfNotExist(db *sqliteRepo) {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (hash VARCHAR(32) PRIMARY KEY, pubkey TEXT, address TEXT)", tableName)
	st, err := db.conn.Prepare(query)
	if err != nil {
		return
	}

	_, err = st.Exec()
	if err != nil {
		return
	}
}

func (r *sqliteRepo) Resolve(addr address.HashAddress) (*Info, error) {
	var (
		h string
		p string
		a string
	)

	query := fmt.Sprintf("SELECT hash, pubkey, address FROM %s WHERE hash LIKE ?", tableName)
	err := r.conn.QueryRow(query, addr.String).Scan(&h, &p, &a)
	if err != nil {
		return nil, err
	}

	pk, err := encrypt.NewPubKey(p)
	if err != nil {
		return nil, err
	}

	return &Info{
		Hash:      h,
		PublicKey: *pk,
		Server:    a,
	}, nil
}

func (r *sqliteRepo) Upload(addr address.HashAddress, pubKey encrypt.PubKey, address, _ string) error {
	query := fmt.Sprintf("INSERT INTO %s(hash, pubkey , address) VALUES (?, ?, ?)", tableName)
	st, err := r.conn.Prepare(query)
	if err != nil {
		return err
	}

	_, err = st.Exec(addr.String(), pubKey.S, address)
	return err
}
