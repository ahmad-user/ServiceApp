package views

import (
	"GoSql/users"
	"fmt"

	"gorm.io/gorm"
)

type Tbl_topUp struct {
	gorm.Model
	JumlahTopUp float64
	Type        string
	Status      string
}

func HistoryTopUp(userID uint, database *gorm.DB) ([]Tbl_topUp, error) {
	var history []Tbl_topUp
	if err := database.Where("ID = ?", userID).Find(&history).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil riwayat top up: %v", err)
	}

	return history, nil
}

func TopUp(userID uint, jumlahTopUp float64, typeTopUp string, database *gorm.DB) error {
	topUp := Tbl_topUp{
		JumlahTopUp: jumlahTopUp,
		Type:        typeTopUp,
		Status:      "Success",
	}
	if err := database.Create(&topUp).Error; err != nil {
		return fmt.Errorf("gagal membuat entri top up: %v", err)
	}

	if err := database.Model(&users.Tbl_user{}).Where("id = ?", userID).Update("balance", gorm.Expr("balance + ?", jumlahTopUp)).Error; err != nil {
		return fmt.Errorf("gagal memperbarui saldo pengguna: %v", err)
	}

	return nil
}
