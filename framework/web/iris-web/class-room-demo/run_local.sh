cd $(dirname $0)

docker-compose -f ./docker-compose.yml up -d

go run main.go
