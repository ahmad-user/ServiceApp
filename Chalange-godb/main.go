package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "12345678"
	dbname   = "laundry"
)

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
var db *sql.DB

type Customer struct {
	PhoneNumber int
	Name        string
}

type Order struct {
	No            int
	TanggalMasuk  time.Time
	TanggalKeluar time.Time
	Diterima      string
}

type Transaction struct {
	No        int
	Pelayanan string
	Jumlah    int
	Satuan    string
	Harga     int
	Total     int
}

func main() {
	var err error
	db, err = connectDB()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	for {
		menuUtama()
	}

}

// Menu utama
func menuUtama() {
	for {
		fmt.Println("\n== Menu Utama ==")
		fmt.Println("1. Menu Customer")
		fmt.Println("2. Menu Pesanan")
		fmt.Println("3. Menu Transaction")
		fmt.Println("0. Keluar")

		var choice int
		fmt.Print("Pilih menu: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			MenuCustomer()
		case 2:
			MenuPesanan()
		case 3:
			MenuTransaction()
		case 0:
			fmt.Println("Keluar dari program.")
			os.Exit(0)
		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}

// tampilan menu customer
func MenuCustomer() {
	for {
		fmt.Println("\n== Menu ==")
		fmt.Println("1. Lihat Customer")
		fmt.Println("2. Tambah Customer")
		fmt.Println("3. Edit Customer")
		fmt.Println("4. Delete Customer")
		fmt.Println(" ")
		fmt.Println("9. Kembali")
		fmt.Println("0. Keluar")

		var choice int
		fmt.Print("Pilih menu: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			viewCustomers()
		case 2:
			insertCustomer()
		case 3:
			editCustomer()
		case 4:
			deleteCustomer()
		case 9:
			menuUtama()
		case 0:
			fmt.Println("Keluar dari program.")
			os.Exit(0)
		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}

// tampilan menu pesanan
func MenuPesanan() {
	for {
		fmt.Println("\n== Menu ==")
		fmt.Println("1. Lihat Pesanan")
		fmt.Println("2. Tambah Pesanan ")
		fmt.Println("3. Edit pesanan")
		fmt.Println("4. Hapus pesanan")
		fmt.Println(" ")
		fmt.Println("9. Kembali")
		fmt.Println("0. Keluar")

		var choice int
		fmt.Print("Pilih menu: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			viewPesanan()
		case 2:
			insertPesanan()
		case 3:
			editPesanan()
		case 4:
			deletePesanan()
		case 9:
			menuUtama()
		case 0:
			fmt.Println("Keluar dari program.")
			os.Exit(0)
		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}
func MenuTransaction() {
	for {
		fmt.Println("\n== Menu Transaction ==")
		fmt.Println("1. Lihat Transaksi")
		fmt.Println("2. Tambah Transaksi")
		fmt.Println(" ")
		fmt.Println("9. Kembali")
		fmt.Println("0. Keluar")

		var choice int
		fmt.Print("Pilih menu: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			viewTransactions()
		case 2:
			insertTransaction()
		case 9:
			menuUtama()
		case 0:
			fmt.Println("Keluar dari program.")
			os.Exit(0)
		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}

// koneksi database
func connectDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("Berhasil terkoneksi ke database!")
	return db, nil
}

// tampilkan costumer
func viewCustomers() {
	rows, err := db.Query("SELECT * FROM tbl_customer")
	if err != nil {
		log.Fatal("Error fetching customers:", err)
	}
	defer rows.Close()

	fmt.Println("\n== Daftar Customer ==")
	var count int
	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.PhoneNumber, &customer.Name)
		if err != nil {
			log.Fatal("Error scanning customers:", err)
		}
		fmt.Printf("Hp: %d, No. Name: %s\n", customer.PhoneNumber, customer.Name)
		count++
	}
	if count == 0 {
		fmt.Println("Tidak ada data customer yang tersedia.")
	}
}

// tambahkan data customer
func insertCustomer() {
	var phone, name string

	fmt.Print("Masukkan nomor HP customer: ")
	fmt.Scanln(&phone)
	fmt.Print("Masukkan nama customer: ")
	fmt.Scanln(&name)

	_, err := strconv.Atoi(phone)
	if err != nil {
		fmt.Println("Nomor HP harus berupa angka.")
		return
	}

	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM tbl_customer WHERE no_hp = $1", phone)
	err = row.Scan(&count)
	if err != nil {
		log.Fatal("Error checking phone number:", err)
	}
	if count > 0 {
		fmt.Println("Nomor HP sudah digunakan oleh customer lain.")
		return
	}

	_, err = db.Exec("INSERT INTO tbl_customer VALUES ($1, $2)", phone, name)
	if err != nil {
		log.Fatal("Error inserting customer:", err)
	}
	fmt.Println("Customer berhasil ditambahkan.")
}
func editCustomer() {
	var phone, newName string

	fmt.Print("Masukkan nomor HP customer yang akan diubah: ")
	fmt.Scanln(&phone)
	fmt.Print("Masukkan nama baru customer: ")
	fmt.Scanln(&newName)

	_, err := strconv.Atoi(phone)
	if err != nil {
		fmt.Println("Nomor HP harus berupa angka.")
		return
	}

	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM tbl_customer WHERE no_hp = $1", phone)
	err = row.Scan(&count)
	if err != nil {
		log.Fatal("Error checking phone number:", err)
	}
	if count == 0 {
		fmt.Println("Customer dengan nomor HP tersebut tidak ditemukan.")
		return
	}

	_, err = db.Exec("UPDATE tbl_customer SET nama_Customer = $1 WHERE no_hp = $2", newName, phone)
	if err != nil {
		log.Fatal("Error updating customer:", err)
	}
	fmt.Println("Customer berhasil diubah.")
}

// delete customer
func deleteCustomer() {
	var no_hp int
	fmt.Print("Masukkan No Customer customer yang akan dihapus: ")
	fmt.Scanln(&no_hp)

	result, err := db.Exec("DELETE FROM tbl_customer WHERE no_hp = $1", no_hp)
	if err != nil {
		log.Fatal("Error deleting customer:", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		fmt.Println("Customer tidak ditemukan.")
	} else {
		fmt.Println("Customer berhasil dihapus.")
	}
}

// tampilkan data pesanan
func viewPesanan() {
	rows, err := db.Query("SELECT * FROM tbl_pesanan")
	if err != nil {
		log.Fatal("Error fetching orders:", err)
	}
	defer rows.Close()

	fmt.Println("\n== Daftar Pesanan ==")
	var count int
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.No, &order.TanggalMasuk, &order.TanggalKeluar, &order.Diterima)
		if err != nil {
			log.Fatal("Error scanning pesanan:", err)
		}
		fmt.Printf("No: %d, Tanggal Masuk: %s, Tanggal Keluar: %s, Diterima oleh: %s\n", order.No, order.TanggalMasuk, order.TanggalKeluar, order.Diterima)
		count++
	}
	if count == 0 {
		fmt.Println("Tidak ada data pesanan yang tersedia.")
	}
}

// tambahkan data pensanan
func insertPesanan() {
	var no int
	var tglMasuk, tglKeluar, diterima string
	fmt.Print("Masukkan nomor pesanan: ")
	fmt.Scanln(&no)
	fmt.Print("Masukkan tanggal masuk (YYYY-MM-DD): ")
	fmt.Scanln(&tglMasuk)
	fmt.Print("Masukkan tanggal keluar (YYYY-MM-DD): ")
	fmt.Scanln(&tglKeluar)
	fmt.Print("Masukkan nama yang menerima pesanan: ")
	fmt.Scanln(&diterima)

	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM tbl_pesanan WHERE no = $1", no)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal("Error checking order number:", err)
	}
	if count > 0 {
		fmt.Println("Nomor pesanan sudah digunakan.")
		return
	}

	_, err = db.Exec("INSERT INTO tbl_pesanan (no, tgl_masuk, tgl_keluar, diterima) VALUES ($1, $2, $3, $4)", no, tglMasuk, tglKeluar, diterima)
	if err != nil {
		log.Fatal("Error inserting pesanan:", err)
	}
	fmt.Println("Pesanan berhasil ditambahkan.")
}

// edit pesanan
func editPesanan() {
	var no int
	fmt.Print("Masukkan nomor pesanan yang akan diedit: ")
	fmt.Scanln(&no)

	var tglMasuk, tglKeluar, diterima string
	fmt.Print("Masukkan tanggal masuk (YYYY-MM-DD): ")
	fmt.Scanln(&tglMasuk)
	fmt.Print("Masukkan tanggal keluar (YYYY-MM-DD): ")
	fmt.Scanln(&tglKeluar)
	fmt.Print("Masukkan nama yang menerima pesanan: ")
	fmt.Scanln(&diterima)

	_, err := db.Exec("UPDATE tbl_pesanan SET tgl_masuk = $1, tgl_keluar = $2, diterima = $3 WHERE no = $4", tglMasuk, tglKeluar, diterima, no)
	if err != nil {
		log.Fatal("Error updating pesanan:", err)
	}

	fmt.Println("Pesanan berhasil diupdate.")
}
func deletePesanan() {
	var no int
	fmt.Print("Masukkan nomor pesanan yang akan dihapus: ")
	fmt.Scanln(&no)

	result, err := db.Exec("DELETE FROM tbl_pesanan WHERE no = $1", no)
	if err != nil {
		log.Fatal("Error deleting pesanan:", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		fmt.Println("Pesanan tidak ditemukan.")
	} else {
		fmt.Println("Pesanan berhasil dihapus.")
	}
}

// tampilkan data transaction
func viewTransactions() {
	rows, err := db.Query("SELECT * FROM tbl_transaction")
	if err != nil {
		log.Fatal("Error fetching transactions:", err)
	}
	defer rows.Close()

	fmt.Println("\n== Daftar Transaksi ==")
	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(&transaction.No, &transaction.Pelayanan, &transaction.Jumlah, &transaction.Satuan, &transaction.Harga, &transaction.Total)
		if err != nil {
			log.Fatal("Error scanning transaction:", err)
		}
		fmt.Printf("No: %d, Pelayanan: %s, Jumlah: %d, Satuan: %s, Harga: %d, Total: %d\n", transaction.No, transaction.Pelayanan, transaction.Jumlah, transaction.Satuan, transaction.Harga, transaction.Total)
	}
}

// tambahkan data transaksi
func insertTransaction() {
	var no, jumlah, harga int
	var pelayanan, satuan string
	fmt.Print("Masukkan nomor transaksi: ")
	fmt.Scanln(&no)
	fmt.Print("Masukkan jenis pelayanan: ")
	fmt.Scanln(&pelayanan)
	fmt.Print("Masukkan jumlah: ")
	fmt.Scanln(&jumlah)
	fmt.Print("Masukkan satuan: ")
	fmt.Scanln(&satuan)
	fmt.Print("Masukkan harga: ")
	fmt.Scanln(&harga)

	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM tbl_transaction WHERE no = $1", no)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal("Error checking transaction number:", err)
	}
	if count > 0 {
		fmt.Println("Nomor transaksi sudah digunakan.")
		return
	}

	if pelayanan == "" {
		fmt.Println("Jenis pelayanan harus diisi.")
		return
	}

	if jumlah <= 0 {
		fmt.Println("Jumlah harus lebih besar dari 0.")
		return
	}

	if satuan == "" {
		fmt.Println("Satuan harus diisi.")
		return
	}

	if harga <= 0 {
		fmt.Println("Harga harus lebih besar dari 0.")
		return
	}

	_, err = db.Exec("INSERT INTO tbl_transaction (no, pelayanan, jumlah, satuan, harga, total) VALUES ($1, $2, $3, $4, $5, $6)", no, pelayanan, jumlah, satuan, harga, harga*jumlah)
	if err != nil {
		log.Fatal("Error inserting transaction:", err)
	}
	fmt.Println("Transaksi berhasil ditambahkan.")
}
