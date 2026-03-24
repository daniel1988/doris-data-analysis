package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/dmp_admin_v2/backend/internal/core"
)

func main() {
	projectAlias := "cjxcx3"
	tables := []string{
		"project_data",
		"project_event",
		"project_event_property",
		"user_properties",
	}

	// 1. Initialize DB connection
	fmt.Println(">>> 正在初始化数据库连接...")
	if err := core.Init(""); err != nil {
		log.Fatalf("!!! 数据库连接失败: %v", err)
	}

	db := core.GetProjectCenter()
	if db == nil {
		log.Fatal("!!! 获取数据库连接失败")
	}

	// 2. Prepare output directory
	outDir := "../docs/sql"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatalf("!!! 创建输出目录失败: %v", err)
	}

	outFile := filepath.Join(outDir, fmt.Sprintf("%s_backup.sql", projectAlias))
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("!!! 创建输出文件失败: %v", err)
	}
	defer f.Close()

	// Write header
	f.WriteString(fmt.Sprintf("-- Backup for project: %s\n", projectAlias))
	f.WriteString("-- Tables: " + strings.Join(tables, ", ") + "\n\n")

	// 3. Export data for each table
	for _, tableName := range tables {
		fmt.Printf(">>> 正在导出表 %s...\n", tableName)

		query := fmt.Sprintf("SELECT * FROM dmp_center.%s WHERE project_alias = ?", tableName)
		rows, err := db.Raw(query, projectAlias).Rows()
		if err != nil {
			log.Printf("!!! 查询表 %s 失败: %v\n", tableName, err)
			continue
		}

		cols, err := rows.Columns()
		if err != nil {
			rows.Close()
			log.Printf("!!! 获取表 %s 列信息失败: %v\n", tableName, err)
			continue
		}

		f.WriteString(fmt.Sprintf("-- Data for table %s\n", tableName))

		rowCount := 0
		for rows.Next() {
			rowCount++
			values := make([]interface{}, len(cols))
			valuePtrs := make([]interface{}, len(cols))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				log.Printf("!!! 读取行数据失败: %v\n", err)
				continue
			}

			// Format values for INSERT
			var valStrs []string
			for _, v := range values {
				if v == nil {
					valStrs = append(valStrs, "NULL")
				} else {
					switch val := v.(type) {
					case []byte:
						valStrs = append(valStrs, fmt.Sprintf("'%s'", escapeSQL(string(val))))
					case string:
						valStrs = append(valStrs, fmt.Sprintf("'%s'", escapeSQL(val)))
					default:
						valStrs = append(valStrs, fmt.Sprintf("'%v'", val))
					}
				}
			}

			insertStmt := fmt.Sprintf("INSERT INTO dmp_center.%s (%s) VALUES (%s);\n",
				tableName,
				strings.Join(cols, ", "),
				strings.Join(valStrs, ", "))
			f.WriteString(insertStmt)
		}
		rows.Close()

		f.WriteString("\n")
		fmt.Printf(">>> 表 %s 导出完成，共 %d 行。\n", tableName, rowCount)
	}

	fmt.Printf(">>> 备份完成！文件保存在: %s\n", outFile)
}

func escapeSQL(s string) string {
	s = strings.ReplaceAll(s, "'", "''")
	s = strings.ReplaceAll(s, "\\", "\\\\")
	return s
}
