package transaction

import (
	"database/sql"

	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/internal/model/maincateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/subcateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/user"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil/efactory"
)

type TransactionFactory struct {
	transaction *efactory.Factory[Transaction]
	user        *efactory.Factory[user.User]
	maincateg   *efactory.Factory[maincateg.MainCateg]
	subcateg    *efactory.Factory[subcateg.SubCateg]
}

func NewTransactionFactory(db *sql.DB) *TransactionFactory {
	return &TransactionFactory{
		transaction: efactory.New(Transaction{}).SetConfig(efactory.Config[Transaction]{
			DB:        &efactory.SQLDB{DB: db},
			TableName: "transactions",
			BluePrint: BluePrint,
		}),
		user: efactory.New(user.User{}).SetConfig(efactory.Config[user.User]{
			DB:        &efactory.SQLDB{DB: db},
			TableName: "users",
		}),
		maincateg: efactory.New(maincateg.MainCateg{}).SetConfig(efactory.Config[maincateg.MainCateg]{
			DB:        &efactory.SQLDB{DB: db},
			TableName: "main_categories",
			BluePrint: maincateg.BluePrint,
		}),
		subcateg: efactory.New(subcateg.SubCateg{}).SetConfig(efactory.Config[subcateg.SubCateg]{
			DB:        &efactory.SQLDB{DB: db},
			TableName: "sub_categories",
		}),
	}
}

func (tf *TransactionFactory) PrepareUserMainAndSubCateg() (user.User, maincateg.MainCateg, subcateg.SubCateg, icon.Icon, error) {
	u := user.User{}
	i := icon.Icon{}
	m, _, err := tf.maincateg.Build().WithOne(&u).WithOne(&i).InsertWithAss()
	if err != nil {
		return user.User{}, maincateg.MainCateg{}, subcateg.SubCateg{}, icon.Icon{}, err
	}

	ow := subcateg.SubCateg{UserID: u.ID, MainCategID: m.ID}
	s, err := tf.subcateg.Build().Overwrite(ow).Insert()
	if err != nil {
		return user.User{}, maincateg.MainCateg{}, subcateg.SubCateg{}, icon.Icon{}, err
	}

	return u, m, s, i, nil
}

func (tf *TransactionFactory) InsertTransactionsWithOneUser(i int, ow ...Transaction) ([]Transaction, user.User, []maincateg.MainCateg, []subcateg.SubCateg, []icon.Icon, error) {
	u := user.User{}

	iconPtrList := make([]interface{}, 0, i)
	for k := 0; k < i; k++ {
		iconPtrList = append(iconPtrList, &icon.Icon{})
	}

	maincategList, _, err := tf.maincateg.BuildList(i).WithOne(&u).WithMany(iconPtrList...).InsertWithAss()
	if err != nil {
		return nil, user.User{}, []maincateg.MainCateg{}, []subcateg.SubCateg{}, []icon.Icon{}, err
	}

	owSub := []subcateg.SubCateg{}
	for _, m := range maincategList {
		owSub = append(owSub, subcateg.SubCateg{UserID: m.UserID, MainCategID: m.ID})
	}

	subcategList, err := tf.subcateg.BuildList(i).Overwrites(owSub...).Insert()
	if err != nil {
		return nil, user.User{}, []maincateg.MainCateg{}, []subcateg.SubCateg{}, []icon.Icon{}, err
	}

	owTrans := []Transaction{}
	for k, m := range maincategList {
		owTrans = append(owTrans, Transaction{
			UserID:      m.UserID,
			MainCategID: m.ID,
			SubCategID:  subcategList[k].ID,
		})
	}

	transList, err := tf.transaction.BuildList(i).Overwrites(owTrans...).Overwrites(ow...).Insert()
	if err != nil {
		return nil, user.User{}, []maincateg.MainCateg{}, []subcateg.SubCateg{}, []icon.Icon{}, err
	}

	iconList := make([]icon.Icon, 0, i)
	for _, v := range iconPtrList {
		iconList = append(iconList, *v.(*icon.Icon))
	}

	return transList, u, maincategList, subcategList, iconList, nil
}

func (tf *TransactionFactory) Reset() {
	tf.transaction.Reset()
	tf.user.Reset()
	tf.maincateg.Reset()
	tf.subcateg.Reset()
}
