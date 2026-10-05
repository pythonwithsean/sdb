package storageengine

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	queryprocessor "sdb/QueryProcessor"
)

var currentDatabase string

func typeToByte(t string) (uint8, error) {
	switch t {
	case "INT":
		return uint8(1), nil
	case "STRING":
		return uint8(2), nil
	default:
		return 0, fmt.Errorf("unknown type: %s", t)
	}
}

// Permission bits 4 for read, 2 for write, 1 for execute
// we have owner, group and others
func executeCreateDatabase(stmt *queryprocessor.CreateDatabase) error {
	fd, err := os.OpenFile("."+stmt.Name+".db", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	defer fd.Close()
	if len(stmt.Name) > math.MaxUint8 {
		return fmt.Errorf("database name is %d bytes; maximum is %d", len(stmt.Name), math.MaxUint8)
	}

	err = binary.Write(fd, binary.LittleEndian, uint8(len(stmt.Name)))
	if err != nil {
		return fmt.Errorf("failed to write database name length: %w", err)
	}

	err = binary.Write(fd, binary.LittleEndian, []byte(stmt.Name))
	if err != nil {
		return fmt.Errorf("failed to write database name: %w", err)
	}
	currentDatabase = stmt.Name
	return nil
}

func executeCreateTable(stmt *queryprocessor.CreateTable) error {
	if currentDatabase == "" {
		return fmt.Errorf("no database selected")
	}

	if len(stmt.Name) > math.MaxUint8 {
		return fmt.Errorf("table name is %d bytes; maximum is %d", len(stmt.Name), math.MaxUint8)
	}
	if len(stmt.Columns) > math.MaxUint8 {
		return fmt.Errorf("column count is %d; maximum is %d", len(stmt.Columns), math.MaxUint8)
	}
	fd, err := os.OpenFile("."+currentDatabase+".db", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open database file: %w", err)
	}
	err = binary.Write(fd, binary.LittleEndian, uint8(len(stmt.Name)))
	if err != nil {
		return fmt.Errorf("failed to write table name length: %w", err)
	}
	err = binary.Write(fd, binary.LittleEndian, []byte(stmt.Name))
	if err != nil {
		return fmt.Errorf("failed to write table name: %w", err)
	}
	err = binary.Write(fd, binary.LittleEndian, uint32(0)) // number of rows, will be updated later
	if err != nil {
		return fmt.Errorf("failed to write number of rows: %w", err)
	}
	err = binary.Write(fd, binary.LittleEndian, uint8(len(stmt.Columns)))
	if err != nil {
		return fmt.Errorf("failed to write number of columns: %w", err)
	}

	max_row_length := 0
	for _, col := range stmt.Columns {
		ty, err := typeToByte(col.Type)
		if err != nil {
			return fmt.Errorf("failed to convert type: %w", err)
		}
		if ty == 1 { // INT
			max_row_length += 4 // 4 bytes for INT
		}
		if ty == 2 { // STRING
			max_row_length += 255 // 255 bytes for STRING
		}
	}

	if max_row_length > math.MaxUint16 {
		return fmt.Errorf("maximum row size is %d bytes; maximum is %d", max_row_length, math.MaxUint16)
	}
	err = binary.Write(fd, binary.LittleEndian, uint16(max_row_length))
	if err != nil {
		return fmt.Errorf("failed to write max row length: %w", err)
	}
	for _, col := range stmt.Columns {
		if len(col.Name) > math.MaxUint8 {
			return fmt.Errorf("column name is %d bytes; maximum is %d", len(col.Name), math.MaxUint8)
		}
		err = binary.Write(fd, binary.LittleEndian, uint8(len(col.Name)))
		if err != nil {
			return fmt.Errorf("failed to write column name length: %w", err)
		}

		err = binary.Write(fd, binary.LittleEndian, []byte(col.Name))
		if err != nil {
			return fmt.Errorf("failed to write column name: %w", err)
		}

		ty, err := typeToByte(col.Type)
		if err != nil {
			return fmt.Errorf("failed to convert column type: %w", err)
		}

		err = binary.Write(fd, binary.LittleEndian, ty)
		if err != nil {
			return fmt.Errorf("failed to write column type: %w", err)
		}
	}

	return nil
}

func executeSelect(stmt *queryprocessor.Select) error {
	return nil
}

func executeInsert(stmt *queryprocessor.Insert) error {
	return nil
}

func executeUpdate(stmt *queryprocessor.Update) error {
	return nil
}

func executeDelete(stmt *queryprocessor.Delete) error {
	return nil
}

func ExecuteStatement(stmt queryprocessor.Statement) error {
	fmt.Printf("Executing statement: %+v\n", stmt)
	switch s := stmt.(type) {
	case *queryprocessor.CreateDatabase:
		return executeCreateDatabase(s)
	case *queryprocessor.CreateTable:
		return executeCreateTable(s)
	case *queryprocessor.Select:
		return executeSelect(s)
	case *queryprocessor.Insert:
		return executeInsert(s)
	case *queryprocessor.Update:
		return executeUpdate(s)
	case *queryprocessor.Delete:
		return executeDelete(s)
	default:
		return fmt.Errorf("unknown statement type: %T", s)
	}
}
