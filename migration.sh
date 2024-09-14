goose -dir ./migrations postgres "postgres://postgres:password@localhost:5432/notes?sslmode=disable" status

goose -dir ./migrations postgres "postgres://postgres:password@localhost:5432/notes?sslmode=disable" up