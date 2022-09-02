package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type User struct {
	gorm.Model
	Name string `gorm:"name;unique"`
}

func main() {
	sleepDuration := time.Millisecond * 500
	db := newDB()

	if err := db.Migrator().DropTable(new(User)); err != nil {
		panic(err)
	}
	time.Sleep(sleepDuration)
	if err := db.Migrator().CreateTable(new(User)); err != nil {
		panic(err)
	}
	time.Sleep(sleepDuration)

	saved := User{Name: "user1"}
	log.Println("//==============================================")
	log.Println("Try to save an user")
	if err := db.Create(&saved).Error; err != nil {
		panic(err)
	}
	log.Println("================================================//")
	time.Sleep(sleepDuration)

	log.Println("//==============================================")
	log.Println("Try to update an user")
	saved.Name = "updated-user1"
	if err := db.Updates(&saved).Error; err != nil {
		panic(err)
	}
	log.Println("================================================//")
	time.Sleep(sleepDuration)

	log.Println("//==============================================")
	log.Println("Try to find an user by calling First()")
	var find User
	if err := db.First(&find, "name = ?", saved.Name).Error; err != nil {
		panic(err)
	}
	log.Println("================================================//")
	time.Sleep(sleepDuration)

	log.Println("//==============================================")
	log.Println("Try to find an user by calling exec()")
	db.Exec("SELECT * FROM users").Rows()
	log.Println("================================================//")
	time.Sleep(sleepDuration)

	log.Println("//==============================================")
	log.Println("Try to find an user with manual switching")
	if err := db.Clauses(dbresolver.Write).First(&find, "name = ?", saved.Name).Error; err != nil {
		panic(err)
	}
	log.Println("================================================//")
	time.Sleep(sleepDuration)

	log.Println("//==============================================")
	log.Println("Try to delete an user")
	if err := db.Delete(&saved).Error; err != nil {
		panic(err)
	}
	log.Println("================================================//")
	time.Sleep(sleepDuration)

	// Output
	//2022/09/01 18:10:54 //==============================================
	//2022/09/01 18:10:54 Try to save an user
	//[Callback - Create] >> current user: Master(mydb_user@192.168.96.1), err: <nil>
	//2022/09/01 18:10:54 ================================================//
	//2022/09/01 18:10:55 //==============================================
	//2022/09/01 18:10:55 Try to update an user
	//[Callback - Update] >> current user: Master(mydb_user@192.168.96.1), err: <nil>
	//2022/09/01 18:10:55 ================================================//
	//2022/09/01 18:10:55 //==============================================
	//2022/09/01 18:10:55 Try to find an user by calling First()
	//[Callback - Query] >> current user: Slave(mydb_slave_user@192.168.96.1), err: <nil>
	//2022/09/01 18:10:55 ================================================//
	//2022/09/01 18:10:56 //==============================================
	//2022/09/01 18:10:56 Try to find an user by calling exec()
	//[Callback - Row] >> current user: Slave(mydb_slave_user@192.168.96.1), err: <nil>
	//2022/09/01 18:10:56 ================================================//
	//2022/09/01 18:10:56 //==============================================
	//2022/09/01 18:10:56 Try to find an user with manual switching
	//[Callback - Query] >> current user: Master(mydb_user@192.168.96.1), err: <nil>
	//2022/09/01 18:10:56 ================================================//
	//2022/09/01 18:10:57 //==============================================
	//2022/09/01 18:10:57 Try to delete an user
	//[Callback - Delete] >> current user: Master(mydb_user@192.168.96.1), err: <nil>
	//2022/09/01 18:10:57 ================================================//
}

func newDB() *gorm.DB {
	var (
		masterDSN = "mydb_user:mydb_pwd@tcp(127.0.0.1:4406)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
		slaveDSN  = "mydb_slave_user:mydb_slave_pwd@tcp(127.0.0.1:5506)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	)

	db, err := gorm.Open(mysql.Open(masterDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}

	db.Use(dbresolver.Register(dbresolver.Config{
		// use `masterDSN` as sources (DB's default connection)
		// `slaveDSN` as replicas
		Replicas: []gorm.Dialector{mysql.Open(slaveDSN)},
	}))

	db.Callback().Create().After("gorm:db_resolver").Register("callback:create", func(db *gorm.DB) {
		printCallbackEvent("[Callback - Create]", db)
	})
	db.Callback().Delete().After("gorm:db_resolver").Register("callback:delete", func(db *gorm.DB) {
		printCallbackEvent("[Callback - Delete]", db)
	})
	db.Callback().Query().After("gorm:db_resolver").Register("callback:query", func(db *gorm.DB) {
		printCallbackEvent("[Callback - Query]", db)
	})
	db.Callback().Row().After("gorm:db_resolver").Register("callback:row", func(db *gorm.DB) {
		printCallbackEvent("[Callback - Row]", db)
	})
	db.Callback().Update().After("gorm:db_resolver").Register("callback:update", func(db *gorm.DB) {
		printCallbackEvent("[Callback - Update]", db)
	})
	return db
}

func printCallbackEvent(title string, db *gorm.DB) {
	user, err := currentDBUser(db.Statement.ConnPool.(*sql.DB))
	fmt.Printf("%s >> current user: %s(%s), err: %v\n", title, resolveUser(user), user, err)
}

func resolveUser(user string) string {
	if strings.Contains(user, "mydb_user") {
		return "Master"
	}
	if strings.Contains(user, "mydb_slave_user") {
		return "Slave"
	}
	return user
}

func currentDBUser(db *sql.DB) (string, error) {
	conn, err := db.Conn(context.Background())
	if err != nil {
		return "", err
	}
	defer conn.Close()

	var user string
	err = conn.QueryRowContext(context.Background(), "SELECT USER()").Scan(&user)
	if err != nil {
		return "", err
	}
	return user, nil
}
