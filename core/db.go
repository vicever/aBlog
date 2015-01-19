package core

import (
	"encoding/json"
	"github.com/lunny/nodb"
	"github.com/lunny/nodb/config"
)

type coreDb struct {
	Db *nodb.DB
}

func newCoreDb(opt dbConfig) *coreDb {
	// init config
	cfg := config.NewConfigDefault()
	cfg.DataDir = opt.Directory

	// init nosql
	db, err := nodb.Open(cfg)
	if err != nil {
		return nil
	}

	// select db
	dbSet, err := db.Select(opt.Index)
	if err != nil {
		return nil
	}
	return &coreDb{dbSet}
}

func (db *coreDb) Set(key string, value []byte) error {
	return db.Db.Set([]byte(key), value)
}

func (db *coreDb) Get(key string) ([]byte, error) {
	return db.Db.Get([]byte(key))
}

func (db *coreDb) SetExpire(key string, duration int64) error {
	_, err := db.Db.SExpire([]byte(key), duration)
	return err
}

func (db *coreDb) Del(key string) error {
	_, err := db.Db.Del([]byte(key))
	return err
}

func (db *coreDb) Exist(key string) bool {
	if i, _ := db.Db.Exists([]byte(key)); i > 0 {
		return true
	}
	return false
}

func (db *coreDb) SetJson(key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return db.Set(key, bytes)
}

func (db *coreDb) GetJson(key string, value interface{}) error {
	bytes, err := db.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, value)
}

func (db *coreDb) HSet(key1, key2 string, value []byte) error {
	_, err := db.Db.HSet([]byte(key1), []byte(key2), value)
	return err
}

func (db *coreDb) HGet(key1, key2 string) ([]byte, error) {
	return db.Db.HGet([]byte(key1), []byte(key2))
}

func (db *coreDb) HGetAll(key1 string) (map[string][]byte, error) {
	data := make(map[string][]byte)
	pValues, err := db.Db.HGetAll([]byte(key1))
	if err != nil {
		return data, err
	}
	for _, v := range pValues {
		data[string(v.Field)] = v.Value
	}
	return data, nil
}

func (db *coreDb) HClear(key1 string) error {
	_, err := db.Db.HClear([]byte(key1))
	return err
}
