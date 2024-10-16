package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Common Node Header Layout
const (
	NODE_TYPE_SIZE          uint32 = 1                                                   // 1 byte
	NODE_TYPE_OFFSET        uint32 = 0                                                   // byte 0
	IS_ROOT_SIZE            uint32 = 1                                                   // 1 byte
	IS_ROOT_OFFSET          uint32 = 1                                                   // byte 1
	PARENT_POINTER_SIZE     uint32 = 4                                                   // 4 bytes
	PARENT_POINTER_OFFSET   uint32 = IS_ROOT_OFFSET + IS_ROOT_SIZE                       // byte 2
	COMMON_NODE_HEADER_SIZE uint32 = NODE_TYPE_SIZE + IS_ROOT_SIZE + PARENT_POINTER_SIZE // byte 6
)

// Free space pointer (2 bytes) added at byte 6 in the header
const (
	FREE_SPACE_POINTER_SIZE   uint32 = 2                                                 // 2 bytes
	FREE_SPACE_POINTER_OFFSET uint32 = COMMON_NODE_HEADER_SIZE                           // byte 6
	NEW_COMMON_HEADER_SIZE    uint32 = COMMON_NODE_HEADER_SIZE + FREE_SPACE_POINTER_SIZE // byte 8
)

// Leaf Node Header Layout
const (
	LEAF_NODE_NUM_CELLS_SIZE   uint32 = 4                                                 // 4 bytes
	LEAF_NODE_NUM_CELLS_OFFSET uint32 = NEW_COMMON_HEADER_SIZE                            // byte 8
	LEAF_NODE_HEADER_SIZE      uint32 = NEW_COMMON_HEADER_SIZE + LEAF_NODE_NUM_CELLS_SIZE // byte 12
)

// Leaf Node Body Layout
const (
	LEAF_NODE_KEY_SIZE   uint32 = 4
	LEAF_NODE_KEY_OFFSET uint32 = LEAF_NODE_HEADER_SIZE

	LEAF_NODE_VALUE_SIZE   uint32 = USERNAME_SIZE + EMAIL_SIZE                // USERNAME + EMAIL: 252 bytes
	LEAF_NODE_VALUE_OFFSET uint32 = LEAF_NODE_KEY_OFFSET + LEAF_NODE_KEY_SIZE // byte 4

	LEAF_NODE_CELL_SIZE uint32 = LEAF_NODE_KEY_SIZE + LEAF_NODE_VALUE_SIZE // 256 bytes

)

func (t *Table) _WriteData(row Row) {
	// Read the current free space pointer
	_, _ = t.pager.file_descriptor.Seek(int64(FREE_SPACE_POINTER_OFFSET), io.SeekStart)
	var freeSpacePointer uint16
	_ = binary.Read(t.pager.file_descriptor, binary.LittleEndian, &freeSpacePointer)

	// Write Key (ID)
	_, _ = t.pager.file_descriptor.Seek(int64(freeSpacePointer), io.SeekStart)
	idBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(idBytes, row.id)
	_, _ = t.pager.file_descriptor.Write(idBytes)

	// Update free space pointer after writing key
	freeSpacePointer += uint16(LEAF_NODE_KEY_SIZE)

	// Write Values (username and email)
	buf := make([]byte, USERNAME_SIZE+EMAIL_SIZE)
	copy(buf[:USERNAME_SIZE], row.username[:])
	copy(buf[USERNAME_SIZE:], row.email[:])
	_, _ = t.pager.file_descriptor.Seek(int64(freeSpacePointer), io.SeekStart)
	_, _ = t.pager.file_descriptor.Write(buf)

	// Update free space pointer after writing values
	freeSpacePointer += uint16(USERNAME_SIZE + EMAIL_SIZE)

	// Write back the updated free space pointer
	_, _ = t.pager.file_descriptor.Seek(int64(FREE_SPACE_POINTER_OFFSET), io.SeekStart)
	_ = binary.Write(t.pager.file_descriptor, binary.LittleEndian, freeSpacePointer)

	//fmt.Println("Final freeSpacePointer:", freeSpacePointer)
}

func (t *Table) Insert(row Row) error {
	fileInfo, err := t.pager.file_descriptor.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()
	//NOTE: If this is the first insert, we need to allocate the initial page
	//	Remember, each page is a node in the b tree
	if fileSize == 0 {
		t.pager.AllocateRoot()
	}
	//NOTE: Checks if there is enough space for a row in the page
	if (PAGE_SIZE - FREE_SPACE_POINTER_OFFSET) < ROW_SIZE {
		return errors.New("Not enough space to insert in Page!")
		//TODO: This is where we would do the split
	}
	t._WriteData(row)
	return nil
}
func (t *Table) Select() error {
	return nil
}

func (t *Table) NodeType(page int) (error, byte) {
	return nil, t.pager.tree.root[page][NODE_TYPE_OFFSET]
}
func (t *Table) IsRoot(page int) (error, byte) {
	return nil, t.pager.tree.root[page][IS_ROOT_OFFSET]
}
func (t *Table) ParentPointer(page int) (error, []byte) {
	return nil, t.pager.tree.root[page][PARENT_POINTER_OFFSET : PARENT_POINTER_OFFSET+PARENT_POINTER_SIZE]
}
func (t *Table) NumCells(page int) (error, []byte) {
	return nil, t.pager.tree.root[page][LEAF_NODE_NUM_CELLS_OFFSET : LEAF_NODE_NUM_CELLS_OFFSET+LEAF_NODE_NUM_CELLS_SIZE]
}
