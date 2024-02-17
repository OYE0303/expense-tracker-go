package transaction

import (
	"database/sql"
	"fmt"
	"time"
)

func BluePrint(i int, last Transaction) Transaction {
	now := time.Now()
	return Transaction{
		Price: float64(i*10.0 + 1.0),
		Note:  "test" + fmt.Sprint(i),
		Date:  now,
	}
}

func Inserter(db *sql.DB, t Transaction) (Transaction, error) {
	stmt := `INSERT INTO transactions (user_id, main_category_id, sub_category_id, price, note, date) VALUES (?, ?, ?, ?, ?, ?)`

	res, err := db.Exec(stmt, t.UserID, t.MainCategID, t.SubCategID, t.Price, t.Note, t.Date)
	if err != nil {
		return Transaction{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Transaction{}, err
	}

	t.ID = id
	return t, nil
}
