package users

import (
	"fmt"

	"gorm.io/gorm"
)

type Tbl_user struct {
	gorm.Model
	HP       string
	Nama     string
	Password string
	Alamat   string
	Balance  float64
}

func (u *Tbl_user) GantiPassword(connection *gorm.DB, newPassword string) (bool, error) {
	query := connection.Model(u).Where("hp = ?", u.HP).Update("password", newPassword)
	if err := query.Error; err != nil {
		return false, err
	}

	return query.RowsAffected > 0, nil
}

func Register(connection *gorm.DB, newUser Tbl_user) (bool, error) {
	query := connection.Create(&newUser)
	if err := query.Error; err != nil {
		return false, err
	}

	return query.RowsAffected > 0, nil
}

func Login(connection *gorm.DB, hp string, password string) (Tbl_user, error) {
	var result Tbl_user
	err := connection.Where("hp = ? AND password = ?", hp, password).First(&result).Error
	if err != nil {
		return Tbl_user{}, err
	}
	return result, nil
}

func CariUsers(connection *gorm.DB, query string) ([]Tbl_user, error) {
	var daftarUsers []Tbl_user
	if err := connection.Where("HP LIKE ?", "%"+query+"%").Find(&daftarUsers).Error; err != nil {
		return nil, err
	}
	return daftarUsers, nil
}
func TambahUsers(connection *gorm.DB, tambah Tbl_user) (bool, error) {
	query := connection.Create(&tambah)
	if err := query.Error; err != nil {
		return false, err
	}
	return query.RowsAffected > 0, nil
}

func LihatUsers(connection *gorm.DB, userID uint) ([]Tbl_user, error) {
	var daftarBarang []Tbl_user
	if err := connection.Where("ID = ?", userID).Find(&daftarBarang).Error; err != nil {
		return nil, err
	}
	return daftarBarang, nil
}

func EditUsers(connection *gorm.DB, id uint, UpdateData map[string]interface{}) (bool, error) {
	var barang Tbl_user
	if err := connection.First(&barang, id).Error; err != nil {
		return false, err
	}
	if err := connection.Model(&barang).Updates(UpdateData).Error; err != nil {
		return false, err
	}
	return true, nil
}
func TopUp(ID uint, balance float64, database *gorm.DB) error {
	// Lakukan proses top up
	var user Tbl_user // Gantilah User dengan struktur data pengguna Anda
	if err := database.Model(&user).Where("id = ?", ID).Update("balance", gorm.Expr("balance + ?", balance)).Error; err != nil {
		return fmt.Errorf("gagal melakukan top up: %v", err)
	}
	fmt.Printf("Top up berhasil dilakukan. Saldo akun telah ditambahkan sebesar %.2f.\n", balance)
	return nil
}
func HistoryTopUp(ID int, db *gorm.DB) ([]Tbl_user, error) {
	var history []Tbl_user
	if err := db.Where("ID = ? ", ID).Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}
func HapusUsers(connection *gorm.DB, id uint) (bool, error) {
	var users Tbl_user
	if err := connection.First(&users, id).Error; err != nil {
		return false, err
	}
	if err := connection.Delete(&users).Error; err != nil {
		return false, err
	}
	return true, nil
}
