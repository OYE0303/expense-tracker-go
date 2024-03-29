package router

import (
	"net/http"

	hd "github.com/OYE0303/expense-tracker-go/internal/handler"
	"github.com/OYE0303/expense-tracker-go/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// New initializes a new router and returns it
func New(handler *hd.Handler) http.Handler {
	r := mux.NewRouter()

	// user
	r.HandleFunc("/v1/user/signup", handler.User.Signup).Methods(http.MethodPost)
	r.HandleFunc("/v1/user/login", handler.User.Login).Methods(http.MethodPost)

	auth := alice.New(middleware.Authenticate)

	// icon
	r.Handle("/v1/icon", auth.ThenFunc(handler.Icon.List)).Methods(http.MethodGet)

	// main category
	r.Handle("/v1/main-category", auth.ThenFunc(handler.MainCateg.CreateMainCateg)).Methods(http.MethodPost)
	r.Handle("/v1/main-category", auth.ThenFunc(handler.MainCateg.GetAllMainCateg)).Methods(http.MethodGet)
	r.Handle("/v1/main-category/{id}", auth.ThenFunc(handler.MainCateg.UpdateMainCateg)).Methods(http.MethodPatch)
	r.Handle("/v1/main-category/{id}", auth.ThenFunc(handler.MainCateg.DeleteMainCateg)).Methods(http.MethodDelete)

	// sub category
	r.Handle("/v1/sub-category", auth.ThenFunc(handler.SubCateg.CreateSubCateg)).Methods(http.MethodPost)
	r.Handle("/v1/sub-category", auth.ThenFunc(handler.SubCateg.GetAllSubCateg)).Methods(http.MethodGet)
	r.Handle("/v1/main-category/{id}/sub-category", auth.ThenFunc(handler.SubCateg.GetByMainCategID)).Methods(http.MethodGet)
	r.Handle("/v1/sub-category/{id}", auth.ThenFunc(handler.SubCateg.UpdateSubCateg)).Methods(http.MethodPatch)
	r.Handle("/v1/sub-category/{id}", auth.ThenFunc(handler.SubCateg.DeleteSubCateg)).Methods(http.MethodDelete)

	// transaction
	r.Handle("/v1/transaction", auth.ThenFunc(handler.Transaction.Create)).Methods(http.MethodPost)
	r.Handle("/v1/transaction", auth.ThenFunc(handler.Transaction.GetAll)).Methods(http.MethodGet)
	r.Handle("/v1/transaction/info", auth.ThenFunc(handler.Transaction.GetAccInfo)).Methods(http.MethodGet)
	r.Handle("/v1/transaction/{id}", auth.ThenFunc(handler.Transaction.Delete)).Methods(http.MethodDelete)
	r.Handle("/v1/transaction/bar-chart", auth.ThenFunc(handler.Transaction.GetBarChartData)).Methods(http.MethodGet)
	r.Handle("/v1/transaction/pie-chart", auth.ThenFunc(handler.Transaction.GetPieChartData)).Methods(http.MethodGet)

	regular := alice.New(middleware.LogRequest, middleware.EnableCORS)

	return regular.Then(r)
}
