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
const EMAIL_SIZE = 220
const ROW_SIZE = ID_SIZE + USERNAME_SIZE + EMAIL_SIZE //NOTE: Size of row = 256
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
	//NOTE: copy(dst, src)
	copy(buf[ID_SIZE:ID_SIZE+USERNAME_SIZE], r.username[:])
	copy(buf[ID_SIZE+USERNAME_SIZE:], r.email[:])

	return buf
}

func Deserialize(buf []byte) *Row {
	var row Row
	row.id = binary.LittleEndian.Uint32(buf[:ID_SIZE])
	copy(row.username[:], buf[ID_SIZE:ID_SIZE+USERNAME_SIZE])
	copy(row.email[:], buf[ID_SIZE+USERNAME_SIZE:])
	return &row
}

type Statement struct {
	statement_type string
	row            Row
}

type Pager struct {
	pages           [TABLE_MAX_PAGES]*[PAGE_SIZE]byte
	file_descriptor *os.File
}

type Table struct {
	num_rows int
	pager    *Pager
}

type Cursor struct {
	table        *Table
	row_num      int
	end_of_table bool
}

func AllocatePage(pager *Pager, num_rows int32) {
	fileInfo, err := pager.file_descriptor.Stat()
	fileSize := fileInfo.Size()
	totalRowSize := num_rows * ROW_SIZE
	if err != nil {
		fmt.Println(err)
	}
	//NOTE: we need to allocate another page
	if fileSize < int64(totalRowSize) {
		fmt.Println("Allocating page...")
		number_of_pages := (table.num_rows / PAGES_MAX_ROWS)
		var startingOffset int = number_of_pages * int(PAGE_SIZE)
		_, err = pager.file_descriptor.Seek(int64(startingOffset), io.SeekStart)
		buf := make([]byte, PAGE_SIZE)
		_, err = pager.file_descriptor.Write(buf)
	}

}

func WriteRowToFile(pager *Pager, row Row, num_rows int32) error {
	//TODO:
	// Check if we need to write a page
	// Once page is written, find offset based on num_rows
	AllocatePage(pager, num_rows+1)

	var err error
	var offset int32
	offset = ROW_SIZE * num_rows
	_, err = pager.file_descriptor.Seek(int64(offset), io.SeekStart)
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
			fmt.Println(selectedRow)
		}
	}
}

// NOTE: Page starts at 0
func read_page(page_number int) error {
	//TODO:
	// 1. If page 0, copy bytes 0:PAGE_SIZE into a buffer
	// 2. Once we have the buffer, loop over the contents and read data into Pager struct
	file, err := os.Open("mydb.db")
	if err != nil {
		return err
	}
	defer file.Close()

	var starting_offset int = page_number * int(PAGE_SIZE)
	var ending_offset int = int(PAGE_SIZE) + starting_offset
	_, err = file.Seek(int64(starting_offset), io.SeekStart)
	if err != nil {
		return err
	}

	bytesToRead := ending_offset - starting_offset
	fileBuff := make([]byte, bytesToRead)
	_, err = io.ReadFull(file, fileBuff[:bytesToRead])
	if err != nil {
		// Handle error
		return err
	}

	for i := 0; i < int(table.num_rows*ROW_SIZE); i += ROW_SIZE {
		dst := make([]byte, ROW_SIZE)
		copy(dst[:], fileBuff[i:i+ROW_SIZE])
		row := Deserialize(dst)
		if row.id == 0 {
			break
		}
		fmt.Printf("(%d %s %s)\n", row.id, row.username, row.email)
	}

	return nil
}

func execute_statement(statement *Statement, writer io.Writer) error {
	if table.num_rows >= TABLE_MAX_PAGES*PAGES_MAX_ROWS {
		return EXECUTE_TABLE_FULL
	}

	rowToInsert := statement.row
	// pageIndex := table.num_rows / PAGES_MAX_ROWS
	// rowIndex := table.num_rows % PAGES_MAX_ROWS
	switch statement.statement_type {
	case "statement_insert":
		{
			//NOTE: We manage the array ourselves, instead of using append, because arrays are fixed sized and slices are dynamically sized
			//Arrays will be much more efficient
			WriteRowToFile(table.pager, rowToInsert, int32(table.num_rows))
			table.num_rows += 1
		}
	case "statement_select":
		{
			//TODO: 1. Check if value is in cache(which is this for loop)
			//print_cache(pageIndex, writer)
			//TODO: 2. If not, get values from file, print, and store in data structure

			fileInfo, _ := table.pager.file_descriptor.Stat()
			number_of_rows := fileInfo.Size() / ROW_SIZE
			number_of_pages := (number_of_rows / PAGES_MAX_ROWS)
			for i := 0; i < int(number_of_pages); i++ {
				err := read_page(i)
				if err != nil {
					fmt.Println("Error reading page from db")
					fmt.Println(err)
				}
			}
			//print_cache(pageIndex, writer)

		}
	default:
		{
			fmt.Println("Unrecognized statement type")
		}
	}
	return nil
}

func table_start(table *Table) *Cursor {
	cursor := &Cursor{
		table:        table,
		row_num:      0,
		end_of_table: table.num_rows == 0,
	}
	return cursor
}
func table_end(table *Table) *Cursor {
	cursor := &Cursor{
		table:        table,
		row_num:      table.num_rows,
		end_of_table: true,
	}
	return cursor
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
		file_descriptor: file, // Example value
	}
	//NOTE: Each row in the file is of size 256 bytes
	number_of_rows := fileInfo.Size() / ROW_SIZE
	table = &Table{
		num_rows: int(number_of_rows),
		pager:    pager,
	}

	return table, nil

}
