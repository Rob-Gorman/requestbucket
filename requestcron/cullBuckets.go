package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reqcron/environment"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
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

var env *environment.Env

func main() {
	env = environment.LoadDotenv()
	log := retrieveLog()
	defer log.Close()

	// psqlconn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", "localhost", env.Pgport, "postgres", env.Pgdbname)
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", env.Pghost, env.Pgport, env.Pguser, env.Pgdbname)
	pgdb, err2 := sql.Open("postgres", psqlconn)
	CheckError(err2)
	defer pgdb.Close()

	cullPGBuckets(pgdb, log)
}

func retrieveLog() *os.File {
	fmt.Println(env.Logfile)
	log, err1 := os.OpenFile(env.Logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	if err1 != nil {
		if os.IsNotExist(err1) {
			log, _ = os.Create(env.Logfile)
		} else {
			errors.Errorf("we have a problem %v", err1)
			return nil
		}
	}
	return log
}

func writeLog(log *os.File, value string) {
	a, b, c := time.Now().Clock()
	log.WriteString(fmt.Sprintf("%d:%d:%d DELETED:%v\n", a, b, c, value))
}

// iterates through oldest buckets and deletes any older than BUCKETS_TTL env var
func cullPGBuckets(db *sql.DB, log *os.File) {
	for rowCount(db, env.Table) > 0 {
		oldestBucket := getOldestBucket(db)

		TTL, err := strconv.ParseFloat(os.Getenv("BUCKET_TTL"), 64)
		CheckError(err)

		if time.Since(*oldestBucket.created_at).Hours() > TTL {
			deleteBucket(db, oldestBucket)
			writeLog(log, *oldestBucket.url)
		}
	}
}

func getOldestBucket(db *sql.DB) *bucket {
	selectstmt := fmt.Sprintf("SELECT * FROM %s order by id limit 1", env.Table)
	oldestRow := bucket{}
	err := db.QueryRow(selectstmt).Scan(&oldestRow.id, &oldestRow.url, &oldestRow.created_at)
	CheckError(err)
	return &oldestRow
}

func deleteBucket(db *sql.DB, oldestBucket *bucket) {
	bucketId := *oldestBucket.id
	cullPGRequests(db, bucketId)
	removeRow(db, bucketId)
}

func cullPGRequests(db *sql.DB, bucketId int) {
	removeMongoIds(db, bucketId)
	deletestmt := fmt.Sprintf("DELETE FROM %s WHERE bucket_id=%d", "requests", bucketId)
	fmt.Println(deletestmt)
	_, err := db.Exec(deletestmt)
	CheckError(err)
}

func removeMongoIds(db *sql.DB, bucketId int) {
	mongoIds := getMongoIds(db, bucketId)

	collection := getMongoCollection()

	for _, id := range mongoIds {
		deleteMongoDoc(collection, id)
	}
}

func getMongoCollection() *mongo.Collection {
	client, err2 := mongo.Connect(context.Background(), options.Client().ApplyURI(env.MongoUri))
	CheckError(err2)

	collection := client.Database(env.Mongodb).Collection(env.MongoColl)
	return collection
}

func getMongoIds(db *sql.DB, bucketId int) (mongoIds []string) {
	querystmt := fmt.Sprintf("SELECT mongo_document_ref FROM requests WHERE bucket_id=%d", bucketId)
	results, err := db.Query(querystmt)
	CheckError(err)
	mongoIds = marshalMongoIds(results)
	return
}

func marshalMongoIds(results *sql.Rows) []string {
	ids := []string{}
	for results.Next() {
		var rowId string
		results.Scan(&rowId)
		ids = append(ids, rowId)
	}
	return ids
}

func deleteMongoDoc(coll *mongo.Collection, id string) {
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	CheckError(err)
	filter := bson.M{"_id": idPrimitive}
	_, err = coll.DeleteOne(context.Background(), filter, nil)
	CheckError(err)
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
		errors.Errorf("\nproblem here! %v", err)
		panic(err) // not great form, but it's low stakes
	}
}
