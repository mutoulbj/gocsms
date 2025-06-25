package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run cmd/gen_migration/main.go <migration_name>")
		os.Exit(1)
	}

	name := os.Args[1]
	timestamp := time.Now().Format("20060102150405")
	prefix := fmt.Sprintf("%s_%s", timestamp, name)
	dir := "migrations"

	upPath := filepath.Join(dir, fmt.Sprintf("%s.up.sql", prefix))
	downPath := filepath.Join(dir, fmt.Sprintf("%s.down.sql", prefix))

	// 保证 migrations 目录存在
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	write := func(path string) {
		if err := os.WriteFile(path, []byte("-- "+strings.ToUpper(filepath.Ext(path)[1:])+" migration\n"), 0644); err != nil {
			panic(err)
		}
		fmt.Println("Created:", path)
	}

	write(upPath)
	write(downPath)
}
