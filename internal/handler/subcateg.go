package handler

import (
	"net/http"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/pkg/ctxutil"
	"github.com/OYE0303/expense-tracker-go/pkg/errutil"
	"github.com/OYE0303/expense-tracker-go/pkg/jsonutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/OYE0303/expense-tracker-go/pkg/validator"
)

type subCategHandler struct {
	SubCateg SubCategUC
}

func newSubCategHandler(s SubCategUC) *subCategHandler {
	return &subCategHandler{
		SubCateg: s,
	}
}

func (s *subCategHandler) CreateSubCateg(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		MainCategID int64  `json:"main_category_id"`
	}
	if err := jsonutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsonutil.ReadJSON failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	categ := domain.SubCateg{
		Name:        input.Name,
		MainCategID: input.MainCategID,
	}

	v := validator.New()
	if !v.CreateSubCateg(&categ) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	user := ctxutil.GetUser(r)
	if err := s.SubCateg.Create(&categ, user.ID); err != nil {
		if err == domain.ErrDataAlreadyExists || err == domain.ErrDataNotFound {
			errutil.BadRequestResponse(w, r, err)
			return
		}

		logger.Error("s.SubCateg.Create failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (s *subCategHandler) GetAllSubCateg(w http.ResponseWriter, r *http.Request) {
	user := ctxutil.GetUser(r)
	categs, err := s.SubCateg.GetAll(user.ID)
	if err != nil {
		logger.Error("s.SubCateg.GetAll failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	respData := map[string]interface{}{
		"sub_categories": categs,
	}
	if err := jsonutil.WriteJSON(w, http.StatusOK, respData, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (s *subCategHandler) GetByMainCategID(w http.ResponseWriter, r *http.Request) {
	id, err := jsonutil.ReadID(r)
	if err != nil {
		logger.Error("jsonutil.ReadID failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	user := ctxutil.GetUser(r)
	categs, err := s.SubCateg.GetByMainCategID(user.ID, id)
	if err != nil {
		logger.Error("s.SubCateg.GetByMainCategID failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	respData := map[string]interface{}{
		"sub_categories": categs,
	}
	if err := jsonutil.WriteJSON(w, http.StatusOK, respData, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (s *subCategHandler) UpdateSubCateg(w http.ResponseWriter, r *http.Request) {
	id, err := jsonutil.ReadID(r)
	if err != nil {
		logger.Error("jsonutil.ReadID failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	var input struct {
		Name string `json:"name"`
	}
	if err := jsonutil.ReadJson(w, r, &input); err != nil {
		logger.Error("jsonutil.ReadJSON failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	categ := domain.SubCateg{
		ID:   id,
		Name: input.Name,
	}

	v := validator.New()
	if !v.UpdateSubCateg(&categ) {
		errutil.VildateErrorResponse(w, r, v.Error)
		return
	}

	user := ctxutil.GetUser(r)
	if err := s.SubCateg.Update(&categ, user.ID); err != nil {
		if err == domain.ErrDataNotFound || err == domain.ErrDataAlreadyExists {
			errutil.BadRequestResponse(w, r, err)
			return
		}

		logger.Error("s.SubCateg.Update failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}

func (s *subCategHandler) DeleteSubCateg(w http.ResponseWriter, r *http.Request) {
	id, err := jsonutil.ReadID(r)
	if err != nil {
		logger.Error("jsonutil.ReadID failed", "package", "handler", "err", err)
		errutil.BadRequestResponse(w, r, err)
		return
	}

	if err := s.SubCateg.Delete(id); err != nil {
		logger.Error("s.SubCateg.Delete failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, nil, nil); err != nil {
		logger.Error("jsonutil.WriteJSON failed", "package", "handler", "err", err)
		errutil.ServerErrorResponse(w, r, err)
		return
	}
}
