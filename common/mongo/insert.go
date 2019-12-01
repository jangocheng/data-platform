package mongo

import "github.com/pkg/errors"

func (c *ClientMongo) InsertOne(docs interface{}, collection ... string) error {
	if len(collection) == 0 && c.collection == nil {
		return errors.New("have not assign collection")
	}
	var err error
	if len(collection) != 0 {
		err = c.dSession.C(collection[0]).Insert(docs)
	} else {
		err = c.collection.Insert(docs)
	}
	if err != nil {
		return errors.Wrap(err, "fail to insert items")
	}
	return nil
}


func (c *ClientMongo) Insert(docs []interface{}, collection ... string) error {
	if len(collection) == 0 && c.collection == nil {
		return errors.New("have not assign collection")
	}
	var err error
	if len(collection) != 0 {
		err = c.dSession.C(collection[0]).Insert(docs...)
	} else {
		err = c.collection.Insert(docs...)
	}
	if err != nil {
		return errors.Wrap(err, "fail to insert items")
	}
	return nil
}