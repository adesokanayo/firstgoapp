package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mongodbSession *mgo.Session

func init() {

	RootCMD.PersistentFlags().String("mongodb_uri", "localhost", "where the db is hosted")
	viper.BindPFlag("mongodb_uri", RootCMD.PersistentFlags().Lookup("mongodb_uri"))
	CreateUniqueIndexes()

}

//DBSession is used to create mongo db session.
func DBSession() *mgo.Session {
	if mongodbSession == nil {
		uri := os.Getenv("MONGODB_URI")
		if uri == "" {
			uri = viper.GetString("mongodb_uri")

			if uri == "" {
				log.Fatalln("No connection uri for MongoDB provided")
			}
		}

		var err error
		mongodbSession, err = mgo.Dial(uri)
		if mongodbSession == nil || err != nil {
			log.Fatalf("Can't connect to mongo, go error %v\n", err)
		}
		mongodbSession.SetSafe(&mgo.Safe{})
	}
	return mongodbSession

}

//DB is Database name
func DB() *mgo.Database {
	return DBSession().DB(viper.GetString("dbname"))
}

//Items functions calls the database collection channels.
func Items() *mgo.Collection {
	return DB().C("items")
}

//Channels functions calls the database Collection channels.
func Channels() *mgo.Collection {
	return DB().C("channels")
}

// CreateUniqueIndexes is used to crete unique indexes for all db operations.
func CreateUniqueIndexes() {

	idx := mgo.Index{

		Key:        []string{"key"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	if err := Items().EnsureIndex(idx); err != nil {
		fmt.Println(err)
	}

	if err := Channels().EnsureIndex(idx); err != nil {
		fmt.Println(err)
	}
}

func AllChannels() []Chnl {
	var channels []Chnl
	r := Channels().Find(bson.M{}).Sort("-lastbuilddate")
	r.All(&channels)
	return channels
}
