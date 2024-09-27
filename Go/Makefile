run:
	go run main.go goqlite.go b_tree.go mydb.db

debug:
	#go run main.go goqlite.go -d -input="select" mydb.db
	go build -gcflags="all=-N -l" -o goqlite main.go goqlite.go b_tree.go



test:
	@if [ -f mydb.db ]; then rm -r mydb.db; fi
	go test *.go # | grep -v "Goqlite>"

bytes:
	xxd -g 1 mydb.db
