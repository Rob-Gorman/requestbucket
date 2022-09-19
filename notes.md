### Postgres

`docker run --name rbpg -d -p5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=requestbucket -v $PWD/requestbucket/schema.sql:/docker-entrypoint-initdb.d/schema.sql postgres`


### MONGO

`docker run --name rbmdb -d -p27017:27017 mongo`
