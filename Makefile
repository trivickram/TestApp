DB_DSN ?= root:Password@123@tcp(localhost:3306)/hospital?parseTime=true

.PHONY: server gateway frontend proto seed init kill all

server:
	DB_DSN="$(DB_DSN)" go run ./server/

gateway:
	go run ./gateway/

frontend:
	cd frontend && npm run dev

proto:
	protoc \
		--go_out=generated --go_opt=paths=import \
		--go-grpc_out=generated --go-grpc_opt=paths=import \
		proto/hospital.proto
	cp generated/hospital/generated/proto/hospital.pb.go    generated/proto/hospital.pb.go
	cp generated/hospital/generated/proto/hospital_grpc.pb.go generated/proto/hospital_grpc.pb.go
	rm -rf generated/hospital

seed:
	mysql -u root -pPassword@123 hospital < scripts/seed.sql

init:
	mysql -u root -pPassword@123 < scripts/init.sql

kill:
	-lsof -ti:50051 -ti:8080 | xargs kill -9 2>/dev/null

all: kill
	DB_DSN="$(DB_DSN)" go run ./server/ &
	go run ./gateway/ &
	cd frontend && npm run dev
