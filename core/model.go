package core

import (
	"ablog/util"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var Model *coreModel

type coreModel struct {
	models map[string]coreModelMeta
}

type coreModelMeta struct {
	Name       string
	Pk         string // pk field, need int64
	PkKey      string
	Index      map[string]string
	IndexKeys  map[string]string
	Unique     map[string]string
	UniqueKeys map[string]string
}

func newCoreModel() *coreModel {
	return &coreModel{make(map[string]coreModelMeta)}
}

// register value
func (m *coreModel) Register(values ...interface{}) error {
	for _, value := range values {
		rType := reflect.TypeOf(value).Elem()
		fieldNum := rType.NumField()
		meta := coreModelMeta{
			Name:       rType.String(),
			Index:      make(map[string]string),
			IndexKeys:  make(map[string]string),
			Unique:     make(map[string]string),
			UniqueKeys: make(map[string]string),
		}
		for i := 0; i < fieldNum; i++ {
			rField := rType.Field(i)
			tag := rField.Tag.Get("model")
			if tag == "pk" && rField.Type.String() == "int64" {
				meta.Pk = rField.Name
				meta.PkKey = fmt.Sprintf("%s:%s", meta.Name, meta.Pk)
				continue
			}
			if tag == "index" {
				meta.Index[rField.Name] = rField.Type.String()
				if rField.Type.String() == "slice" && rField.Type.Elem().String() == "byte" {
					meta.Index[rField.Name] = "[]byte"
				}
				meta.IndexKeys[rField.Name] = fmt.Sprintf("%s:%s", meta.Name, rField.Name)
				continue
			}
			if tag == "unique" {
				meta.Unique[rField.Name] = rField.Type.String()
				if rField.Type.String() == "slice" && rField.Type.Elem().String() == "byte" {
					meta.Unique[rField.Name] = "[]byte"
				}
				meta.UniqueKeys[rField.Name] = fmt.Sprintf("%s:%s", meta.Name, rField.Name)
				continue
			}
		}
		if meta.Pk == "" {
			return errors.New(rType.String() + " need pk field")
		}
		m.models[rType.String()] = meta
	}
	return nil
}

// save data
func (m *coreModel) Save(v interface{}) error {
	// get or register meta
	rv := reflect.ValueOf(v).Elem()
	typeName := rv.Type().String()
	meta, ok := m.models[typeName]
	if !ok {
		if err := m.Register(v); err != nil {
			return err
		}
		meta = m.models[typeName]
	}
	// save data with pk value
	rfv := rv.FieldByName(meta.Pk)
	pkValue := rfv.Int()
	pkBytes := util.Int642Bytes(pkValue)

	key := fmt.Sprintf("%s:%d", meta.PkKey, pkValue)
	if err := Db.SetJson(key, v); err != nil {
		return err
	}

	// save pk list
	if err := Db.ZSet(
		meta.PkKey, pkValue, pkBytes); err != nil {
		return err
	}

	// save unique
	for uniqueKey, uniqueType := range meta.Unique {
		rfv := rv.FieldByName(uniqueKey)
		bytes := reflectValue2Bytes(rfv, uniqueType)
		// check unique
		if b, _ := Db.HGet(meta.UniqueKeys[uniqueKey], bytes); len(b) > 0 && util.Bytes2Int64(b) != pkValue {
			return errors.New(typeName + "'s " + uniqueKey + " is conflicted")
		}
		if err := Db.HSet(meta.UniqueKeys[uniqueKey], bytes, pkBytes); err != nil {
			return err
		}
	}

	// save index
	for indexKey, indexType := range meta.Index {
		rfv := rv.FieldByName(indexKey)
		bytes := reflectValue2Bytes(rfv, indexType)
		if err := Db.HSet(meta.IndexKeys[indexKey], bytes, pkBytes); err != nil {
			return err
		}
	}
	return nil
}

func (m *coreModel) Remove(v interface{}) error {
	// get or register meta
	rv := reflect.ValueOf(v).Elem()
	typeName := rv.Type().String()
	meta, ok := m.models[typeName]
	if !ok {
		if err := m.Register(v); err != nil {
			return err
		}
		meta = m.models[typeName]
	}

	// remove data with pk value
	rfv := rv.FieldByName(meta.Pk)
	pkValue := rfv.Int()
	pkBytes := util.Int642Bytes(pkValue)

	key := fmt.Sprintf("%s:%d", meta.PkKey, pkValue)
	if err := Db.Del(key); err != nil {
		return err
	}

	// delete pk list
	if err := Db.ZDel(
		meta.PkKey, pkBytes); err != nil {
		return err
	}

	// delete unique
	for uniqueKey, uniqueType := range meta.Unique {
		rfv := rv.FieldByName(uniqueKey)
		bytes := reflectValue2Bytes(rfv, uniqueType)
		if err := Db.HDel(meta.UniqueKeys[uniqueKey], bytes); err != nil {
			return err
		}
	}

	// delete index
	for indexKey, indexType := range meta.Unique {
		rfv := rv.FieldByName(indexKey)
		bytes := reflectValue2Bytes(rfv, indexType)
		if err := Db.HDel(meta.IndexKeys[indexKey], bytes); err != nil {
			return err
		}
	}
	return nil
}

func (m *coreModel) Get(v interface{}) error {
	// get meta
	rv := reflect.ValueOf(v).Elem()
	meta, ok := m.models[rv.Type().String()]
	if !ok {
		return errors.New("unregistered value")
	}

	// get pk value
	rfv := rv.FieldByName(meta.Pk)
	pkValue := rfv.Int()

	// to value
	key := fmt.Sprintf("%s:%d", meta.PkKey, pkValue)
	return Db.GetJson(key, v)
}

func (m *coreModel) GetBy(v interface{}, field string) error {
	// get meta
	rv := reflect.ValueOf(v).Elem()
	typeName := rv.Type().String()
	meta, ok := m.models[typeName]
	if !ok {
		return errors.New("unregistered value")
	}

	// get field type
	fType := meta.Unique[field]
	if fType == "" {
		fType = meta.Index[field]
	}
	if fType == "" {
		return errors.New("un-indexed field " + field)
	}

	rfv := rv.FieldByName(field)
	bytes := reflectValue2Bytes(rfv, fType)
	key := fmt.Sprintf("%s:%s", typeName, field)
	resultBytes, err := Db.HGet(key, bytes)
	if err != nil {
		return err
	}
	id := util.Bytes2Int64(resultBytes)
	rv.FieldByName(meta.Pk).SetInt(id)
	reflect.ValueOf(v).Elem().Set(rv)
	return m.Get(v)
}

func reflectValue2Bytes(rv reflect.Value, realType string) []byte {
	if realType == "[]byte" {
		return rv.Bytes()
	}
	if realType == "string" {
		return []byte(rv.String())
	}
	if realType == "int" || realType == "int64" {
		return util.Int642Bytes(rv.Int())
	}
	if realType == "float64" {
		return util.Float2Bytes(rv.Float())
	}
	if realType == "bool" {
		if rv.Bool() {
			return []byte("true")
		}
		return []byte("false")
	}
	b, _ := json.Marshal(rv.Interface())
	return b
}
