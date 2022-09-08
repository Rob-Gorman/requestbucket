package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"requestbucket/environment"
	"time"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type bucket struct {
	id         *int
	url        *string
	created_at *time.Time
}

var env *environment.Env = environment.LoadDotenv()

func main() {
	fmt.Printf("%+v", env)
	fmt.Print(env.MongoUri)
	log := retrieveLog()
	defer log.Close()

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", env.Host, env.Port, env.User, env.Pgdbname)
	pgdb, err2 := sql.Open("postgres", psqlconn)
	CheckError(err2)
	defer pgdb.Close()

	cullPGBuckets(pgdb, log)
}

func retrieveLog() *os.File {
	log, err1 := os.OpenFile(env.Logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	CheckError(err1)
	return log
}

func writeLog(log *os.File, value string) {
	a, b, c := time.Now().Clock()
	log.WriteString(fmt.Sprintf("%d:%d:%d DELETED:%v\n", a, b, c, value))
}

func cullPGBuckets(db *sql.DB, log *os.File) {
	for rowCount(db, env.Table) > 0 {
		selectstmt := fmt.Sprintf("SELECT * FROM %s order by id limit 1", env.Table)
		oldestRow := bucket{}
		err3 := db.QueryRow(selectstmt).Scan(&oldestRow.id, &oldestRow.url, &oldestRow.created_at)
		CheckError(err3)

		if time.Since(*oldestRow.created_at).Hours() > 48 {
			bucketId := *oldestRow.id
			cullPGRequests(db, bucketId)
			removeRow(db, bucketId)
			writeLog(log, *oldestRow.url)
		} else {
			break
		}
	}
}

func cullPGRequests(db *sql.DB, bucketId int) {
	removeMongoIds(db, bucketId)
	deletestmt := fmt.Sprintf("DELETE FROM %s WHERE bucket_id=%d", "requests", bucketId)
	fmt.Println(deletestmt)
	_, err := db.Exec(deletestmt)
	CheckError(err)
}

func removeMongoIds(db *sql.DB, bucketId int) {
	querystmt := fmt.Sprintf("SELECT mongo_document_ref FROM requests WHERE bucket_id=%d", bucketId)
	mongoIds := []string{}
	results, err := db.Query(querystmt)
	CheckError(err)
	results.Scan(mongoIds)

	client, err2 := mongo.Connect(context.TODO(), options.Client().ApplyURI(env.MongoUri))
	CheckError(err2)
	collection := client.Database(env.Mongodb).Collection(env.MongoColl)
	for _, id := range mongoIds {
		deleteMongoDoc(collection, id)
	}
}

func deleteMongoDoc(coll *mongo.Collection, id string) {
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	CheckError(err)
	filter := bson.M{"_id": idPrimitive}
	_, err2 := coll.DeleteOne(context.TODO(), filter, nil)
	CheckError(err2)
}

func rowCount(db *sql.DB, table string) (count int) {
	selectstmt := fmt.Sprintf("SELECT COUNT(*) as count FROM %s", table)
	err := db.QueryRow(selectstmt).Scan(&count)
	CheckError(err)
	return count
}

func removeRow(db *sql.DB, id int) {
	deletestmt := fmt.Sprintf("DELETE FROM %s WHERE id=%d", env.Table, id)
	fmt.Println(deletestmt)
	_, err := db.Exec(deletestmt)
	CheckError(err)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
