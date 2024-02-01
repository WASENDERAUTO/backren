package backren

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var (
	user     User
	response Response
)

func SqlConn() (*sql.DB, error) {
	db, err := sql.Open("mysql", "wacybero_waren:Cyberorenserver12@tcp(bogor.gusti.id:3306)/wacybero_wa")
	if err != nil {
		return db, fmt.Errorf("db: %v", err)
	}
	// Tes koneksi
	err = db.Ping()
	if err != nil {
		return db, fmt.Errorf("tes koneksi: %v", err)
	}
	return db, err
}

func LogIn(PASETOPRIVATEKEYENV string, r *http.Request) string {
	response.Status = 400
	db, err := SqlConn()
	if err != nil {
		response.Message = "error 1: " + err.Error()
		return GCFReturnStruct(response)
	}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(response)
	}
	if user.Email == "" || user.Password == "" {
		response.Message = "mohon untuk melengkapi"
		return GCFReturnStruct(response)
	}
	var email string
	var password string
	var username string
	err = db.QueryRow("SELECT username, email, password FROM users_store WHERE email = ?", user.Email).Scan(&username, &email, &password)
	if err != nil {
		response.Message = "error: email tidak ada" + err.Error()
		return GCFReturnStruct(response)
	}
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if err != nil {
		response.Message = "error: Kata sandi tidak cocok." + err.Error()
		return GCFReturnStruct(response)
	}
	tokenstring, err := Encode(username, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		response.Message = "error 3: " + err.Error()
		return GCFReturnStruct(response)
	}
	data := map[string]interface{}{
		"status":  200,
		"message": "Berhasil Login",
		"data": map[string]interface{}{
			"token":    tokenstring,
			"username": username,
		},
	}
	return GCFReturnStruct(data)
}

func InsertUserStore(r *http.Request) string {
	response.Status = 400
	db, err := SqlConn()
	if err != nil {
		response.Message = "error 1: " + err.Error()
		return GCFReturnStruct(response)
	}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(response)
	}
	if user.Email == "" || user.Password == "" || user.Name == "" || user.PhoneNumber == "" || user.Username == "" {
		response.Message = "mohon untuk melengkapi"
		return GCFReturnStruct(response)
	}
	rows_username, err := db.Query("SELECT * FROM users_store WHERE username = ?", user.Username)
	if err != nil {
		response.Message = "error: " + err.Error()
		return GCFReturnStruct(response)
	}
	var count_username int
	for rows_username.Next() {
		count_username++
	}
	if count_username > 0 {
		response.Message = "error: username sudah terdaftar" + err.Error()
		return GCFReturnStruct(response)
	}
	if err := rows_username.Err(); err != nil {
		response.Message = "error: " + err.Error()
		return GCFReturnStruct(response)
	}
	rows_email, err := db.Query("SELECT * FROM users_store WHERE email = ?", user.Email)
	if err != nil {
		response.Message = "error: " + err.Error()
		return GCFReturnStruct(response)
	}
	var count_email int
	for rows_email.Next() {
		count_email++
	}
	if count_email > 0 {
		response.Message = "error: email sudah terdaftar" + err.Error()
		return GCFReturnStruct(response)
	}
	if err := rows_email.Err(); err != nil {
		response.Message = "error: " + err.Error()
		return GCFReturnStruct(response)
	}
	rows_phone_number, err := db.Query("SELECT * FROM users_store WHERE phone_number = ?", user.PhoneNumber)
	if err != nil {
		response.Message = "error: " + err.Error()
		return GCFReturnStruct(response)
	}
	var count_phone_number int
	for rows_phone_number.Next() {
		count_phone_number++
	}
	if count_phone_number > 0 {
		response.Message = "error: nomor hp sudah terdaftar" + err.Error()
		return GCFReturnStruct(response)
	}
	if err := rows_phone_number.Err(); err != nil {
		response.Message = "error: " + err.Error()
		return GCFReturnStruct(response)
	}
	password := user.Password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Message = "error 2: " + err.Error()
		return GCFReturnStruct(response)
	}

	query := "INSERT INTO users_store (name, username, email, phone_number, license, password) VALUES (?, ?, ?, ?, ?, ?)"

	// Lakukan penyisipan data menggunakan Prepare statement
	stmt, err := db.Prepare(query)
	if err != nil {
		// panic(err.Error())
		response.Message = "error 3: " + err.Error()
		return GCFReturnStruct(response)
	}
	defer stmt.Close()

	// Eksekusi perintah untuk menyisipkan data
	_, err = stmt.Exec(user.Name, user.Username, user.Email, user.PhoneNumber, "kosong", string(hashedPassword))
	if err != nil {
		response.Message = "error 4: " + err.Error()
		return GCFReturnStruct(response)
	}

	response.Message = "Berhasil SignUp"
	data := map[string]interface{}{
		"status":  200,
		"message": response.Message,
		"data": map[string]interface{}{
			"name":         user.Name,
			"username":     user.Username,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
		},
	}
	return GCFReturnStruct(data)
}

