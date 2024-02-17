package transaction

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/internal/model/maincateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/subcateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/user"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
)

type TransactionFactory struct {
	transaction *testutil.Factory[Transaction]
	user        *testutil.Factory[user.User]
	maincateg   *testutil.Factory[maincateg.MainCateg]
	subcateg    *testutil.Factory[subcateg.SubCateg]
}

func NewTransactionFactory(db *sql.DB) *TransactionFactory {
	return &TransactionFactory{
		transaction: testutil.NewFactory(db, Transaction{}, BluePrint, Inserter),
		user:        testutil.NewFactory(db, user.User{}, user.Blueprint, user.Inserter),
		maincateg:   testutil.NewFactory(db, maincateg.MainCateg{}, maincateg.BluePrint, maincateg.Inserter),
		subcateg:    testutil.NewFactory(db, subcateg.SubCateg{}, subcateg.Blueprint, subcateg.Inserter),
	}
}

func (tf *TransactionFactory) PrepareUserMainAndSubCateg() (user.User, maincateg.MainCateg, subcateg.SubCateg, error) {
	u := user.User{}
	i := icon.Icon{}
	m, _, err := tf.maincateg.Build().WithOne(&u).WithOne(&i).InsertWithAss()
	if err != nil {
		return user.User{}, maincateg.MainCateg{}, subcateg.SubCateg{}, err
	}

	ow := subcateg.SubCateg{UserID: u.ID, MainCategID: m.ID}
	s, err := tf.subcateg.Build().Overwrite(ow).Insert()
	if err != nil {
		return user.User{}, maincateg.MainCateg{}, subcateg.SubCateg{}, err
	}

	return u, m, s, nil
}

func (tf *TransactionFactory) Reset() {
	tf.transaction.Reset()
	tf.user.Reset()
	tf.maincateg.Reset()
}
