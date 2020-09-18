package resolve

import (
	"database/sql"
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
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
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (hash VARCHAR(32) PRIMARY KEY, pubkey TEXT, routing TEXT)", tableName)
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
		h  string
		p  string
		rt string
	)

	query := fmt.Sprintf("SELECT hash, pubkey, routing FROM %s WHERE hash LIKE ?", tableName)
	err := r.conn.QueryRow(query, addr.String).Scan(&h, &p, &rt)
	if err != nil {
		return nil, err
	}

	pk, err := bmcrypto.NewPubKey(p)
	if err != nil {
		return nil, err
	}

	return &Info{
		Hash:      h,
		PublicKey: *pk,
		Routing:   rt,
	}, nil
}

func (r *sqliteRepo) Upload(info *Info, _ bmcrypto.PrivKey, pow proofofwork.ProofOfWork) error {
	query := fmt.Sprintf("INSERT INTO %s(hash, pubkey , routing) VALUES (?, ?, ?)", tableName)
	st, err := r.conn.Prepare(query)
	if err != nil {
		return err
	}

	_, err = st.Exec(info.Hash, info.PublicKey.S, info.Routing)
	return err
}

func (r *sqliteRepo) Delete(info *Info, privKey bmcrypto.PrivKey) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE hash LIKE ?", tableName)
	st, err := r.conn.Prepare(query)
	if err != nil {
		return err
	}

	_, err = st.Exec(info.Hash)
	return err
}
