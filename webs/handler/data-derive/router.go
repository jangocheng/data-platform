package data_derive

func (d *DataDeriveWeb) RegisterHandler() {
	dataRouter := d.Router.PathPrefix("/data").Subrouter()
	dataRouter.HandleFunc("/queryDerive", d.handlerDataGet).Methods("POST")
	dataRouter.HandleFunc("/queryDeriveSet", d.handlerDataGet).Methods("POST")
	dataRouter.HandleFunc("/params", d.handlerDataApiParams).Methods("GET")
	dataRouter.HandleFunc("/verifyDerive", d.handlerDataGet).Methods("POST")
	dataRouter.HandleFunc("/verifyDeriveSet", d.handlerDataGet).Methods("POST")
	dataRouter.Use(d.addStartTimeMiddleware)
	dataRouter.Use(d.addSwiftNumberMiddleware)
}