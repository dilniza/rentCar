
migration-up:
	migrate -path ./migrations/postgres -database 'postgres://admin:admin@localhost:5432/rentcar?sslmode=disable' up
	
migration-down:
	migrate -path ./migrations/postgres -database 'postgres://admin:admin@localhost:5432/rentcar?sslmode=disable' down
	
migration-force-1v:
	migrate -path ./migrations/postgres -database 'postgres://admin:admin@localhost:5432/rentcar?sslmode=disable' force 1
	
	
	migrate -path ./migrations/postgres -database 'postgres://admin:admin@localhost:5432/rentcar?sslmode=disable' up 1 01_add_column.up

swag-init:
	swag init -g api/router.go -o api/docs
