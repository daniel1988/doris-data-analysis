package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"gitee.com/dmp_admin_v2/backend/internal/core"
)

// getSQLDir 获取 SQL 文件所在目录
// 优先使用命令行参数 -sql-dir，否则根据源代码位置推断 docs/sql 目录
func getSQLDir(flagDir string) string {
	if flagDir != "" {
		return flagDir
	}
	// 通过 runtime.Caller 获取当前文件所在目录，向上回退到项目根目录
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("!!! 无法获取当前文件路径")
	}
	// main.go 位于 backend/cmd/db_init/main.go，向上3层到项目根目录
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))
	return filepath.Join(projectRoot, "docs", "sql")
}

// readSQLFile 读取指定 SQL 文件内容
func readSQLFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("!!! 读取 SQL 文件失败 [%s]: %v", path, err)
	}
	return string(data)
}

// parseSQLStatements 将 SQL 文件内容按分号拆分为独立语句
// 过滤掉注释、空行、DROP TABLE 语句
func parseSQLStatements(content string) []string {
	// 按分号拆分
	rawParts := strings.Split(content, ";")
	var stmts []string

	for _, part := range rawParts {
		// 逐行过滤注释
		lines := strings.Split(part, "\n")
		var cleanLines []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "--") {
				continue
			}
			cleanLines = append(cleanLines, line)
		}
		stmt := strings.TrimSpace(strings.Join(cleanLines, "\n"))
		if stmt == "" {
			continue
		}

		// 跳过 DROP TABLE 语句
		upper := strings.ToUpper(stmt)
		if strings.HasPrefix(upper, "DROP") {
			fmt.Printf("[跳过] DROP 语句: %s\n", truncate(stmt, 60))
			continue
		}

		stmts = append(stmts, stmt)
	}
	return stmts
}

