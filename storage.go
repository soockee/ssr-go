package main

import (
	"zombiezen.com/go/sqlite"
)

type Storage interface {
}

type SQLiteStore struct {
	db *sqlite.Conn
}

func NewSQLiteStore() (*SQLiteStore, error) {
	conn, err := sqlite.OpenConn(":memory:", sqlite.OpenReadWrite)
	if err != nil {
		return nil, err
	}

	return &SQLiteStore{
		db: conn,
	}, nil
}
