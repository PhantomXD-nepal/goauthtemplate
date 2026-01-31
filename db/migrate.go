package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/PhantomXD-nepal/goauthtemplate/internal/config"

	"github.com/golang-migrate/migrate/v4"
	mysqldriver "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true&multiStatements=true",
		config.Envs.DBUser,
		config.Envs.DBPassword,
		config.Envs.DBAddress,
		config.Envs.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping db:", err)
	}

	driver, err := mysqldriver.WithInstance(db, &mysqldriver.Config{})
	if err != nil {
		log.Fatal("failed to create migration driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal("failed to create migrate instance:", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("usage: migrate [up|down|force|version]")
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Println("✓ migrations applied")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Println("✓ migrations rolled back")

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("usage: migrate force <version>")
		}
		var version int
		fmt.Sscanf(os.Args[2], "%d", &version)

		if err := m.Force(version); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("✓ forced version %d\n", version)

	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			log.Fatal(err)
		}
		if dirty {
			fmt.Printf("version %d (dirty)\n", v)
		} else {
			fmt.Printf("version %d\n", v)
		}

	default:
		log.Fatal("invalid command: use up, down, force, version")
	}
}
