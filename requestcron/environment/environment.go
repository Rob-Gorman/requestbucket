package environment

import (
	"os"

	"github.com/subosito/gotenv"
)

type Env struct {
	Logfile   string
	Host      string
	Port      string
	User      string
	Password  string
	Pgdbname  string
	Table     string
	MongoUri  string
	Mongodb   string
	MongoColl string
}

func LoadDotenv() *Env {
	gotenv.Load("./environment/.env")
	env := &Env{
		Logfile:   os.Getenv("LOGFILE"),
		Host:      os.Getenv("HOST"),
		Port:      os.Getenv("PORT"),
		User:      os.Getenv("USER"),
		Password:  os.Getenv("PASSWORD"),
		Pgdbname:  os.Getenv("PGDBNAME"),
		Table:     os.Getenv("PGTABLE"),
		MongoUri:  os.Getenv("MONGODB_URI"),
		Mongodb:   os.Getenv("MONGODB"),
		MongoColl: os.Getenv("MONGODB_COLL"),
	}
	return env
}
