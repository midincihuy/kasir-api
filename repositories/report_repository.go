package repositories

import (
	"database/sql"
	"kasir-api/models"
	"errors"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetReportToday() ([]models.Report, error) {
	query := "SELECT COALESCE(SUM(total_amount),0) as total_revenue, COUNT(id) as total_transaksi FROM transactions WHERE created_at::date = CURRENT_DATE"
	
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	query2 := "SELECT p.name, COALESCE(SUM(quantity),0) as qty_terjual FROM transactions t JOIN transaction_details td ON td.transaction_id = t.id JOIN products p ON p.id = td.product_id WHERE t.created_at::date = CURRENT_DATE GROUP BY p.name ORDER BY SUM(quantity) DESC LIMIT 1"
	
	rows2, err := repo.db.Query(query2)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()


	var pd models.ProdukTerlaris
	for rows2.Next() {
		err := rows2.Scan(&pd.Nama, &pd.QtyTerjual)
		if err != nil {
			return nil, err
		}
	}
	
	reports := make([]models.Report, 0)
	for rows.Next() {
		var p models.Report
		err := rows.Scan(&p.TotalRevenue, &p.TotalTransaksi)
		if err != nil {
			return nil, err
		}
		p.ProdukTerlaris = pd
		reports = append(reports, p)
	}

	return reports, nil
}

func (repo *ReportRepository) GetReport(start_date string, end_date string) ([]models.Report, error) {
	query := "SELECT COALESCE(SUM(total_amount),0) as total_revenue, COUNT(id) as total_transaksi FROM transactions t "

	query2 := "SELECT p.name, COALESCE(SUM(quantity),0) as qty_terjual FROM transactions t JOIN transaction_details td ON td.transaction_id = t.id JOIN products p ON p.id = td.product_id "
	args := []interface{}{}
	if start_date != "" {
		query += " WHERE t.created_at::date >= $1"
		query2 += " WHERE t.created_at::date >= $1"
		args = append(args, "%"+start_date+"%")
		if end_date != "" {
			query += " AND t.created_at::date <= $2"
			query2 += " AND t.created_at::date <= $2"
			args = append(args, "%"+end_date+"%")
		}
	}else if end_date != "" {
		query += " WHERE t.created_at::date <= $1"
		query2 += " WHERE t.created_at::date <= $1"
		args = append(args, "%"+end_date+"%")
	}
	
	
	rows, err := repo.db.Query(query, args...)
	if err == sql.ErrNoRows {
		return nil, errors.New("Data tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	query2 += " GROUP BY p.name ORDER BY SUM(quantity) DESC LIMIT 1"
	rows2, err := repo.db.Query(query2, args...)
	if err == sql.ErrNoRows {
		return nil, errors.New("Data tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	defer rows2.Close()


	var pd models.ProdukTerlaris
	for rows2.Next() {
		err := rows2.Scan(&pd.Nama, &pd.QtyTerjual)
		if err != nil {
			return nil, err
		}
	}
	
	reports := make([]models.Report, 0)
	for rows.Next() {
		var p models.Report
		err := rows.Scan(&p.TotalRevenue, &p.TotalTransaksi)
		if err != nil {
			return nil, err
		}
		p.ProdukTerlaris = pd
		reports = append(reports, p)
	}

	return reports, nil
}


