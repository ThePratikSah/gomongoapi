package router

import (
	"github.com/ThePratikSah/gomongoapi/controller"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/calculate", controller.GetAggregationData).Methods("GET")
	router.HandleFunc("/calculate-without-aggr", controller.GetDataWithAgg).Methods("GET")

	return router
}
