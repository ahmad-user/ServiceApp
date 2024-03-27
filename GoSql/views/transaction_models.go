package views

import (
	"GoSql/users"
	"fmt"

	"gorm.io/gorm"
)

type Tbl_transfer struct {
	gorm.Model
	SenderHP   string
	ReceiverHP string
	Amount     float64
	Type       string
	Status     string
}

func TransferBalanceHp(senderHP, receiverHP string, amount float64, typeHP string, db *gorm.DB) error {
	var sender, receiver users.Tbl_user

	if err := db.Where("HP = ?", senderHP).First(&sender).Error; err != nil {
		return fmt.Errorf("failed to find sender: %v", err)
	}
	if err := db.Where("HP = ?", receiverHP).First(&receiver).Error; err != nil {
		return fmt.Errorf("failed to find receiver: %v", err)
	}
	if sender.Balance < amount {
		return fmt.Errorf("dana anda tidak mencukupi")
	}
	sender.Balance -= amount
	if err := db.Save(&sender).Error; err != nil {
		return fmt.Errorf("failed to update sender's balance: %v", err)
	}
	receiver.Balance += amount
	if err := db.Save(&receiver).Error; err != nil {
		return fmt.Errorf("failed to update receiver's balance: %v", err)
	}
	transfer := Tbl_transfer{
		SenderHP:   senderHP,
		ReceiverHP: receiverHP,
		Amount:     amount,
		Type:       typeHP,
		Status:     "Sukses",
	}
	if err := db.Create(&transfer).Error; err != nil {
		return fmt.Errorf("failed to save transfer history: %v", err)
	}

	return nil
}

// func TransferBalanceHP(receiverHP string, amount float64, typeHP string, db *gorm.DB) error {
// 	var sender, receiver users.Tbl_user

// 	if err := db.First(&sender).Error; err != nil {
// 		return fmt.Errorf("gagal menemukan pengirim: %v", err)
// 	}
// 	if err := db.Where("HP = ?", receiverHP).First(&receiver).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return fmt.Errorf("penerima dengan nomor telepon %s tidak ditemukan", receiverHP)
// 		}
// 		return fmt.Errorf("gagal menemukan penerima: %v", err)
// 	}
// 	if sender.Balance < amount {
// 		return fmt.Errorf("dana pengirim tidak mencukupi")
// 	}
// 	sender.Balance -= amount
// 	if err := db.Save(&sender).Error; err != nil {
// 		return fmt.Errorf("gagal mengupdate saldo pengirim: %v", err)
// 	}
// 	receiver.Balance += amount
// 	if err := db.Save(&receiver).Error; err != nil {
// 		return fmt.Errorf("gagal mengupdate saldo penerima: %v", err)
// 	}
// 	transfer := Tbl_transfer{
// 		ReceiverHP: receiverHP,
// 		Amount:     amount,
// 		Type:       typeHP,
// 		Date:       time.Now(),
// 	}
// 	if err := db.Create(&transfer).Error; err != nil {
// 		return fmt.Errorf("gagal menyimpan riwayat transfer: %v", err)
// 	}

// 	return nil
// }

func HistoryTransfer(senderHP uint, database *gorm.DB) ([]Tbl_transfer, error) {
	var history []Tbl_transfer
	if err := database.Where("ID = ?", senderHP).Find(&history).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil riwayat top up: %v", err)
	}

	return history, nil
}
