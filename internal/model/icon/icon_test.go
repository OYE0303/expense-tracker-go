package icon

import (
	"database/sql"
	"testing"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/usecase/interfaces"
	"github.com/OYE0303/expense-tracker-go/pkg/dockerutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil/efactory"
	"github.com/golang-migrate/migrate"
	"github.com/stretchr/testify/suite"
)

type IconSuite struct {
	suite.Suite
	db      *sql.DB
	migrate *migrate.Migrate
	model   interfaces.IconModel
	// f       *testutil.Factory[Icon]
	f *efactory.Factory[Icon]
}

func TestIconSuite(t *testing.T) {
	suite.Run(t, new(IconSuite))
}

func (s *IconSuite) SetupSuite() {
	port := dockerutil.RunDocker()
	db, migrate := testutil.ConnToDB(port)
	logger.Register()
	s.model = NewIconModel(db)
	s.db = db
	s.migrate = migrate
	s.f = efactory.New(Icon{}).SetConfig(efactory.Config[Icon]{
		DB:        db,
		TableName: "icons",
	})
}

func (s *IconSuite) TearDownSuite() {
	s.db.Close()
	s.migrate.Close()
	dockerutil.PurgeDocker()
}

func (s *IconSuite) SetupTest() {
	s.model = NewIconModel(s.db)
	s.f = efactory.New(Icon{}).SetConfig(efactory.Config[Icon]{
		DB:        s.db,
		TableName: "icons",
	})
}

func (s *IconSuite) TearDownTest() {
	tx, err := s.db.Begin()
	s.Require().NoError(err)
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM icons")
	s.Require().NoError(err)

	s.Require().NoError(tx.Commit())
	s.f.Reset()
}

func (s *IconSuite) TestList() {
	for scenario, fn := range map[string]func(s *IconSuite, desc string){
		"when has icons, return all":    list_WithIcons_ReturnAll,
		"when has no icons, return nil": list_WithoutIcons_ReturnNil,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func list_WithIcons_ReturnAll(s *IconSuite, desc string) {
	icons, err := s.f.BuildList(2).Insert()
	s.Require().NoError(err, desc)

	expRes := []domain.Icon{
		{
			ID:  icons[0].ID,
			URL: icons[0].URL,
		},
		{
			ID:  icons[1].ID,
			URL: icons[1].URL,
		},
	}

	res, err := s.model.List()
	s.Require().NoError(err, desc)
	s.Require().Equal(expRes, res, desc)
}

func list_WithoutIcons_ReturnNil(s *IconSuite, desc string) {
	res, err := s.model.List()
	s.Require().NoError(err, desc)
	s.Require().Nil(res, desc)
}

func (s *IconSuite) TestGetByID() {
	for scenario, fn := range map[string]func(s *IconSuite, desc string){
		"when has icon, return icon":   getByID_WithIcon_ReturnIcon,
		"when has no icon, return err": getByID_WithoutIcon_ReturnErr,
	} {
		s.Run(testutil.GetFunName(fn), func() {
			s.SetupTest()
			fn(s, scenario)
			s.TearDownTest()
		})
	}
}

func getByID_WithIcon_ReturnIcon(s *IconSuite, desc string) {
	icons, err := s.f.BuildList(2).Insert()
	s.Require().NoError(err, desc)

	expRes := domain.Icon{
		ID:  icons[0].ID,
		URL: icons[0].URL,
	}

	res, err := s.model.GetByID(icons[0].ID)
	s.Require().NoError(err, desc)
	s.Require().Equal(expRes, res, desc)
}

func getByID_WithoutIcon_ReturnErr(s *IconSuite, desc string) {
	_, err := s.f.BuildList(2).Insert()
	s.Require().NoError(err, desc)

	expRes := domain.Icon{}

	res, err := s.model.GetByID(999)
	s.Require().Equal(expRes, res, desc)
	s.Require().Equal(domain.ErrIconNotFound, err, desc)
}
