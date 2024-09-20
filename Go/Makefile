run:
	go run main.go goqlite.go mydb.db

test:
	@if [ -f mydb.db ]; then rm -r mydb.db; fi
	go test *.go # | grep -v "Goqlite>"

bytes:
	xxd -g 1 mydb.db
