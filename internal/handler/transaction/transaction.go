package transaction

import (
	"errors"
	"net/http"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/usecase/interfaces"
	"github.com/OYE0303/expense-tracker-go/pkg/ctxutil"
	"github.com/OYE0303/expense-tracker-go/pkg/errutil"
	"github.com/OYE0303/expense-tracker-go/pkg/jsonutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/validator"
)

const (
	packageName = "handler/transaction"
)

type TransactionHandler struct {
	transaction interfaces.TransactionUC
}

func NewTransactionHandler(t interfaces.TransactionUC) *TransactionHandler {
	return &TransactionHandler{
		transaction: t,
	}
}

func (t *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input createTransactionReq
	if err := jsonutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsonutil.ReadJSON failed", "package", packageName, "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	user := ctxutil.GetUser(r)
	trans := domain.CreateTransactionInput{
		UserID:      user.ID,
		Type:        domain.CvtToTransactionType(input.Type),
		MainCategID: input.MainCategID,
		SubCategID:  input.SubCategID,
		Price:       input.Price,
		Date:        input.Date,
		Note:        input.Note,
	}

	v := validator.New()
	if !v.CreateTransaction(trans) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	ctx := r.Context()
	if err := t.transaction.Create(ctx, trans); err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			errutil.BadRequestResponse(w, r, err)
			return
		}

		errutil.ServerErrorResponse(w, r, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusCreated, nil, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", packageName, "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (t *TransactionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	query, err := genGetAllQuery(r)
	if err != nil {
		errutil.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if !v.GetTransaction(query) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	user := ctxutil.GetUser(r)
	ctx := r.Context()
	transactions, err := t.transaction.GetAll(ctx, query, *user)
	if err != nil {
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	transResp := cvtToGetTransactionResp(transactions)
	respData := map[string]interface{}{
		"transactions": transResp.Transactions,
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, respData, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", packageName, "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (t *TransactionHandler) GetAccInfo(w http.ResponseWriter, r *http.Request) {
	query := genGetAccInfoQuery(r)
	v := validator.New()
	if !v.GetAccInfo(query) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	user := ctxutil.GetUser(r)
	ctx := r.Context()
	info, err := t.transaction.GetAccInfo(ctx, query, *user)
	if err != nil {
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	resp := map[string]interface{}{
		"total_income":  info.TotalIncome,
		"total_expense": info.TotalExpense,
		"total_balance": info.TotalBalance,
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, resp, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", packageName, "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (t *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := jsonutil.ReadID(r)
	if err != nil {
		errutil.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if !v.Delete(id) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	ctx := r.Context()
	user := ctxutil.GetUser(r)
	if err := t.transaction.Delete(ctx, id, *user); err != nil {
		if errors.Is(err, domain.ErrTransactionDataNotFound) {
			errutil.BadRequestResponse(w, r, err)
			return
		}

		errutil.ServerErrorResponse(w, r, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, nil, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", packageName, "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (t *TransactionHandler) GetBarChartData(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	rawTransactionType := r.URL.Query().Get("type")

	dateRange := domain.ChartDateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}
	transactionType := domain.CvtToTransactionType(rawTransactionType)

	v := validator.New()
	if !v.GetChartData(dateRange, transactionType) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	user := ctxutil.GetUser(r)
	ctx := r.Context()
	data, err := t.transaction.GetBarChartData(ctx, dateRange, transactionType, *user)
	if err != nil {
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	resp := map[string]interface{}{
		"chart_data": data,
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, resp, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", packageName, "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (t *TransactionHandler) GetPieChartData(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	rawTransactionType := r.URL.Query().Get("type")
	transactionType := domain.CvtToTransactionType(rawTransactionType)

	dateRange := domain.ChartDateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}

	v := validator.New()
	if !v.GetPieChartData(dateRange, transactionType) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	user := ctxutil.GetUser(r)
	data, err := t.transaction.GetPieChartData(r.Context(), dateRange, transactionType, *user)
	if err != nil {
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	resp := map[string]interface{}{
		"chart_data": data,
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, resp, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", packageName, "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}
