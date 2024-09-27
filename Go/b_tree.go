package main

import ()

// Common Node Header Layout
const (
	NODE_TYPE_SIZE          uint32 = 1                                                   // 1
	NODE_TYPE_OFFSET        uint32 = 0                                                   // 0
	IS_ROOT_SIZE            uint32 = 1                                                   // 1
	IS_ROOT_OFFSET          uint32 = NODE_TYPE_SIZE                                      // 1
	PARENT_POINTER_SIZE     uint32 = 4                                                   // 4
	PARENT_POINTER_OFFSET   uint32 = IS_ROOT_OFFSET + IS_ROOT_SIZE                       // 2
	COMMON_NODE_HEADER_SIZE uint32 = NODE_TYPE_SIZE + IS_ROOT_SIZE + PARENT_POINTER_SIZE // 6
)

// Leaf Node Header Layout
const (
	LEAF_NODE_NUM_CELLS_SIZE   uint32 = 4                                                  // 4
	LEAF_NODE_NUM_CELLS_OFFSET uint32 = COMMON_NODE_HEADER_SIZE                            // 6
	LEAF_NODE_HEADER_SIZE      uint32 = COMMON_NODE_HEADER_SIZE + LEAF_NODE_NUM_CELLS_SIZE // 10
)

// Leaf Node Body Layout
const (
	LEAF_NODE_KEY_SIZE        uint32 = 4                                               // 4
	LEAF_NODE_KEY_OFFSET      uint32 = 0                                               // 0
	LEAF_NODE_VALUE_SIZE      uint32 = ROW_SIZE                                        // 291 (assuming ROW_SIZE = 291)
	LEAF_NODE_VALUE_OFFSET    uint32 = LEAF_NODE_KEY_OFFSET + LEAF_NODE_KEY_SIZE       // 4
	LEAF_NODE_CELL_SIZE       uint32 = LEAF_NODE_KEY_SIZE + LEAF_NODE_VALUE_SIZE       // 295
	LEAF_NODE_SPACE_FOR_CELLS uint32 = PAGE_SIZE - LEAF_NODE_HEADER_SIZE               // 4086
	LEAF_NODE_MAX_CELLS       uint32 = LEAF_NODE_SPACE_FOR_CELLS / LEAF_NODE_CELL_SIZE // 13
)

func (t *Pager) Insert(row Row) error {
	return nil
}
func (t *BPlusTree) Select() error {
	return nil
}

func (t *BPlusTree) NodeType(page int) (error, byte) {
	return nil, t.root[page][NODE_TYPE_OFFSET]
}
func (t *BPlusTree) IsRoot(page int) (error, byte) {
	return nil, t.root[page][IS_ROOT_OFFSET]
}
func (t *BPlusTree) ParentPointer(page int) (error, []byte) {
	return nil, t.root[page][PARENT_POINTER_OFFSET : PARENT_POINTER_OFFSET+PARENT_POINTER_SIZE]
}
func (t *BPlusTree) NumCells(page int) (error, []byte) {
	return nil, t.root[page][LEAF_NODE_NUM_CELLS_OFFSET : LEAF_NODE_NUM_CELLS_OFFSET+LEAF_NODE_NUM_CELLS_SIZE]
}

func insert(table *Table, row Row) error {
	table.pager.Insert(row)
	return nil
}
