package backren

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func TestConn31(t *testing.T) {
	// Buat string koneksi
	// Format: "username:password@protocol(address)/dbname?param=value"
	// Misalnya: "root:password@tcp(127.0.0.1:3306)/dbname"
	// db, err := sql.Open("mysql", "wacybero_waren:Cyberorenserver12@tcp(bogor.gusti.id:3306)/wacybero_wa")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer db.Close()
	db, err := SqlConn()
	if err != nil {
		panic(err.Error())
	}

	// Tes koneksi
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to the database")

	// Query sederhana
	rows, err := db.Query("SELECT name, username FROM users")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	// fmt.Println("rows")
	// fmt.Println(rows)
	var users []User
	for rows.Next() {
		// var id int
		// var name string
		var user User
		if err := rows.Scan(&user.Name, &user.Username); err != nil {
			panic(err.Error())
		}
		users = append(users, user)

		fmt.Printf("ID: %s, Username: %s\n", user.Name, user.Username)
	}
	fmt.Println("users")
	fmt.Println(users)
	var name string
	var username string
	var role string
	var limit_device string

	err = db.QueryRow("SELECT name, username, role, limit_device FROM users WHERE username = ?", "jx").Scan(&name, &username, &role, &limit_device)
	if err != nil {
		fmt.Printf("ID: %d, Username: %s, error:  %s\n", 0, "", err)
	}
	fmt.Printf("name: %s, Username: %s, role: %s, limit_device: %s\n", name, username, role, limit_device)

}

func TestInserUser(t *testing.T) {
	// Buat string koneksi
	// Format: "username:password@protocol(address)/dbname?param=value"
	// Misalnya: "root:password@tcp(127.0.0.1:3306)/dbname"
	db, err := sql.Open("mysql", "wacybero_waren:Cyberorenserver12@tcp(bogor.gusti.id:3306)/wacybero_wa")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Tes koneksi
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to the database")

	query := "INSERT INTO users (name, username, role, password, limit_device) VALUES (?, ?, ?, ?, ?)"

	password := "passwordbaru"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	// Siapkan data yang akan disisipkan
	name := "daslan"
	username := "jx"
	role := "user"
	limit_device := "5"

	// Lakukan penyisipan data menggunakan Prepare statement
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	// Eksekusi perintah untuk menyisipkan data
	_, err = stmt.Exec(name, username, role, string(hashedPassword), limit_device)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Data berhasil disisipkan ke dalam tabel 'users'")

	// Query sederhana
	rows, err := db.Query("SELECT id, username FROM users")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			panic(err.Error())
		}
		fmt.Printf("ID: %d, Username: %s\n", id, name)
	}
}

func TestGenerateKey(t *testing.T) {
	privateKey, publicKey := GenerateKey()
	fmt.Println("privateKey : ", privateKey)
	fmt.Println("publicKey : ", publicKey)
}
