package transaction_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/OYE0303/expense-tracker-go/internal/model/transaction"
	"github.com/OYE0303/expense-tracker-go/pkg/dockerutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/golang-migrate/migrate"
	"github.com/stretchr/testify/suite"
)

type TransactionSuite struct {
	suite.Suite
	db               *sql.DB
	migrate          *migrate.Migrate
	transactionModel *transaction.TransactionModel
	f                *transaction.TransactionFactory
}

func TestTransactionSuite(t *testing.T) {
	suite.Run(t, new(TransactionSuite))
}

func (s *TransactionSuite) SetupSuite() {
	port := dockerutil.RunDocker()
	db, migrate := testutil.ConnToDB(port)
	logger.Register()

	s.db = db
	s.migrate = migrate
	s.transactionModel = transaction.NewTransactionModel(s.db)
	s.f = transaction.NewTransactionFactory(db)
}

func (s *TransactionSuite) TearDownSuite() {
	s.db.Close()
	s.migrate.Close()
	dockerutil.PurgeDocker()
}

func (s *TransactionSuite) SetupTest() {
	s.transactionModel = transaction.NewTransactionModel(s.db)
	s.f = transaction.NewTransactionFactory(s.db)
}

func (s *TransactionSuite) TearDownTest() {
	tx, err := s.db.Begin()
	if err != nil {
		s.Require().NoError(err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM transactions"); err != nil {
		s.Require().NoError(err)
	}

	s.Require().NoError(tx.Commit())
	s.f.Reset()
}

func (s *TransactionSuite) TestCreate() {
	user, maincateg, subcateg, err := s.f.PrepareUserMainAndSubCateg()
	s.Require().NoError(err)

	fmt.Println("user: ", user)
	fmt.Println("maincateg: ", maincateg)
	fmt.Println("subcateg: ", subcateg)
}
