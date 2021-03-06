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
		Log.Fatal("open nosql fail : %v", err)
		return nil
	}

	// select db
	dbSet, err := db.Select(opt.Index)
	if err != nil {
		Log.Fatal("select nosql fail : %v", err)
		return nil
	}
	return &coreDb{dbSet}
}

/*
==== kv pair
*/

func (db *coreDb) Set(key string, value []byte) error {
	Log.Debug("[Db -> Set][%s]", key)
	return db.Db.Set([]byte(key), value)
}

func (db *coreDb) Get(key string) ([]byte, error) {
	Log.Debug("[Db -> Get][%s]", key)
	return db.Db.Get([]byte(key))
}

func (db *coreDb) SetExpire(key string, duration int64) error {
	Log.Debug("[Db -> SetEx][%s][%d]", key, duration)
	_, err := db.Db.SExpire([]byte(key), duration)
	return err
}

func (db *coreDb) Del(key string) error {
	Log.Debug("[Db -> Del][%s]", key)
	_, err := db.Db.Del([]byte(key))
	return err
}

func (db *coreDb) Exist(key string) bool {
	Log.Debug("[Db -> Exist][%s]", key)
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

/*
===== hset
*/

func (db *coreDb) HSet(key1 string, key2, value []byte) error {
	Log.Debug("[Db -> HSet][%s:%v]", key1, key2)
	_, err := db.Db.HSet([]byte(key1), []byte(key2), value)
	return err
}

func (db *coreDb) HGet(key1 string, key2 []byte) ([]byte, error) {
	Log.Debug("[Db -> HGet][%s:%v]", key1, key2)
	return db.Db.HGet([]byte(key1), []byte(key2))
}

func (db *coreDb) HDel(key1 string, key2 []byte) error {
	Log.Debug("[Db -> HDel][%s:%v]", key1, key2)
	_, err := db.Db.HDel([]byte(key1), []byte(key2))
	return err
}

func (db *coreDb) HExist(key1 string, key2 []byte) bool {
	Log.Debug("[Db -> HExist][%s:%v]", key1, key2)
	bytes, _ := db.HGet(key1, key2)
	if len(bytes) == 0 {
		return false
	}
	return true
}

func (db *coreDb) HGetAll(key1 string) ([]nodb.FVPair, error) {
	Log.Debug("[Db -> HGetAll][%s]", key1)
	return db.Db.HGetAll([]byte(key1))
}

func (db *coreDb) HClear(key1 string) error {
	Log.Debug("[Db -> HClear][%s]", key1)
	_, err := db.Db.HClear([]byte(key1))
	return err
}

/*
===== zset
*/

func (db *coreDb) ZSet(key string, number int64, value []byte) error {
	Log.Debug("[Db -> ZSet][%s:%v]", key, value)
	scorePair := nodb.ScorePair{number, value}
	_, err := db.Db.ZAdd([]byte(key), scorePair)
	return err
}

func (db *coreDb) ZDel(key string, value []byte) error {
	Log.Debug("[Db -> ZDel][%s:%v]", key, value)
	_, err := db.Db.ZRem([]byte(key), value)
	return err
}

func (db *coreDb) ZClear(key string) error {
	Log.Debug("[Db -> ZClear][%s]", key)
	_, err := db.Db.ZClear([]byte(key))
	return err
}

func (db *coreDb) ZCount(key string) int64 {
	Log.Debug("[Db -> ZCount][%s]", key)
	count, _ := db.Db.ZCard([]byte(key))
	return count
}

func (db *coreDb) ZAllAsc(key string) ([]nodb.ScorePair, error) {
	Log.Debug("[Db -> ZAllAsc][%s]", key)
	return db.Db.ZRange([]byte(key), 0, -1)
}

func (db *coreDb) ZAllDesc(key string) ([]nodb.ScorePair, error) {
	Log.Debug("[Db -> ZAllDesc][%s]", key)
	return db.Db.ZRevRange([]byte(key), 0, -1)
}

func (db *coreDb) ZPageAsc(key string, page, size int) ([]nodb.ScorePair, error) {
	Log.Debug("[Db -> ZPageAsc][%s][%d,%d]", key, page, size)
	all := int(db.ZCount(key))
	if all == 0 {
		return nil, nil
	}

	// begin index
	begin := (page - 1) * size
	if begin < 0 {
		begin = 0
	}

	// end index
	end := page*size - 1
	if end > all {
		end = all
	}
	return db.Db.ZRange([]byte(key), begin, end)
}

func (db *coreDb) ZPageDesc(key string, page, size int) ([]nodb.ScorePair, error) {
	Log.Debug("[Db -> ZPageDesc][%s][%d,%d]", key, page, size)
	all := int(db.ZCount(key))
	if all == 0 {
		return nil, nil
	}

	// begin index
	begin := (page - 1) * size
	if begin < 0 {
		begin = 0
	}

	// end index
	end := page * size
	if end > all {
		end = all
	}
	return db.Db.ZRevRange([]byte(key), begin, end)
}
