runBlog:
	go run cmd/blog/main.go

runAuth:
	go run cmd/auth/main.go

migrateUp:
	go run ./cmd/migrator/main.go \
   		--storage-path="postgres://root:root_pass@localhost:5432/root_db?sslmode=disable" \
    	--migrations-path="./db/migrations" \
   		--migrations-table="migrations" \
   		--direction="up"

migrateDown:
	go run ./cmd/migrator/main.go \
   		--storage-path="postgres://root:root_pass@localhost:5432/root_db?sslmode=disable" \
    	--migrations-path="./db/migrations" \
   		--migrations-table="migrations" \
   		--direction="down"


generateProto:
	protoc \
 	 -I protos/proto \
 	 protos/proto/auth/auth.proto \
 	 --go_out=protos/gen/go --go_opt=paths=source_relative \
 	 --go-grpc_out=protos/gen/go --go-grpc_opt=paths=source_relative