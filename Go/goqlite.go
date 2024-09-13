package main

import (
	"encoding/binary"
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
const ID_SIZE = 4
const USERNAME_SIZE = 32
const EMAIL_SIZE = 255
const ROW_SIZE = ID_SIZE + USERNAME_SIZE + EMAIL_SIZE
const PAGE_SIZE int32 = 4096

var table *Table

type Row struct {
	id       uint32
	username [USERNAME_SIZE]byte
	email    [EMAIL_SIZE]byte
}

func (r *Row) Serialize() []byte {
	buf := make([]byte, ROW_SIZE)
	binary.LittleEndian.PutUint32(buf[:ID_SIZE], r.id)
	copy(buf[ID_SIZE:ID_SIZE+USERNAME_SIZE], r.username[:])
	copy(buf[ID_SIZE+USERNAME_SIZE:], r.email[:])

	fmt.Println("Length of serialized bytes:", len(buf)) // Print the length of bytes
	return buf
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

func WriteRowToFile(pager *Pager, row Row, num_rows int32) error {
	var err error
	var offset int32
	if num_rows == 0 {
		offset = 0
	} else {
		offset = ROW_SIZE * (num_rows + 1)
	}
	_, err = pager.file_descriptor.Seek(int64(offset), 0)
	if err != nil {
		return err
	}

	_, err = pager.file_descriptor.Write(row.Serialize())
	return err
}

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

func print_cache(pageIndex int, writer io.Writer) {
	for i := 0; i <= pageIndex; i++ {
		startIndex := i * PAGES_MAX_ROWS
		endIndex := min((i+1)*PAGES_MAX_ROWS, table.num_rows)
		for j := startIndex; j < endIndex; j++ {
			selectedRow := (*table.pager.pages[i])[j%PAGES_MAX_ROWS] // Adjust index for page
			fmt.Fprintf(writer, "(%d %s %s)\n", selectedRow.id, selectedRow.username, selectedRow.email)
		}
	}
}

// NOTE: Page starts at 1
func read_page(page_number int) error {
	fmt.Println("IN READ_PAGES")
	//TODO:
	// 1. If page 1, copy bytes 0:PAGE_SIZE into a buffer
	// 2. Once we have the buffer, loop over the contents and read data into Pager struct
	file, err := os.Open("mydb.db")
	if err != nil {
		return err
	}
	defer file.Close()

	var starting_offset int = 0
	var ending_offset int = table.num_rows * ROW_SIZE
	_, err = file.Seek(int64(starting_offset), io.SeekStart)
	if err != nil {
		return err
	}

	bytesToRead := ending_offset - starting_offset
	buf := make([]byte, bytesToRead)
	_, err = io.ReadFull(file, buf[:bytesToRead])
	if err != nil {
		// Handle error
		return err
	}
	fmt.Println(buf)

	return nil
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
			WriteRowToFile(table.pager, rowToInsert, int32(table.num_rows))
			table.num_rows += 1
		}
	case "statement_select":
		{
			//TODO: 1. Check if value is in cache(which is this for loop)
			//print_cache(pageIndex, writer)
			//TODO: 2. If not, get values from file, print, and store in data structure
			number_of_pages := (table.num_rows / PAGES_MAX_ROWS) + 1
			for i := 1; i <= number_of_pages; i++ {
				err := read_page(i)
				if err != nil {
					fmt.Println("Error reading page from db")
					fmt.Println(err)
				}
			}
			print_cache(pageIndex, writer)

		}
	default:
		{
			fmt.Println("Unrecognized statement type")
		}
	}
	return nil
}

func db_open(filename string) (*Table, error) {
	table, err := table_open(filename)
	if err != nil {
		return nil, err
	}
	return table, nil
}

func db_close(table *Table) {
	table.pager.file_descriptor.Close()
}

func table_open(filename string) (*Table, error) {
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
	//NOTE: Each row in the file is of size 291 bytes
	number_of_rows := fileInfo.Size() / ROW_SIZE
	table = &Table{
		num_rows: int(number_of_rows),
		pager:    pager,
	}

	return table, nil

}
