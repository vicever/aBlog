package core

import (
	"encoding/binary"
	"encoding/json"
	"github.com/lunny/nodb"
	"github.com/lunny/nodb/config"
	"math"
)

var (
	Db *noDbClient
)

type noDbClient struct {
	Db    *nodb.Nodb
	DbSet *nodb.DB
}

func init() {
	ConnectDb()

}

func ConnectDb() {
	cfg := config.NewConfigDefault()
	cfg.DataDir = Vars.DataDirectory

	NoDb, err := nodb.Open(cfg)
	if err != nil {
		panic(err)
	}

	DbSet, err := NoDb.Select(Vars.DbIndex)
	if err != nil {
		panic(err)
	}

	Db = &noDbClient{NoDb, DbSet}
}

// bytes
func (c *noDbClient) Set(key string, value []byte) error {
	return c.DbSet.Set([]byte(key), value)
}

func (c *noDbClient) Get(key string) ([]byte, error) {
	return c.DbSet.Get([]byte(key))
}

// int
func (c *noDbClient) SetInt64(key string, value int64) error {
	bytes := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(bytes, value)
	return c.DbSet.Set([]byte(key), bytes)
}

func (c *noDbClient) SetInt(key string, value int) error {
	return c.SetInt64(key, int64(value))
}

func (c *noDbClient) GetInt64(key string) (int64, error) {
	bytes, err := c.Get(key)
	if err != nil {
		return 0, err
	}
	value, _ := binary.Varint(bytes)
	return value, nil
}

func (c *noDbClient) GetInt(key string) (int, error) {
	value, err := c.GetInt64(key)
	return int(value), err
}

// bool
func (c *noDbClient) SetBool(key string, value bool) error {
	var bytes []byte
	if value {
		bytes = []byte("true")
	} else {
		bytes = []byte("false")
	}
	return c.Set(key, bytes)
}

func (c *noDbClient) GetBool(key string) (bool, error) {
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
func (c *noDbClient) SetFloat(key string, value float64) error {
	bits := math.Float64bits(value)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return c.DbSet.Set([]byte(key), bytes)
}

func (c *noDbClient) GetFloat(key string) (float64, error) {
	bytes, err := c.Get(key)
	if err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float, nil
}

// json
func (c *noDbClient) SetJson(key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(key, bytes)
}

func (c *noDbClient) GetJson(key string, value interface{}) error {
	bytes, err := c.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, value)
}

// set expire
func (c *noDbClient) SExpire(key string, duration int64) error {
	_, err := c.DbSet.SExpire([]byte(key), duration)
	return err
}

// check exist
func (c *noDbClient) Exist(key string) bool {
	v, err := c.DbSet.Exists([]byte(key))
	if err != nil {
		return false
	}
	return v > 0
}
