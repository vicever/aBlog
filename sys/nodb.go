package sys

import (
	"encoding/binary"
	"encoding/json"
	"github.com/lunny/nodb"
	"github.com/lunny/nodb/config"
	"math"
)

var NoDb *nodbClient

type nodbClient struct {
	db    *nodb.Nodb
	dbSet *nodb.DB
}

func newNodbClient() *nodbClient {
	cfg := config.NewConfigDefault()
	cfg.DataDir = Config.NoDb.Directory

	db, err := nodb.Open(cfg)
	if err != nil {
		Fatal("[NoDB] open database error : %v", err)
		return nil
	}

	dbSet, err := db.Select(Config.NoDb.Index)
	if err != nil {
		Fatal("[NoDB] open database error : %v", err)
		return nil
	}

	return &nodbClient{db, dbSet}

}

// bytes
func (c *nodbClient) Set(key string, value []byte) error {
	return c.dbSet.Set([]byte(key), value)
}

func (c *nodbClient) Get(key string) ([]byte, error) {
	return c.dbSet.Get([]byte(key))
}

// int
func (c *nodbClient) SetInt64(key string, value int64) error {
	bytes := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(bytes, value)
	return c.dbSet.Set([]byte(key), bytes)
}

func (c *nodbClient) SetInt(key string, value int) error {
	return c.SetInt64(key, int64(value))
}

func (c *nodbClient) GetInt64(key string) (int64, error) {
	bytes, err := c.Get(key)
	if err != nil {
		return 0, err
	}
	value, _ := binary.Varint(bytes)
	return value, nil
}

func (c *nodbClient) GetInt(key string) (int, error) {
	value, err := c.GetInt64(key)
	return int(value), err
}

// bool
func (c *nodbClient) SetBool(key string, value bool) error {
	var bytes []byte
	if value {
		bytes = []byte("true")
	} else {
		bytes = []byte("false")
	}
	return c.Set(key, bytes)
}

func (c *nodbClient) GetBool(key string) (bool, error) {
	bytes, err := c.Get(key)
	if err != nil {
		return false, err
	}
	return string(bytes) == "true", nil
}

/*
func Float64frombytes(bytes []byte) float64 {
    bits := binary.LittleEndian.Uint64(bytes)
    float := math.Float64frombits(bits)
    return float
}

func Float64bytes(float float64) []byte {
    bits := math.Float64bits(float)
    bytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(bytes, bits)
    return bytes
}
*/

// float
func (c *nodbClient) SetFloat(key string, value float64) error {
	bits := math.Float64bits(value)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return c.dbSet.Set([]byte(key), bytes)
}

func (c *nodbClient) GetFloat(key string) (float64, error) {
	bytes, err := c.Get(key)
	if err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float, nil
}

// json
func (c *nodbClient) SetJson(key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(key, bytes)
}

func (c *nodbClient) GetJson(key string, value interface{}) error {
	bytes, err := c.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, value)
}

// set expire
func (c *nodbClient) SExpire(key string, duration int64) error {
	_, err := c.dbSet.SExpire([]byte(key), duration)
	return err
}

// check exist
func (c *nodbClient) Exist(key string) bool {
	v, err := c.dbSet.Exists([]byte(key))
	if err != nil {
		return false
	}
	return v > 0
}

// del key
func (c *nodbClient) Del(key string) error {
	_, err := c.dbSet.Del([]byte(key))
	return err
}
