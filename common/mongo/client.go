package mongo

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"strings"
)

type ClientMongo struct {
	session 		*mgo.Session
	dSession         *mgo.Database
	collection  	*mgo.Collection
	db          	string
	collectionName  string
}

type ConnectConfig struct {
	IpPortList      []string
	UserName        string
	Password 		string
	Collection      string
	DB              string
}


func NewClient(cnf ConnectConfig) (*ClientMongo, error) {
	//dailInfo := &mgo.DialInfo{
	//	Addrs:   		cnf.IpPortList,
	//	Direct:  		false,
	//	Database: 		cnf.DB,
	//	Username:       cnf.UserName,
	//	Password:       cnf.Password,
	//	PoolLimit:      cnf.PoolLimit,
	//	Timeout:        16,
	//}
	//tlsConfig := &tls.Config{InsecureSkipVerify:cnf.SkipVerifyTLS}
	//dailInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
	//	conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
	//	return conn, err
	//}
	//session, err := mgo.DialWithInfo(dailInfo)
	mgoUrl := fmt.Sprintf("mongodb://%s:%s@%s/%s",
		cnf.UserName, cnf.Password,
		strings.Join(cnf.IpPortList, ","), cnf.DB)
	session, err := mgo.Dial(mgoUrl)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create mongo client")
	}
	session.SetMode(mgo.Monotonic, true)
	dSession := session.DB(cnf.DB)
	var collection *mgo.Collection
	if cnf.Collection != "" {
		collection = dSession.C(cnf.Collection)
	}
	return &ClientMongo{collection:collection, db: cnf.DB, collectionName:cnf.Collection,
		session:session, dSession: dSession}, nil

}

func (c *ClientMongo) Close() {
	if c.session != nil {
		c.session.Close()
	}
}
