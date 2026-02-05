package repositories

import (
	"database/sql"
	"kasir-api/models"
	"fmt"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	var (
		res *models.Transaction
	)

	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail,0)

	for _, item := range items {
		var productName string
		var productID, productPrice, stock int

		err := tx.QueryRow("SELECT id, name, price, stock FROM products WHERE id=$1", item.ProductID).Scan(&productID, &productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Product %d not found", item.ProductID)
		}

		if err != nil {
			return nil, err
		}

		subtotal := item.Quantity * productPrice
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, productID)
		
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID: productID,
			ProductName: productName,
			Quantity: item.Quantity,
			Subtotal: subtotal,
		})


	}
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING ID", totalAmount).Scan(&transactionID)
	
	if err != nil {
		return nil, err
	}

	// for i, detail := range details {
	// 	details[i].TransactionID = transactionID
	// 	_, err := tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)", transactionID, detail.ProductID, detail.Quantity, detail.Subtotal)
		
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	
	stmt, err := tx.Prepare(`
		INSERT INTO transaction_details 
		(transaction_id, product_id, quantity, subtotal) 
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	for _, d := range details {
		d.TransactionID = transactionID

		res, err := stmt.Exec(
			d.TransactionID,
			d.ProductID,
			d.Quantity,
			d.Subtotal,
		)
		if err != nil {
			return nil, err
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}
		if rows == 0 {
			return nil, fmt.Errorf("failed to insert transaction detail for product_id %d", d.ProductID)
		}
	}
	
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	
	res = &models.Transaction{
		ID: transactionID,
		TotalAmount: totalAmount,
		Details: details,
	}

	return res,nil
}
