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
	models map[string]*coreModelMeta
}

type coreModelMeta struct {
	Name string

	Pk    string // pk field, need int64
	PkKey string

	Index     map[string]string
	IndexKeys map[string]string

	Unique     map[string]string
	UniqueKeys map[string]string

	One2One   map[string]string
	IsOne2One bool

	One2N map[string]string
	N2One []string
}

func newCoreModel() *coreModel {
	return &coreModel{make(map[string]*coreModelMeta)}
}

func (m *coreModel) registerType(rType reflect.Type, isPk bool) error {
	fieldNum := rType.NumField()
	meta := coreModelMeta{
		Name: rType.String(),
	}
	for i := 0; i < fieldNum; i++ {
		rField := rType.Field(i)
		tag := rField.Tag.Get("model")

		// pk tag
		if tag == "pk" && rField.Type.Kind() == reflect.Int64 {
			meta.Pk = rField.Name
			meta.PkKey = fmt.Sprintf("%s:%s", meta.Name, meta.Pk)
			continue
		}

		// unique
		if tag == "unique" {
			if len(meta.Unique) == 0 {
				meta.Unique = make(map[string]string)
				meta.UniqueKeys = make(map[string]string)
			}
			meta.Unique[rField.Name] = rField.Type.String()
			if rField.Type.Kind() == reflect.Slice && rField.Type.Elem().String() == "byte" {
				meta.Unique[rField.Name] = "[]byte"
			}
			meta.UniqueKeys[rField.Name] = fmt.Sprintf("%s:%s", meta.Name, rField.Name)
			continue
		}

		// index tag
		if tag == "index" {
			if len(meta.Index) == 0 {
				meta.Index = make(map[string]string)
				meta.IndexKeys = make(map[string]string)
			}
			meta.Index[rField.Name] = rField.Type.String()
			if rField.Type.Kind() == reflect.Slice && rField.Type.Elem().String() == "byte" {
				meta.Index[rField.Name] = "[]byte"
			}
			meta.IndexKeys[rField.Name] = fmt.Sprintf("%s:%s", meta.Name, rField.Name)
			continue
		}

		// 1-1 tag
		if tag == "1-1" {
			if rField.Type.Kind() != reflect.Ptr || rField.Type.Elem().Kind() != reflect.Struct {
				return errors.New(rType.String() + "'s 1-1 field need struct point")
			}
			if len(meta.One2One) == 0 {
				meta.One2One = make(map[string]string)
			}
			m.registerType(rField.Type.Elem(), false)

			name := rField.Type.Elem().String()
			meta.One2One[rField.Name] = name
			m.models[name].IsOne2One = true
		}

		// 1-n tag
		if tag == "1-n" {
			if rField.Type.Kind() != reflect.Slice || rField.Type.Elem().Kind() != reflect.Ptr || rField.Type.Elem().Elem().Kind() != reflect.Struct {
				return errors.New(rType.String() + "'s 1-1 field need struct point slice")
			}
			if len(meta.One2N) == 0 {
				meta.One2N = make(map[string]string)
			}
			m.registerType(rField.Type.Elem().Elem(), true)

			name := rField.Type.Elem().Elem().String()
			meta.One2N[rField.Name] = name

			// 1-n means n-1 to positive meta
			if len(m.models[name].N2One) == 0 {
				m.models[name].N2One = []string{meta.Name}
			} else {
				m.models[name].N2One = append(m.models[name].N2One, meta.Name)
			}
		}
	}
	if isPk && meta.Pk == "" {
		return errors.New(rType.String() + " need pk field")
	}
	m.models[rType.String()] = &meta
	return nil
}

