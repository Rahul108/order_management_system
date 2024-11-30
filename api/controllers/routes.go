package controllers

import "github.com/rahul108/order_management_system/api/middlewares"

func (s *Server) initializeRoutes() {

	apiV1 := s.Router.PathPrefix("/api/v1").Subrouter()

	// Login Route
	apiV1.HandleFunc("/signup", middlewares.SetMiddlewareJSON(s.Signup)).Methods("POST")
	apiV1.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")
	apiV1.HandleFunc("/refresh", middlewares.SetMiddlewareJSON(s.RefreshToken)).Methods("POST")

	// order management
	apiV1.HandleFunc("/orders", middlewares.AuthMiddleware(s.CreateOrder)).Methods("POST")
	apiV1.HandleFunc("/orders", middlewares.AuthMiddleware(s.ListOrders)).Methods("GET")
	apiV1.HandleFunc("/orders/{CONSIGNMENT_ID}/cancel", middlewares.AuthMiddleware(s.CancelOrder)).Methods("PUT")
}
