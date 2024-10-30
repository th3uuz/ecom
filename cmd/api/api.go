package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/th3uuz/ecom/service/cart"
	"github.com/th3uuz/ecom/service/order"
	"github.com/th3uuz/ecom/service/product"
	"github.com/th3uuz/ecom/service/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
	store *user.Store // Campo para armazenar o Store
}

func NewAPIServer(addr string, db *sql.DB, store *user.Store) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
		store: store, // Armazenando o Store
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(s.store)
	userHandler.RegisterRoutes(subrouter)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	orderStore := order.NewStore(s.db)

	cartHandler := cart.NewHandler(orderStore, productStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}