package database

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	err "github.com/baxromumarov/redis-go/errors"
)

type Database struct {
	Store       map[string]string
	Mu          sync.RWMutex
	Expirations map[string]time.Time
}

func Init() *Database {
	return &Database{
		Store:       make(map[string]string),
		Expirations: make(map[string]time.Time),
		Mu:          sync.RWMutex{},
	}
}

// ttl always in seconds
func (db *Database) Set(key, val string, ttl int) error {
	if key == "" {
		return err.ErrEmptyKey
	}
	if val == "" {
		return err.ErrEmpyVal
	}
	db.Mu.Lock()
	defer db.Mu.Unlock()

	db.Store[key] = val
	if ttl > 0 {
		expirationTime := time.Now().Add(time.Duration(ttl) * time.Second)
		db.Expirations[key] = expirationTime
	} else {
		// If no expiration time is provided, remove the expiration entry (persistent)
		delete(db.Expirations, key)
	}
	return nil
}

func (db *Database) Get(key string) (string, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	if expiration, exists := db.Expirations[key]; exists && time.Now().After(expiration) {
		db.Del(key)
		return "", err.ErrKeyHasExpired
	}

	val, exist := db.Store[key]
	if !exist {
		return "", err.ErrContentNotFound
	}

	return val, nil
}

func (db *Database) Del(key string) {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	delete(db.Store, key)
}

func (db *Database) Exist(key string) bool {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	_, exist := db.Store[key]
	return exist
}

func (db *Database) Expire(key string, seconds int) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	if _, exists := db.Store[key]; !exists {
		return err.ErrKeyNotFound
	}

	db.Expirations[key] = time.Now().Add(time.Duration(seconds) * time.Second)

	return nil
}

func (db *Database) StartExpirationCleaner(interval time.Duration) {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	go func() {
		for {
			time.Sleep(interval)
			for key, expiration := range db.Expirations {
				if time.Now().After(expiration) {
					db.Del(key)
				}
			}
		}
	}()
}

// TTL returns the remaining time to live of a key
func (db *Database) TTL(key string) (string, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	expiration, exists := db.Expirations[key]
	if !exists {
		return "", err.ErrNoExpirationSet
	}

	remaining := time.Until(expiration)
	if remaining <= 0 {
		db.Del(key)
		return "", err.ErrKeyHasExpired
	}
	return fmt.Sprintf("%f", remaining.Seconds()), nil
}

func HandleCommand(db *Database, command string) (string, error) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", err.ErrEmptyCommand
	}

	switch strings.ToUpper(parts[0]) {

	case "SET":
		if len(parts) < 3 {
			return "", err.ErrSETCommand
		}

		key := parts[1]
		val := parts[2]
		ttl := 0 // Default TTL: No expiration

		if len(parts) == 4 {
			fmt.Println(parts)
			// If a fourth argument is provided, treat it as the TTL
			ttlParsed, error := strconv.Atoi(parts[3])
			if error != nil {
				return "", err.ErrInvalidTTLVal
			}

			ttl = ttlParsed
		}

		db.Set(key, val, ttl)

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

	case "TTL":
		if len(parts) != 2 {
			return "", err.ErrTTLRequiresOneArg
		}

		return db.TTL(parts[1])

	default:
		return "", err.ErrUnknownCommand
	}
}
