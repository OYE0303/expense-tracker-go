package maincateg

import (
	"context"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/interfaces"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

type MainCategUC struct {
	MainCateg interfaces.MainCategModel
	Icon      interfaces.IconModel
}

func NewMainCategUC(m interfaces.MainCategModel, i interfaces.IconModel) *MainCategUC {
	return &MainCategUC{
		MainCateg: m,
		Icon:      i,
	}
}

func (m *MainCategUC) Create(categ domain.MainCateg, userID int64) error {
	// check if the icon exists
	if _, err := m.Icon.GetByID(categ.Icon.ID); err != nil {
		return err
	}

	return m.MainCateg.Create(&categ, userID)
}

func (m *MainCategUC) GetAll(ctx context.Context, userID int64, transType domain.TransactionType) ([]domain.MainCateg, error) {
	return m.MainCateg.GetAll(ctx, userID, transType)
}

func (m *MainCategUC) Update(categ *domain.MainCateg, userID int64) error {
	// check if the main category exists
	if _, err := m.MainCateg.GetByID(categ.ID, userID); err != nil {
		return err
	}

	// check if the icon exists
	if _, err := m.Icon.GetByID(categ.Icon.ID); err != nil {
		return err
	}

	if err := m.MainCateg.Update(categ); err != nil {
		return err
	}

	return nil
}

func (m *MainCategUC) Delete(id int64) error {
	if err := m.MainCateg.Delete(id); err != nil {
		logger.Error("m.MainCateg.Delete failed", "package", "usecase", "err", err)
		return err
	}

	return nil
}