func InsertUserApp(PASETOPUBLICKEYENV string, r *http.Request) string {
	response.Status = 400
	db, err := SqlConn()
	if err != nil {
		response.Message = "error: " + err.Error()
		return GCFReturnStruct(response)
	}

	user, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		response.Message = err.Error()
		return GCFReturnStruct(response)
	}

	if user.Username != "admin" {
		response.Message = "Kamu bukan admin"
		return GCFReturnStruct(response)
	}

	username := GetUsername(r)
	if username == "" {
		response.Message = "Wrong parameter"
		return GCFReturnStruct(response)
	}

	var name string
	var password string

	err = db.QueryRow("SELECT name, username, password FROM users_store WHERE username = ?", username).Scan(&name, &username, &password)
	if err != nil {
		response.Message = "error: " + err.Error()
		return GCFReturnStruct(response)
	}

	random := make([]byte, 16)
	_, err = rand.Read(random)
	if err != nil {
		response.Message = "error random: " + err.Error()
		return GCFReturnStruct(response)
	}

	license := username + hex.EncodeToString(random)

	_, err = db.Exec("UPDATE users_store SET license = ? WHERE username = ?", license, username)
	if err != nil {
		response.Message = "error : " + err.Error()
		return GCFReturnStruct(response)
	}

	query := "INSERT INTO users (name, username, role, password, limit_device) VALUES (?, ?, ?, ?, ?)"

	// Siapkan data yang akan disisipkan
	role := "user"
	limit_device := "5"

	// Lakukan penyisipan data menggunakan Prepare statement
	stmt, err := db.Prepare(query)
	if err != nil {
		response.Message = "error : " + err.Error()
		return GCFReturnStruct(response)
	}
	defer stmt.Close()

	// Eksekusi perintah untuk menyisipkan data
	_, err = stmt.Exec(name, username, role, password, limit_device)
	if err != nil {
		response.Message = "error : " + err.Error()
		return GCFReturnStruct(response)
	}

	response.Message = "Berhasil Tambah User App"
	data := map[string]interface{}{
		"status":  200,
		"message": response.Message,
		"data": map[string]interface{}{
			"name":         name,
			"username":     username,
			"role":         role,
			"limit_device": limit_device,
		},
	}
	return GCFReturnStruct(data)
}

func GetuserByAdmin(PASETOPUBLICKEYENV string, r *http.Request) string {
	response.Status = 400
	db, err := SqlConn()
	if err != nil {
		response.Message = "error: " + err.Error()
		return GCFReturnStruct(response)
	}

	user, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		response.Message = err.Error()
		return GCFReturnStruct(response)
	}

	if user.Username != "admin" {
		response.Message = "Kamu bukan admin"
		return GCFReturnStruct(response)
	}
	rows, err := db.Query("SELECT * FROM users_store")
	if err != nil {
		response.Message = "error: " + err.Error()
		return GCFReturnStruct(response)
	}
	var users []User
	for rows.Next() {
		var id int
		var user User
		if err := rows.Scan(&id, &user.Name, &user.Username, &user.Email, &user.PhoneNumber, &user.License, &user.Password); err != nil {
			panic(err.Error())
		}
		users = append(users, user)

		// fmt.Printf("ID: %s, Username: %s\n", user.Name, user.Username)
	}
	response.Message = "Get Success"
	data := map[string]interface{}{
		"status":  200,
		"message": response.Message,
		"data":    users,
	}
	return GCFReturnStruct(data)
}

// get user login
func GetUserLogin(PASETOPUBLICKEYENV string, r *http.Request) (Payload, error) {
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		return payload, err
	}
	return payload, nil
}

// get id
func GetUsername(r *http.Request) string {
	return r.URL.Query().Get("username")
}

// return json string
func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}
