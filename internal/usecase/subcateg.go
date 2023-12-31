package usecase

import (
	"errors"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

type subCategUC struct {
	SubCateg  SubCategModel
	MainCateg MainCategModel
}

func newSubCategUC(s SubCategModel, m MainCategModel) *subCategUC {
	return &subCategUC{
		SubCateg:  s,
		MainCateg: m,
	}
}

func (s *subCategUC) Create(categ *domain.SubCateg, userID int64) error {
	// check if the main category exists
	_, err := s.MainCateg.GetByID(categ.MainCategID, userID)
	if errors.Is(err, domain.ErrDataNotFound) {
		return domain.ErrDataNotFound
	}
	if err != nil {
		logger.Error("s.MainCateg.GetByID failed", "package", "usecase", "err", err)
		return err
	}

	// check if the sub category name is already taken
	categByUserID, err := s.SubCateg.GetOne(categ, userID)
	if err != nil && !errors.Is(err, domain.ErrDataNotFound) {
		logger.Error("s.SubCateg.GetOneByUserID failed", "package", "usecase", "err", err)
		return err
	}
	if categByUserID != nil {
		return domain.ErrDataAlreadyExists
	}

	if err := s.SubCateg.Create(categ, userID); err != nil {
		logger.Error("s.SubCateg.Create failed", "package", "usecase", "err", err)
		return err
	}

	return nil
}

func (s *subCategUC) GetAll(userID int64) ([]*domain.SubCateg, error) {
	categs, err := s.SubCateg.GetAll(userID)
	if err != nil {
		logger.Error("s.SubCateg.GetAll failed", "package", "usecase", "err", err)
		return nil, err
	}

	return categs, nil
}

func (s *subCategUC) GetByMainCategID(userID, mainCategID int64) ([]*domain.SubCateg, error) {
	categs, err := s.SubCateg.GetByMainCategID(userID, mainCategID)
	if err != nil {
		logger.Error("s.SubCateg.GetByMainCategID failed", "package", "usecase", "err", err)
		return nil, err
	}

	return categs, nil
}

func (s *subCategUC) Update(categ *domain.SubCateg, userID int64) error {
	// check if the sub category exists
	categByID, err := s.SubCateg.GetByID(categ.ID, userID)
	if errors.Is(err, domain.ErrDataNotFound) {
		return domain.ErrDataNotFound
	}
	if err != nil {
		logger.Error("s.SubCateg.GetByID failed", "package", "usecase", "err", err)
		return err
	}

	/*
		Have to add mainCategID after getting the sub category by ID because
		1. We don't expect the front end to pass the main category ID
		2. GetOne() needs the main category ID to check if the sub category name is already taken
	*/
	categ.MainCategID = categByID.MainCategID

	// check if the sub category name is already taken
	categByUserID, err := s.SubCateg.GetOne(categ, userID)
	if err != nil && !errors.Is(err, domain.ErrDataNotFound) {
		logger.Error("s.SubCateg.GetOneByUserID failed", "package", "usecase", "err", err)
		return err
	}
	if categByUserID != nil {
		return domain.ErrDataAlreadyExists
	}

	if err := s.SubCateg.Update(categ); err != nil {
		logger.Error("s.SubCateg.Update failed", "package", "usecase", "err", err)
		return err
	}

	return nil
}

func (s *subCategUC) Delete(id int64) error {
	if err := s.SubCateg.Delete(id); err != nil {
		logger.Error("s.SubCateg.Delete failed", "package", "usecase", "err", err)
		return err
	}

	return nil
}
