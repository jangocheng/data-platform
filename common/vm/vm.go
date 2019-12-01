package vm

type ClientVM interface {

	Run(string, ...interface{}) (Value, error)

	SetFunc(string, interface{}) error

	Close()

}

type Value interface {

	GetString() string

	GetFloat() (float64, error)

	GetObjectList() ([]Value, error)

	GetObjectKey(key string) (Value, error)

	GetInt() (int64, error)

	GetBool() (bool, error)

	IsNil() bool

}
