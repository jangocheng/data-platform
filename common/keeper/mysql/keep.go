package mysql

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func (k *Keeper) Watch(instance interface{}, tableName string, idx string) {
	ptrKind := reflect.TypeOf(instance).Kind().String()
	mapKind := reflect.TypeOf(instance).Elem().String()
	if ptrKind != "ptr" || !strings.Contains(mapKind, "map[string]") {
		panic("receive error format instance, the instance be kept must be kind map[string]Struct")
	}
	var jsonStatus string
	var firstSync bool
	go func() {
		time1 := time.Now().Unix()
		for {
			if k.killer.KillNow() {
				break
			}
			if !firstSync || int32(time.Now().Unix() - time1) >= k.internal {
				firstSync = true
				time1 = time.Now().Unix()
				elemType := reflect.TypeOf(instance).Elem().Elem()
				sliceType := reflect.SliceOf(elemType)
				slicePtr := reflect.New(sliceType)
				sliceInterface := slicePtr.Interface()
				if k.engine == nil {
					break
				}
				err := k.SycData(sliceInterface, tableName)
				if err != nil {
					fmt.Printf("fail to syc data from mysql: %s\n", err)
				}
				jsonStatusNew, err := json.Marshal(sliceInterface)
				if err != nil {
					fmt.Printf("fail to digest data to json: %s\n", err)
				}
				if string(jsonStatusNew) == jsonStatus {
					continue
				} else {
					jsonStatus = string(jsonStatusNew)
				}
				sliceElemValue := reflect.ValueOf(sliceInterface).Elem()
				mapReflection := reflect.MakeMapWithSize(reflect.TypeOf(instance).Elem(), sliceElemValue.Len())
				for i:=0; i<sliceElemValue.Len(); i++ {
					mapReflection.SetMapIndex(sliceElemValue.Index(i).FieldByName(idx), sliceElemValue.Index(i))
				}

				reflect.ValueOf(instance).Elem().Set(mapReflection)
			}
			time.Sleep(time.Second * 1)
		}
	}()
}

func (k *Keeper) SycData(instance interface{}, tableName string) error {
	err := k.engine.Table(tableName).Find(instance)
	if err != nil {
		return err
	}
	return nil
}