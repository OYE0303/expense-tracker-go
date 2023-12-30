package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/ctxutil"
	"github.com/OYE0303/expense-tracker-go/pkg/errutil"
	"github.com/OYE0303/expense-tracker-go/pkg/jsonutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/validator"
)

type transactionHandler struct {
	transaction TransactionUC
}

func newTransactionHandler(t TransactionUC) *transactionHandler {
	return &transactionHandler{
		transaction: t,
	}
}

func (t *transactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Type        string     `json:"type"`
		MainCategID int64      `json:"main_category_id"`
		SubCategID  int64      `json:"sub_category_id"`
		Price       int64      `json:"price"`
		Date        *time.Time `json:"date"`
		Note        string     `json:"note"`
	}
	if err := jsonutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsonutil.ReadJSON failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	transaction := domain.Transaction{
		Type:        input.Type,
		MainCategID: input.MainCategID,
		SubCategID:  input.SubCategID,
		Price:       input.Price,
		Date:        input.Date,
		Note:        input.Note,
	}

	v := validator.New()
	if !v.CreateTransaction(&transaction) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	ctx := r.Context()
	user := ctxutil.GetUser(r)
	if err := t.transaction.Create(ctx, user, &transaction); err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			errutil.BadRequestResponse(w, r, err)
			return
		}

		logger.Error("t.transaction.Create failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusCreated, nil, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

}
