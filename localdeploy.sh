docker build -t requestbucket .
docker build -t reqcron -f cron.Dockerfile .
docker compose up -d
echo 'RequestBucket running on port 3000'