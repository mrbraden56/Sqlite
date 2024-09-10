package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// var PREPARE_UNRECOGNIZED_STATEMENT = errors.New("Prepare Unrecognized Statement")
// var PREPARE_PARSING_ERROR = errors.New("Prepare Parsing Error")
// var EXECUTE_TABLE_FULL = errors.New("Execute Table Full")

const TABLE_MAX_PAGES = 100
const PAGES_MAX_ROWS = 14

var table Table

type Row struct {
	id       int
	username string
	email    string
}

type Statement struct {
	statement_type string
	row            Row
}

type Table struct {
	num_rows int
	pages    [TABLE_MAX_PAGES]*[PAGES_MAX_ROWS]Row
}

// NOTE: This is our Sql Compiler
func prepare_statement(input string, statement *Statement) error {
	fields := strings.Fields(input)
	var err error
	if fields[0] == "select" {
		//NOTE: Example: select 1 cstack foo@bar.com
		statement.statement_type = "statement_select"
		statement.row.id, err = strconv.Atoi(fields[1])
		if err != nil {
			fmt.Println("Can't convert this to an int!")
			return PREPARE_PARSING_ERROR
		}
		statement.row.username = fields[2]
		statement.row.email = fields[3]
		return nil
	}
	if fields[0] == "insert" {
		//NOTE: Example: insert 1 cstack foo@bar.com
		statement.statement_type = "statement_insert"
		if fields[1] == "-1" {
			return PREPARE_NEGATIVE_ID
		}
		statement.row.id, err = strconv.Atoi(fields[1])
		if err != nil {
			fmt.Println("Can't convert this to an int!")
			return PREPARE_PARSING_ERROR
		}
		statement.row.username = fields[2]
		statement.row.email = fields[3]
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
	if table.pages[pageIndex] == nil {
		table.pages[pageIndex] = new([PAGES_MAX_ROWS]Row)
	}
	switch statement.statement_type {
	case "statement_insert":
		{
			//NOTE: We manage the array ourselves, instead of using append, because arrays are fixed sized and slices are dynamically sized
			//Arrays will be much more efficient
			(*table.pages[pageIndex])[rowIndex] = rowToInsert
			table.num_rows += 1
		}
	case "statement_select":
		{
			rowId := statement.row.id
			var selectedRow Row
		OuterLoop:
			for i := 0; i < TABLE_MAX_PAGES; i++ {
				for j := 0; j < PAGES_MAX_ROWS; j++ {
					currentRowId := (*table.pages[i])[j].id
					if rowId == currentRowId {
						selectedRow = (*table.pages[i])[j]
						break OuterLoop
					}
				}
			}
			fmt.Fprintf(writer, "(%d %s %s)\n", selectedRow.id, selectedRow.username, selectedRow.email)
		}
	default:
		{
			fmt.Println("Unrecognized statement type")
		}
	}
	return nil
}

func NewTable() *Table {
	return &Table{}
}
