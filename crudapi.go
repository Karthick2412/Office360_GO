package main

import (
	"log"
	"net/http"
	"taskupdate/controllers"
	"taskupdate/middleware"

	"github.com/gorilla/mux"
)

func initialRouter() {
	rout := mux.NewRouter().StrictSlash(true)
	rout.Use(middleware.CommonMiddleware)
	rout.HandleFunc("/login", controllers.Login).Methods("POST")
	rout.HandleFunc("/", Testing).Methods("GET")
	rout.HandleFunc("/userCreate", controllers.Signup).Methods("POST")
	rout.HandleFunc("/forgotpassword", controllers.ForgotPassword).Methods("POST")
	rout.HandleFunc("/resetpassword", controllers.ResetPassword).Methods("POST")
	auth := rout.NewRoute().Subrouter()
	auth.Use(middleware.JwtVerify)

	auth.HandleFunc("/dailyTask", controllers.DailyTask).Methods("POST")
	auth.HandleFunc("/task", controllers.GetTasks).Methods("GET")
	auth.HandleFunc("/report/{StDate}/{EndDate}", controllers.Report).Methods("GET")
	auth.HandleFunc("/validate", controllers.Validate).Methods("GET")
	//r.Use(middleware.CommonMiddleware)
	log.Fatal(http.ListenAndServe(":8080", rout))
}

func Testing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Working fine")
	json.NewEncoder(w).Encode("Working Good")

}
