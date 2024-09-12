run:
	go run main.go goqlite.go mydb.db

test:
	go test *.go | grep -v "Goqlite>"

bytes:
	xxd -g 1 mydb.db
