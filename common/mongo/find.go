package mongo

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

func (c *ClientMongo) FindOne(condition map[string]interface{}, result interface{}, collection ... string) error {
	if len(collection) == 0 && c.collection == nil {
		return errors.New("have not assign collection")
	}

	instanceKind := reflect.TypeOf(result).Kind().String()
	if instanceKind != "Struct" && instanceKind != "struct" && !strings.HasPrefix(instanceKind, "map") {
		return errors.New(fmt.Sprintf("receive wrong result format %s", instanceKind))
	}
	var err error
	if len(collection) != 0 {
		err = c.dSession.C(collection[0]).Find(condition).One(result)
	} else {
		err = c.collection.Find(condition).One(result)
	}


	if err != nil {
		return errors.Wrap(err, "fail to find item")
	}
	return nil
}


func (c *ClientMongo) Find(condition map[string]interface{}, results interface{}, collection ... string) error {
	if len(collection) == 0 && c.collection == nil {
		return errors.New("have not assign collection")
	}

	var err error
	if len(collection) != 0 {
		err = c.dSession.C(collection[0]).Find(condition).All(results)
	} else {
		err = c.collection.Find(condition).All(results)
	}
	if err != nil {
		return errors.Wrap(err, "fail to find items")
	}
	return nil
}


func (c *ClientMongo) FindRange(condition map[string]interface{}, results interface{}, offSet int, limit int,
	collection ... string) error {
	if len(collection) == 0 && c.collection == nil {
		return errors.New("have not assign collection")
	}
	if limit == 0 {
		return errors.New("the limit can not be 0")
	}

	var err error
	if len(collection) != 0 {
		err = c.dSession.C(collection[0]).Find(condition).Skip(offSet).Limit(limit).All(results)
	} else {
		err = c.collection.Find(condition).Skip(offSet).Limit(limit).All(results)
	}
	if err != nil {
		return errors.Wrap(err, "fail to find items")
	}
	return nil
}

