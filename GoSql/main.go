package main

import (
	"GoSql/config"
	"GoSql/users"
	"GoSql/views"
	"fmt"
)

func main() {
	// input database
	database := config.InitMysql()
	// migrate database
	config.Migrate(database)
	var input = -1
	for input != 0 {
		fmt.Println(" Pilih menu")
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Println("3. Lupa Password")
		fmt.Println("0. Exit")
		fmt.Print("Masukkan pilihan: ")
		fmt.Scanln(&input)
		// menu login
		if input == 1 {
			var isRunning bool = true
			for isRunning {
				var hp string
				var password string
				fmt.Println("Masukkan HP")
				fmt.Scanln(&hp)
				fmt.Println("Masukkan Password")
				fmt.Scanln(&password)
				loggedIn, err := users.Login(database, hp, password)
				if err == nil {
					fmt.Println("Selamat Datang,", loggedIn.Nama)
					if loggedIn != (users.Tbl_user{}) {
						//  login berhasil
						var menuInput = -1
						for menuInput != 0 {
							fmt.Println("Pilih menu:")
							fmt.Printf("Saldo anda : %.2f\n", loggedIn.Balance)
							fmt.Println("1. Tambah User")
							fmt.Println("2. Lihat User")
							fmt.Println("3. Edit User")
							fmt.Println("4. Hapus User")
							fmt.Println("5. Top Up Saldo")
							fmt.Println("6. Transfer")
							fmt.Println("7. History Transfer")
							fmt.Println("8. History Top Up")
							fmt.Println("9. Cari User")
							fmt.Println("0. Keluar")
							fmt.Print("Masukkan pilihan: ")
							fmt.Scanln(&menuInput)

							switch menuInput {
							case 1:
								// form untuk menambahkan users
								var newUser users.Tbl_user
								fmt.Print("Masukkan nama : ")
								fmt.Scanln(&newUser.Nama)
								fmt.Print("Masukkan nomor HP : ")
								fmt.Scanln(&newUser.HP)
								fmt.Print("Masukkan password : ")
								fmt.Scanln(&newUser.Password)
								fmt.Print("Masukkan alamat : ")
								fmt.Scanln(&newUser.Alamat)

								success, err := users.TambahUsers(database, newUser)
								if err != nil {
									fmt.Println("Gagal menambahkan barang:", err)
								} else if success {
									fmt.Println("--Barang berhasil ditambahkan--")
								} else {
									fmt.Println("Gagal menambahkan barang")
								}
							case 2:
								// tampilkan daftar users
								daftarBarang, err := users.LihatUsers(database, loggedIn.ID)
								if err != nil {
									fmt.Println("Gagal mendapatkan daftar user:", err)
								} else {
									fmt.Println("Daftar Users:")
									for i, users := range daftarBarang {
										fmt.Printf("%d. HP: %s, Nama: %s, Password: %s, Alamat: %s\n", i+1, users.HP, users.Nama, users.Password, users.Alamat)
									}
								}
							case 3:
								// edit form users
								var idEdit uint
								fmt.Print("Masukan id yg ingin diedit : ")
								fmt.Scanln(&idEdit)
								var hp string
								fmt.Print("Masukan no telp : ")
								fmt.Scanln(&hp)
								var namaUser string
								fmt.Print("Masukan nama user : ")
								fmt.Scanln(&namaUser)
								var password string
								fmt.Print("Masukan password : ")
								fmt.Scanln(&password)
								var alamat string
								fmt.Print("Masukan alamat  : ")
								fmt.Scanln(&alamat)
								updateData := map[string]interface{}{
									"HP":       hp,
									"Nama":     namaUser,
									"Password": password,
									"Alamat":   alamat,
								}
								success, err := users.EditUsers(database, idEdit, updateData)
								if err != nil {
									fmt.Println("Gagal mengedit users : ", err)
								} else if success {
									fmt.Println("Users berhasil diubah")
								} else {
									fmt.Println("Tidak ada perubahan pada user")
								}
							case 4:
								// kode hapus users
								var idHapus uint
								fmt.Println("Masukan id user yg ingin anda hapus : ")
								fmt.Scanln(&idHapus)

								success, err := users.HapusUsers(database, idHapus)
								if err != nil {
									fmt.Println("Gagal menghapus user ", err)
								} else if success {
									fmt.Println("Anda berhasil menghapus user")
								} else {
									fmt.Println("Tidak ada id user")
								}
							case 5:
								// top up saldo
								fmt.Print("Masukan jumlah Top Up : ")
								var jumlahTopUp float64
								fmt.Scanln(&jumlahTopUp)

								fmt.Print("Masukan metode pembayaran : ")
								var typeTopUp string
								fmt.Scanln(&typeTopUp)

								loggedIn.Balance += jumlahTopUp

								err := views.TopUp(loggedIn.ID, jumlahTopUp, typeTopUp, database)
								if err != nil {
									fmt.Println("gagal melakukan Top Up:", err)
								} else {
									fmt.Println("Sukses")
								}
							case 6:
								// Meminta informasi transfer dari pengguna
								var receiverHP string
								var amount float64
								var typeHP string

								fmt.Println("Masukkan nomor penerima:")
								fmt.Scanln(&receiverHP)

								fmt.Println("Masukkan jumlah yang ditransfer:")
								fmt.Scanln(&amount)

								fmt.Println("Masukan Metode Pembayaran")
								fmt.Scanln(&typeHP)

								senderHP := loggedIn.HP
								loggedIn.Balance -= amount

								// Melakukan transfer
								err := views.TransferBalanceHp(senderHP, receiverHP, amount, typeHP, database)
								if err != nil {
									fmt.Println("Gagal melakukan transfer:", err)
									return
								} else {

									fmt.Println("Transfer dana berhasil")
								}
							case 7:
								// history Transfer
								history, err := views.HistoryTransfer(loggedIn.ID, database)
								if err != nil {
									fmt.Println("Gagal melihat history Transfer: ", err)
								} else {
									fmt.Println("Transfer History:")
									for _, entry := range history {
										fmt.Printf("ID: %d, ReceiverHP: %s, Amount: %.2f, Type: %s, Status: %s\n",
											entry.ID, entry.ReceiverHP, entry.Amount, entry.Type, entry.Status)
									}
								}

							case 8:
								//history TopUp
								var panggil views.Tbl_transfer
								panggil.SenderHP += loggedIn.HP
								history, err := views.HistoryTopUp(loggedIn.ID, database)
								if err != nil {
									fmt.Println("gagal melihat history Top Up ", err)
									return
								}
								fmt.Println("Top Up History:")
								for _, entry := range history {
									fmt.Printf("ID: %d, Jumlah Top Up: %.2f, Tipe: %s, Status: %s\n", entry.ID, entry.JumlahTopUp, entry.Type, entry.Status)
								}
							case 9:
								// Cari users berdasarkan id
								var cari string
								fmt.Print("Masukan no Hp Users : ")
								fmt.Scanln(&cari)

								hasilPencarian, err := users.CariUsers(database, cari)
								if err != nil {
									fmt.Println("Gagal melakukan pencarian : ", err)
								} else if len(hasilPencarian) > 0 {
									fmt.Println("Hasil pencarian : ")
									for i, users := range hasilPencarian {
										fmt.Printf("%d. HP: %s, Nama: %s, Password: %s, Alamat: %s\n", i+1, users.HP, users.Nama, users.Password, users.Alamat)
									}

								} else {
									fmt.Println("Tidak ada kriteria barang")
								}
							case 99:
								fmt.Println("Keluar")
							default:
								fmt.Println("Pilihan tidak valid")
							}
						}
					}
					isRunning = false
				} else {
					var inputExit string
					fmt.Println("Ketik 'EXIT' untuk kembali ke menu sebelumnya")
					fmt.Scanln(&inputExit)
					if inputExit == "EXIT" {
						isRunning = false
					}
				}
			}
			// kalo sukses welcome, kalo gagal isi lagi
		} else if input == 2 {
			// menu register
			var newUser users.Tbl_user
			fmt.Print("Masukkan nama : ")
			fmt.Scanln(&newUser.Nama)
			fmt.Print("Masukkan nomor HP : ")
			fmt.Scanln(&newUser.HP)
			fmt.Print("Masukkan password : ")
			fmt.Scanln(&newUser.Password)
			fmt.Print("Masukkan alamat : ")
			fmt.Scanln(&newUser.Alamat)
			success, err := users.Register(database, newUser)
			if err != nil {
				fmt.Println("terjadi kesalahan(tidak bisa mendaftarkan pengguna)", err.Error())
			}

			if success {
				fmt.Println("--selamat anda telah terdaftar--")

			}
		} else if input == 3 {
			// lupa kata sandi
			var hp string
			fmt.Print("Masukkan nomor HP Anda: ")
			fmt.Scanln(&hp)

			var currentPassword string
			fmt.Print("Masukkan kata sandi saat ini: ")
			fmt.Scanln(&currentPassword)

			var newPassword string
			fmt.Print("Masukkan kata sandi baru: ")
			fmt.Scanln(&newPassword)

			// Verifikasi login pengguna
			user, err := users.Login(database, hp, currentPassword)
			if err != nil {
				fmt.Println("Terjadi kesalahan:", err)
				return
			}
			success, err := user.GantiPassword(database, newPassword)
			if err != nil {
				fmt.Println("Terjadi kesalahan:", err)
				return
			}

			if success {
				fmt.Println("--Password berhasil diubah--")
			} else {
				fmt.Println("Gagal mengubah password.")
			}
		}
	}
	fmt.Println("=== Terimakasih  telah bertransaksi !!! ===")

}