// truncate 截断字符串用于日志输出
func truncate(s string, maxLen int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

// reTableName 匹配 CREATE TABLE 语句中的表名（带或不带反引号）
var reTableName = regexp.MustCompile("(?i)CREATE\\s+TABLE\\s+(?:IF\\s+NOT\\s+EXISTS\\s+)?(`?)(\\w+)(`?)")

// addDatabasePrefix 为 CREATE TABLE 语句添加数据库前缀，并确保 IF NOT EXISTS
func addDatabasePrefix(stmt string, dbName string) string {
	// 替换 CREATE TABLE `table_name` 为 CREATE TABLE IF NOT EXISTS `dbName`.`table_name`
	result := reTableName.ReplaceAllStringFunc(stmt, func(match string) string {
		sub := reTableName.FindStringSubmatch(match)
		if len(sub) < 4 {
			return match
		}
		tableName := sub[2]
		return fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`.`%s`", dbName, tableName)
	})
	return result
}

func main() {
	sqlDir := flag.String("sql-dir", "", "SQL文件所在目录路径 (默认自动推断)")
	flag.Parse()

	// 1. 初始化核心配置（数据库连接）
	fmt.Println(">>> 正在初始化数据库连接...")
	if err := core.Init(""); err != nil {
		log.Fatalf("!!! 数据库连接失败: %v", err)
	}

	db := core.GetProjectCenter()
	if db == nil {
		log.Fatal("!!! 获取数据库连接失败")
	}

	dir := getSQLDir(*sqlDir)
	fmt.Printf(">>> SQL 文件目录: %s\n", dir)

	// ============================================================
	// 第一阶段：基于 dmp_center.sql 创建中心库表
	// ============================================================
	fmt.Println("\n========================================")
	fmt.Println(">>> 第一阶段：初始化 dmp_center 数据库")
	fmt.Println("========================================")

	centerSQLPath := filepath.Join(dir, "dmp_center.sql")
	centerSQL := readSQLFile(centerSQLPath)
	centerStmts := parseSQLStatements(centerSQL)

	centerSuccess := 0
	centerSkip := 0
	centerFail := 0

	for _, stmt := range centerStmts {
		upper := strings.ToUpper(strings.TrimSpace(stmt))

		// 对 USE 语句直接跳过（我们已连接到正确的数据库）
		if strings.HasPrefix(upper, "USE ") {
			fmt.Printf("[跳过] USE 语句\n")
			centerSkip++
			continue
		}

		// 对 CREATE TABLE 语句添加 IF NOT EXISTS（如果没有的话）
		if strings.Contains(upper, "CREATE TABLE") && !strings.Contains(upper, "IF NOT EXISTS") {
			// 在 CREATE TABLE 后插入 IF NOT EXISTS
			stmt = strings.Replace(stmt, "CREATE TABLE", "CREATE TABLE IF NOT EXISTS", 1)
		}

		// 对 CREATE TABLE 中没有数据库前缀的表添加 dmp_center 前缀
		if strings.Contains(upper, "CREATE TABLE") && !strings.Contains(stmt, "dmp_center.") {
			stmt = addDatabasePrefix(stmt, "dmp_center")
		}

		desc := truncate(stmt, 80)
		fmt.Printf("[执行] %s\n", desc)

		if err := db.Exec(stmt).Error; err != nil {
			fmt.Printf("  !!! 执行失败: %v\n", err)
			centerFail++
			continue
		}
		centerSuccess++
	}

	fmt.Printf("\n>>> dmp_center 初始化完成: 成功=%d, 跳过=%d, 失败=%d\n", centerSuccess, centerSkip, centerFail)

	// ============================================================
	// 第二阶段：查询 project_data，为每个项目创建数据库和表
	// ============================================================
	fmt.Println("\n========================================")
	fmt.Println(">>> 第二阶段：创建项目日志数据库")
	fmt.Println("========================================")

	// 查询所有项目
	var projects []struct {
		ProjectAlias string `gorm:"column:project_alias"`
		ProjectName  string `gorm:"column:project_name"`
	}
	if err := db.Raw("SELECT project_alias, project_name FROM dmp_center.project_data").Scan(&projects).Error; err != nil {
		fmt.Printf("!!! 查询 project_data 失败: %v\n", err)
		fmt.Println(">>> 跳过项目数据库创建（project_data 表可能为空或不存在）")
		return
	}

	if len(projects) == 0 {
		fmt.Println(">>> project_data 中暂无项目，跳过项目数据库创建")
		return
	}

	fmt.Printf(">>> 发现 %d 个项目\n", len(projects))

	// 读取项目表结构 SQL
	projectSQLPath := filepath.Join(dir, "doris_project_table.sql")
	projectSQL := readSQLFile(projectSQLPath)
	projectStmts := parseSQLStatements(projectSQL)

	for _, proj := range projects {
		alias := proj.ProjectAlias
		if alias == "" {
			continue
		}
		dbName := core.GetProjectDatabaseName(alias) // dmp_{alias}

		fmt.Printf("\n--- 项目: %s (%s) -> 数据库: %s ---\n", proj.ProjectName, alias, dbName)

		// 创建数据库
		createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName)
		if err := db.Exec(createDBSQL).Error; err != nil {
			fmt.Printf("  !!! 创建数据库 %s 失败: %v\n", dbName, err)
			continue
		}
		fmt.Printf("  [OK] 数据库 %s 就绪\n", dbName)

		// 执行项目级建表语句
		projSuccess := 0
		projSkip := 0
		for _, stmt := range projectStmts {
			upper := strings.ToUpper(strings.TrimSpace(stmt))

			// 跳过非 CREATE TABLE 语句（如 USE, CREATE DATABASE 等）
			if !strings.Contains(upper, "CREATE TABLE") {
				projSkip++
				continue
			}

			// 添加数据库前缀和 IF NOT EXISTS
			finalStmt := addDatabasePrefix(stmt, dbName)
			if !strings.Contains(strings.ToUpper(finalStmt), "IF NOT EXISTS") {
				finalStmt = strings.Replace(finalStmt, "CREATE TABLE", "CREATE TABLE IF NOT EXISTS", 1)
			}

			desc := truncate(finalStmt, 80)
			fmt.Printf("  [执行] %s\n", desc)

			if err := db.Exec(finalStmt).Error; err != nil {
				fmt.Printf("    !!! 执行失败: %v\n", err)
				continue
			}
			projSuccess++
		}
		fmt.Printf("  >>> 项目 %s 完成: 成功=%d, 跳过=%d\n", alias, projSuccess, projSkip)
	}

	fmt.Println("\n========================================")
	fmt.Println(">>> 全部初始化完成！")
	fmt.Println("========================================")
}
