package transaction

import (
	"context"
	"errors"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/interfaces"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

const (
	PackageName = "usecase/transaction"
)

var (
	weekDayFormat = "Mon"
)

type TransactionUC struct {
	Transaction interfaces.TransactionModel
	MainCateg   interfaces.MainCategModel
	SubCateg    interfaces.SubCategModel
}

func NewTransactionUC(t interfaces.TransactionModel, m interfaces.MainCategModel, s interfaces.SubCategModel) *TransactionUC {
	return &TransactionUC{
		Transaction: t,
		MainCateg:   m,
		SubCateg:    s,
	}
}

func (t *TransactionUC) Create(ctx context.Context, trans domain.CreateTransactionInput) error {
	// check if the main category exists
	mainCateg, err := t.MainCateg.GetByID(trans.MainCategID, trans.UserID)
	if errors.Is(err, domain.ErrDataNotFound) {
		return domain.ErrDataNotFound
	}
	if err != nil {
		return err
	}

	// check if the type in main category matches the transaction type
	if trans.Type != mainCateg.Type {
		logger.Error("Create Transaction failed", "package", PackageName, "err", domain.ErrTypeNotConsistent)
		return domain.ErrTypeNotConsistent
	}

	// check if the sub category exists
	subCateg, err := t.SubCateg.GetByID(trans.SubCategID, trans.UserID)
	if errors.Is(err, domain.ErrDataNotFound) {
		return domain.ErrDataNotFound
	}
	if err != nil {
		return err
	}

	// check if the sub category matches the main category
	if subCateg.MainCategID != trans.MainCategID {
		logger.Error("Create Transaction failed", "package", PackageName, "err", domain.ErrMainCategNotConsistent)
		return domain.ErrMainCategNotConsistent
	}

	if err := t.Transaction.Create(ctx, trans); err != nil {
		return err
	}

	return nil
}

func (t *TransactionUC) GetAll(ctx context.Context, query domain.GetQuery, user domain.User) ([]domain.Transaction, error) {
	return t.Transaction.GetAll(ctx, query, user.ID)
}

func (t *TransactionUC) Update(ctx context.Context, trans domain.UpdateTransactionInput, user domain.User) error {
	// check if the main category exists
	mainCateg, err := t.MainCateg.GetByID(trans.MainCategID, user.ID)
	if err != nil {
		return err
	}

	// check if the type in main category matches the transaction type
	if trans.Type != mainCateg.Type {
		logger.Error("Update Transaction failed", "package", PackageName, "err", domain.ErrTypeNotConsistent)
		return domain.ErrTypeNotConsistent
	}

	// check if the sub category exists
	subCateg, err := t.SubCateg.GetByID(trans.SubCategID, user.ID)
	if err != nil {
		return err
	}

	// check if the sub category matches the main category
	if trans.MainCategID != subCateg.MainCategID {
		logger.Error("Update Transaction failed", "package", PackageName, "err", domain.ErrMainCategNotConsistent)
		return domain.ErrMainCategNotConsistent
	}

	// check permission
	if _, err := t.Transaction.GetByIDAndUserID(ctx, trans.ID, user.ID); err != nil {
		return err
	}

	return t.Transaction.Update(ctx, trans)
}

func (t *TransactionUC) Delete(ctx context.Context, id int64, user domain.User) error {
	// check permission
	if _, err := t.Transaction.GetByIDAndUserID(ctx, id, user.ID); err != nil {
		return err
	}

	return t.Transaction.Delete(ctx, id)
}

func (t *TransactionUC) GetAccInfo(ctx context.Context, query domain.GetAccInfoQuery, user domain.User) (domain.AccInfo, error) {
	return t.Transaction.GetAccInfo(ctx, query, user.ID)
}

func (t *TransactionUC) GetBarChartData(ctx context.Context, chartDateRange domain.ChartDateRange, transactionType domain.TransactionType, user domain.User) (domain.ChartData, error) {
	dataByWeekday, err := t.Transaction.GetBarChartData(ctx, chartDateRange, transactionType, user.ID)
	if err != nil {
		return domain.ChartData{}, err
	}

	start, end, err := cvtDateToTime(chartDateRange.StartDate, chartDateRange.EndDate)
	if err != nil {
		logger.Error("cvtDateToTime failed", "package", PackageName, "err", err)
		return domain.ChartData{}, err
	}

	var chartData domain.ChartData
	for t := start; t.Before(end) || t.Equal(end); t = t.AddDate(0, 0, 1) {
		weekDay := t.Format(weekDayFormat)

		chartData.Labels = append(chartData.Labels, weekDay)

		// if there is no data for the weekday, append 0
		if _, ok := dataByWeekday[weekDay]; !ok {
			chartData.Datasets = append(chartData.Datasets, 0)
		} else {
			chartData.Datasets = append(chartData.Datasets, dataByWeekday[weekDay])
		}
	}

	return chartData, nil
}

func (t *TransactionUC) GetPieChartData(ctx context.Context, chartDateRange domain.ChartDateRange, transactionType domain.TransactionType, user domain.User) (domain.ChartData, error) {
	return t.Transaction.GetPieChartData(ctx, chartDateRange, transactionType, user.ID)
}
