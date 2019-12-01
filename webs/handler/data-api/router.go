package data_api

func (d *DataApiWeb) RegisterHandler() {
	dataRouter := d.Router.PathPrefix("/data").Subrouter()
	dataRouter.HandleFunc("/queryApi", d.handlerDataGet).Methods("POST")
	dataRouter.HandleFunc("/params", d.handlerDataSourceParams).Methods("GET")
	dataRouter.HandleFunc("/verify", d.handlerDataGet).Methods("POST")
	dataRouter.Use(d.addStartTimeMiddleware)
	dataRouter.Use(d.addSwiftNumberMiddleware)
}