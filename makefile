runBlog:
	go run cmd/blog/main.go

runAuth:
	go run cmd/auth/main.go

runNotification:
	go run cmd/notification/main.go

#migrateUp:
#	go run ./cmd/migrator/main.go \
#   		--storage-path="postgres://root:root_pass@localhost:5432/root_db?sslmode=disable" \
#    	--migrations-path="./dbcodes/migrations" \
#   		--migrations-table="migrations" \
#   		--direction="up"
#
#migrateDown:
#	go run ./cmd/migrator/main.go \
#   		--storage-path="postgres://root:root_pass@localhost:5432/root_db?sslmode=disable" \
#    	--migrations-path="./dbcodes/migrations" \
#   		--migrations-table="migrations" \
#   		--direction="down"


authGenerateProto:
	protoc \
 	 -I protos/proto \
 	 protos/proto/auth/auth.proto \
 	 --go_out=protos/gen/go --go_opt=paths=source_relative \
 	 --go-grpc_out=protos/gen/go --go-grpc_opt=paths=source_relative

notificationGenerateProto:
	protoc \
 	 -I protos/proto \
 	 protos/proto/notification/notification.proto \
 	 --go_out=protos/gen/go --go_opt=paths=source_relative \
 	 --go-grpc_out=protos/gen/go --go-grpc_opt=paths=source_relative

authMigrateUp:
	migrate -source file://db/auth/migrations \
	-database "postgres://root:root_pass@localhost:5432/auth_db?sslmode=disable" \
	up


authMigrateDown:
	migrate -source file://db/auth/migrations \
	-database "postgres://root:root_pass@localhost:5432/auth_db?sslmode=disable" \
	down


blogMigrateUp:
	migrate -source file://db/blog/migrations \
	-database "postgres://root:root_pass@localhost:5432/blog_db?sslmode=disable" \
	up

blogMigrateDown:
	migrate -source file://db/blog/migrations \
	-database "postgres://root:root_pass@localhost:5432/blog_db?sslmode=disable" \
	down
