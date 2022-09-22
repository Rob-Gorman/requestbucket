package environment

import (
	"fmt"
	"os"
)

type Env struct {
	Logfile   string
	Pghost    string
	Pgport    string
	Pguser    string
	Password  string
	Pgdbname  string
	Table     string
	MongoUri  string
	Mongodb   string
	MongoColl string
}

func LoadDotenv() *Env {
	// gotenv.Load()
	// MONGO_URI := fmt.Sprintf("mongodb://localhost:27017/")
	MONGO_URI := fmt.Sprintf("mongodb://%s:%s/",
		os.Getenv("MONGODB_HOST"), os.Getenv("MONGODB_PORT"))
	env := &Env{
		Logfile:   os.Getenv("LOGFILE"),
		Pghost:    os.Getenv("PGHOST"),
		Pgport:    os.Getenv("PGPORT"),
		Pguser:    os.Getenv("PGUSER"),
		Password:  os.Getenv("PASSWORD"),
		Pgdbname:  os.Getenv("PGDATABASE"),
		Table:     os.Getenv("PGTABLE"),
		MongoUri:  MONGO_URI,
		Mongodb:   os.Getenv("MONGODB"),
		MongoColl: os.Getenv("MONGODB_COLL"),
	}
	return env
}
