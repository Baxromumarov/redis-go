package database

import (
	"strings"

	err "github.com/baxromumarov/redis-go/errors"
)

type Database struct {
	Store map[string]string
	Expiration
}

func Init() *Database {
	return &Database{
		Store: make(map[string]string),
	}
}

func (db *Database) Set(key, val string) error {
	if key == "" {
		return err.ErrEmptyKey
	}
	if val == "" {
		return err.ErrEmpyVal
	}

	db.Store[key] = val

	return nil
}

func (db *Database) Get(key string) (string, error) {
	val, exist := db.Store[key]
	if !exist {
		return "", err.ErrContentNotFound
	}

	return val, nil
}

func (db *Database) Del(key string) {
	delete(db.Store, key)
}

func (db *Database) Exist(key string) bool {
	_, exist := db.Store[key]
	return exist
}

func HandleCommand(db *Database, command string) (string, error) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", err.ErrEmptyCommand
	}

	switch strings.ToUpper(parts[0]) {
	case "SET":
		if len(parts) != 3 {
			return "", err.ErrSetRequiresTwoArgs
		}

		db.Set(parts[1], parts[2])

		return "", nil

	case "GET":
		if len(parts) != 2 {
			return "", err.ErrGetRequiresOneArg
		}

		return db.Get(parts[1])

	case "DEL":
		if len(parts) != 2 {
			return "", err.ErrDelRequiresOneArg
		}
		db.Del(parts[1])

		return "", nil

	default:
		return "", err.ErrUnknownCommand
	}
}
