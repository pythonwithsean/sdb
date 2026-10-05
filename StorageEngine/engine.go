package storageengine

import (
	"fmt"
	"os"
	queryprocessor "sdb/QueryProcessor"
)

// Permission bits 4 for read, 2 for write, 1 for execute
// we have owner, group and others
func executeCreateDatabase(stmt *queryprocessor.CreateDatabase) error {
	_, err := os.OpenFile("."+stmt.Name+".db", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	return nil
}

func executeCreateTable(stmt *queryprocessor.CreateTable) error {
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
