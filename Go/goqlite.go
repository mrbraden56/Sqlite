package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// var PREPARE_UNRECOGNIZED_STATEMENT = errors.New("Prepare Unrecognized Statement")
// var PREPARE_PARSING_ERROR = errors.New("Prepare Parsing Error")
// var EXECUTE_TABLE_FULL = errors.New("Execute Table Full")

const TABLE_MAX_PAGES = 100
const PAGES_MAX_ROWS = 14
const USERNAME_SIZE = 32
const EMAIL_SIZE = 32

var table *Table

type Row struct {
	id       uint32
	username [USERNAME_SIZE]byte
	email    [EMAIL_SIZE]byte
}

type Statement struct {
	statement_type string
	row            Row
}

type Pager struct {
	pages           [TABLE_MAX_PAGES]*[PAGES_MAX_ROWS]Row
	file_descriptor *os.File
	file_length     os.FileInfo
}

type Table struct {
	num_rows int
	pager    *Pager
}

// NOTE: This is our Sql Compiler
func prepare_statement(input string, statement *Statement) error {
	fields := strings.Fields(input)
	var err error
	if fields[0] == "select" {
		//NOTE: Example: select 1 cstack foo@bar.com
		statement.statement_type = "statement_select"
		return nil
	}
	if fields[0] == "insert" {
		//NOTE: Example: insert 1 cstack foo@bar.com
		statement.statement_type = "statement_insert"
		if fields[1] == "-1" {
			return PREPARE_NEGATIVE_ID
		}
		var int_id int
		int_id, err = strconv.Atoi(fields[1])
		if err != nil {
			fmt.Println("Can't convert this to an int!")
			return PREPARE_PARSING_ERROR
		}
		statement.row.id = uint32(int_id)
		if err != nil {
			fmt.Println("Can't convert this to an int!")
			return PREPARE_PARSING_ERROR
		}
		copy(statement.row.username[:], fields[2])
		copy(statement.row.email[:], fields[3])
		return nil
	}

	return PREPARE_UNRECOGNIZED_STATEMENT
}

func execute_statement(statement *Statement, writer io.Writer) error {
	if table.num_rows >= TABLE_MAX_PAGES*PAGES_MAX_ROWS {
		return EXECUTE_TABLE_FULL
	}

	rowToInsert := statement.row
	pageIndex := table.num_rows / PAGES_MAX_ROWS
	rowIndex := table.num_rows % PAGES_MAX_ROWS
	if table.pager.pages[pageIndex] == nil {
		table.pager.pages[pageIndex] = new([PAGES_MAX_ROWS]Row)
	}
	switch statement.statement_type {
	case "statement_insert":
		{
			//NOTE: We manage the array ourselves, instead of using append, because arrays are fixed sized and slices are dynamically sized
			//Arrays will be much more efficient
			(*table.pager.pages[pageIndex])[rowIndex] = rowToInsert
			table.num_rows += 1
		}
	case "statement_select":
		{
			for i := 0; i <= pageIndex; i++ {
				startIndex := i * PAGES_MAX_ROWS
				endIndex := min((i+1)*PAGES_MAX_ROWS, table.num_rows)
				for j := startIndex; j < endIndex; j++ {
					selectedRow := (*table.pager.pages[i])[j%PAGES_MAX_ROWS] // Adjust index for page
					fmt.Fprintf(writer, "(%d %s %s)\n", selectedRow.id, selectedRow.username, selectedRow.email)
				}
			}
		}
	default:
		{
			fmt.Println("Unrecognized statement type")
		}
	}
	return nil
}

func db_open(filename string) (*Table, error) {
	pager, err := pager_open(filename)
	num_rows := pager.file_length
	fmt.Println(num_rows)
	if err != nil {
		return nil, err
	}
	table = &Table{
		num_rows: 0,
		pager:    pager,
	}
	return table, nil
}

func db_close(table *Table) {
	table.pager.file_descriptor.Close()
}

func pager_open(filename string) (*Pager, error) {
	file, err := os.OpenFile("mydb.db", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, PAGER_OPENING_ERROR
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return nil, PAGER_OPENING_ERROR
	}
	pager := &Pager{
		pages:           [TABLE_MAX_PAGES]*[PAGES_MAX_ROWS]Row{}, // Initialize array of pointers
		file_descriptor: file,                                    // Example value
		file_length:     fileInfo,                                // Example value
	}
	return pager, nil

}
