package db

// TODO: act as coordinator?
type DB struct {
	schemas []Schema
	logger  any
	config  any
	wal     any
}

func NewDB() {

}

func (db *DB) CreateSchema() {

}

func (db *DB) DeleteSchema() {

}
