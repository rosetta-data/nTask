package API

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	globalStructs "github.com/r4ulcl/NetTask/globalStructs"
	"github.com/r4ulcl/NetTask/manager/database"
	"github.com/r4ulcl/NetTask/manager/utils"
)

// @Summary Handle callback from slave
// @Description Handle callback from slave
// @Tags callback
// @Accept json
// @Produce json
// @Param Authorization header string true "OAuth Key" default(WLJ2xVQZ5TXVw4qEznZDnmEEV)
// @Success 200 "OK"
// @Failure 400 {string} string "Invalid callback body"
// @Failure 401 {string} string "Unauthorized"
// @Router /callback [post]
func HandleCallback(w http.ResponseWriter, r *http.Request, config *utils.ManagerConfig, db *sql.DB) {
	oauthKey := r.Header.Get("Authorization")
	if incorrectOauth(oauthKey, config.OAuthToken) && incorrectOauthWorker(oauthKey, config.OauthTokenWorkers) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var result globalStructs.Task
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, "Invalid callback body", http.StatusBadRequest)
		return
	}

	fmt.Println(result)

	fmt.Printf("Received result (ID: %s) from :\n %s with output: %s\n", result.ID, result.WorkerName, result.Output)

	// Update task with the worker one
	database.UpdateTask(db, result)

	// Set worker to iddle now
	database.SetWorkerworkingToString(false, db, result.WorkerName)

	// Handle the result as needed

	w.WriteHeader(http.StatusOK)
}

// @Summary Get workers
// @Description Handle worker request
// @Tags worker
// @Accept json
// @Produce json
// @Param Authorization header string true "OAuth Key" default(WLJ2xVQZ5TXVw4qEznZDnmEEV)
// @Success 200 {string} string "OK"
// @Router /worker [get]
func HandleWorkerGet(w http.ResponseWriter, r *http.Request, config *utils.ManagerConfig, db *sql.DB) {
	oauthKey := r.Header.Get("Authorization")
	if incorrectOauth(oauthKey, config.OAuthToken) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//get workers
	workers, err := database.GetWorkers(db)
	if err != nil {
		http.Error(w, "Invalid callback body", http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(workers)
	if err != nil {
		http.Error(w, "Invalid callback body", http.StatusBadRequest)
		return
	}

	// Print the JSON data
	//fmt.Println(string(jsonData))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(jsonData))
}

// @Summary Add a worker
// @Description Add a worker, normally done by the worker
// @Tags worker
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Router /worker [post]
// @Type basic
// @In header
// @Name Authorization
// @Param Authorization header string true "OAuth Key" default(WLJ2xVQZ5TXVw4qEznZDnmEEV)
// @Param worker body globalStructs.Worker true "Worker object to create"
func HandleWorkerPost(w http.ResponseWriter, r *http.Request, config *utils.ManagerConfig, db *sql.DB) {
	oauthKey := r.Header.Get("Authorization")
	if incorrectOauth(oauthKey, config.OAuthToken) && incorrectOauthWorker(oauthKey, config.OauthTokenWorkers) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request globalStructs.Worker
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid callback body", http.StatusBadRequest)
		return
	}

	request.IP = ReadUserIP(r)

	fmt.Println(request.Name, request.IP, request.Name)

	err = database.AddWorker(db, &request)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 { // MySQL error number for duplicate entry
				database.SetWorkerUPto(true, db, &request)
				database.SetWorkerCount(0, db, &request)
				return
			} else {
				message := "Invalid worker info: " + err.Error()
				http.Error(w, message, http.StatusBadRequest)
				return
			}
		}

	}

	// Handle the result as needed
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Worker with Name %s added", request.Name)
}

// @Summary Remove a worker
// @Description Remove a worker from the system
// @Tags worker
// @Accept json
// @Produce json
// @Param Authorization header string true "OAuth Key" default(WLJ2xVQZ5TXVw4qEznZDnmEEV)
// @Success 200 {array} string
// @Router /worker/{NAME} [delete]
// @Param NAME path string false "Worker NAME"
func HandleWorkerDeleteName(w http.ResponseWriter, r *http.Request, config *utils.ManagerConfig, db *sql.DB) {
	oauthKey := r.Header.Get("Authorization")
	if incorrectOauth(oauthKey, config.OAuthToken) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	name := vars["NAME"]

	err := database.RmWorkerName(db, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "")
}

// @Summary Get status of worker
// @Description Get status of worker
// @Tags worker
// @Accept json
// @Produce json
// @Param Authorization header string true "OAuth Key" default(WLJ2xVQZ5TXVw4qEznZDnmEEV)
// @Success 200 {array} globalStructs.Worker
// @Router /worker/{NAME} [get]
// @Param NAME path string false "Worker NAME"
func HandleWorkerStatus(w http.ResponseWriter, r *http.Request, config *utils.ManagerConfig, db *sql.DB) {
	oauthKey := r.Header.Get("Authorization")
	if incorrectOauth(oauthKey, config.OAuthToken) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	name := vars["NAME"]

	fmt.Println("NAME " + name)

	worker, err := database.GetWorker(db, name)
	if err != nil {
		http.Error(w, "Invalid callback body"+err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(worker)
	if err != nil {
		http.Error(w, "Invalid callback body"+err.Error(), http.StatusBadRequest)
		return
	}

	// Print the JSON data
	//fmt.Println(string(jsonData))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(jsonData))
}

// Other functions
//ReadUserIP read user IP from request
func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	// Split IP address and port
	ip, _, err := net.SplitHostPort(IPAddress)
	if err == nil {
		return ip
	}

	// If there's an error (e.g., no port found), return the original address
	return IPAddress
}
