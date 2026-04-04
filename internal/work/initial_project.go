package work

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

/*
 After i create the structure for the project
i get the path for the project from projectName
then go inside it.
Then install the dependencies
*/

func InstallDependencies(projectName, framework, dbDriver, sqlDriver string) error {
	projectPath, err := filepath.Abs(projectName)
	if err != nil {
		return fmt.Errorf("cannot get project path: %w", err)
	}

	if err := os.Chdir(projectPath); err != nil {
		return fmt.Errorf("cannot change to project directory: %w", err)
	}

	fmt.Println("Current directory:", projectPath)

	if framework != "None" {
		frameworkPkg := getFrameworkPackage(framework)
		fmt.Println("Installing framework:", frameworkPkg)
		if frameworkPkg != " " {
			cmd := exec.Command("go", "get", frameworkPkg)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install framework %s: %w", framework, err)
			}
		}
	}

	if dbDriver != "None" {
		driverPackage := getDBDriverPackage(dbDriver)
		if driverPackage != "" {
			fmt.Println("Installing DB driver:", driverPackage)
			cmd := exec.Command("go", "get", driverPackage)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install DB driver %s: %w", dbDriver, err)
			}
		}
	}

	// 5️⃣ Install SQL/ORM driver
	if sqlDriver != "None" {
		sqlPackage := getSQLDriverPackage(sqlDriver)
		if sqlPackage != "" {
			fmt.Println("Installing SQL driver:", sqlPackage)
			cmd := exec.Command("go", "get", sqlPackage)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install SQL driver %s: %w", sqlDriver, err)
			}
		}
	}

	fmt.Println("✅ All dependencies installed successfully!")
	return nil
}

func getFrameworkPackage(name string) string {
	switch name {
	case "Gin":
		return "github.com/gin-gonic/gin"
	case "Echo":
		return "github.com/labstack/echo/v4"
	case "Fiber":
		return "github.com/gofiber/fiber/v2"
	case "Chi":
		return "github.com/go-chi/chi/v5"
	default:
		return ""
	}
}

func getDBDriverPackage(db string) string {
	switch db {
	case "PostgreSQL":
		return "github.com/jackc/pgx/v5"
	case "MySQL":
		return "github.com/go-sql-driver/mysql"
	case "SQLite":
		return "github.com/mattn/go-sqlite3"
	case "MongoDB":
		return "go.mongodb.org/mongo-driver/mongo"
	default:
		return ""
	}
}

func getSQLDriverPackage(sql string) string {
	switch sql {
	case "GORM":
		return "gorm.io/gorm"
	case "sqlx":
		return "github.com/jmoiron/sqlx"
	case "sqlc":
		return "github.com/kyleconroy/sqlc"
	case "pgx":
		return "github.com/jackc/pgx/v5"
	default:
		return ""
	}
}