// register value
func (m *coreModel) Register(values ...interface{}) error {
	for _, value := range values {
		rType := reflect.TypeOf(value).Elem()
		m.registerType(rType, true)
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

	// get pk value and bytes
	rfv := rv.FieldByName(meta.Pk)
	pkValue := rfv.Int()
	pkBytes := util.Int642Bytes(pkValue)

	// save unique, check unique first
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

	// save data with pk value
	key := fmt.Sprintf("%s:%d", meta.PkKey, pkValue)
	if err := Db.SetJson(key, v); err != nil {
		return err
	}

	// save pk list
	if err := Db.ZSet(
		meta.PkKey, pkValue, pkBytes); err != nil {
		return err
	}

	// save index
	for indexKey, indexType := range meta.Index {
		rfv := rv.FieldByName(indexKey)
		bytes := reflectValue2Bytes(rfv, indexType)
		if err := Db.HSet(meta.IndexKeys[indexKey], bytes, pkBytes); err != nil {
			return err
		}
	}

	// save 1-1
	for o2oKey, o2oMeta := range meta.One2One {
		rfv := rv.FieldByName(o2oKey)
		// only save not null value
		if !rfv.IsNil() {
			key := fmt.Sprintf(o2oMeta+":%d", pkValue)
			if err := Db.SetJson(key, rfv.Interface()); err != nil {
				return err
			}
		}
	}

	// save 1-n
	for o2nKey, o2nMeta := range meta.One2N {
		// get value
		rfv := rv.FieldByName(o2nKey)
		if rfv.IsNil() {
			continue
		}
		// get slice length
		l := rfv.Len()
		if l == 0 {
			continue
		}

		for i := 0; i < l; i++ {
			rfiv := rfv.Index(i)
			m.Save(rfiv.Interface())
			rfiv2 := rfiv.Elem().FieldByName(m.models[o2nMeta].Pk)
			pkValue2 := rfiv2.Int()

			key1 := fmt.Sprintf("%s:%d:%s", typeName, pkValue, o2nMeta)
			if err := Db.ZSet(key1, pkValue2, util.Int642Bytes(pkValue2)); err != nil {
				return err
			}
			key2 := fmt.Sprintf("%s:%d:%s", o2nMeta, pkValue2, typeName)
			if err := Db.ZSet(key2, pkValue, util.Int642Bytes(pkValue)); err != nil {
				return err
			}
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
	if meta.IsOne2One {
		return errors.New(typeName + " belongs to another model")
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
	for indexKey, indexType := range meta.Index {
		rfv := rv.FieldByName(indexKey)
		bytes := reflectValue2Bytes(rfv, indexType)
		if err := Db.HDel(meta.IndexKeys[indexKey], bytes); err != nil {
			return err
		}
	}

	// delete 1-1
	for _, o2oMeta := range meta.One2One {
		key := fmt.Sprintf(o2oMeta+":%d", pkValue)
		if err := Db.Del(key); err != nil {
			return err
		}
	}

	// delete 1-n
	for _, o2nMeta := range meta.One2N {
		key := fmt.Sprintf("%s:%d:%s", typeName, pkValue, o2nMeta)
		all, err := Db.ZAllAsc(key)
		if err != nil {
			return err
		}
		for _, a := range all {
			key2 := fmt.Sprintf("%s:%d:%s", o2nMeta, a.Score, typeName)
			if err := Db.ZDel(key2, pkBytes); err != nil {
				return err
			}
		}
		if err := Db.ZClear(key); err != nil {
			return err
		}
	}

	// delete n-1
	for _, n2oMeta := range meta.N2One {
		key := fmt.Sprintf("%s:%d:%s", typeName, pkValue, n2oMeta)
		all, err := Db.ZAllAsc(key)
		if err != nil {
			return err
		}
		for _, a := range all {
			key2 := fmt.Sprintf("%s:%d:%s", n2oMeta, a.Score, typeName)
			if err := Db.ZDel(key2, pkBytes); err != nil {
				return err
			}
		}
		if err := Db.ZClear(key); err != nil {
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
	if err := Db.GetJson(key, v); err != nil {
		return err
	}

	// get 1-1 data
	if len(meta.One2One) > 0 {
		for o2oKey, o2oMeta := range meta.One2One {
			key := fmt.Sprintf(o2oMeta+":%d", pkValue)
			rfv := rv.FieldByName(o2oKey)
			rv2 := reflect.New(rfv.Type().Elem())
			if err := Db.GetJson(key, rv2.Interface()); err != nil {
				return err
			}
			rfv.Set(rv2)
		}
	}

	// get 1-n data
	if len(meta.One2N) > 0 {
		for o2nKey, o2nMeta := range meta.One2N {
			key := fmt.Sprintf("%s:%d:%s", meta.Name, pkValue, o2nMeta)
			rfv := rv.FieldByName(o2nKey)
			all, err := Db.ZAllAsc(key)
			if err != nil {
				return err
			}
			meta2 := m.models[o2nMeta]
			sliceRv := reflect.New(reflect.SliceOf(rfv.Type().Elem())).Elem()
			for _, a := range all {
				rv2 := reflect.New(rfv.Type().Elem().Elem())
				key := fmt.Sprintf("%s:%s:%d", o2nMeta, meta2.Pk, a.Score)
				if err := Db.GetJson(key, rv2.Interface()); err != nil {
					return err
				}
				sliceRv = reflect.Append(sliceRv, rv2)
			}
			rfv.Set(sliceRv)
		}
	}
	return nil
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
