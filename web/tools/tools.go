package tools

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// IsAPI Tell whether request requires API reply or not.
func IsAPI(req *http.Request) bool {
	return strings.Contains(req.URL.String(), "/api/")
}

// GetInt parse request to get int value.
func GetInt(req *http.Request, key string) (int, error) {
	vars := mux.Vars(req)
	value, err := strconv.Atoi(vars[key])
	if err != nil {
		log.Printf("Web: requested key: %s , not found in: %s", key, req.URL)
		return 0, errors.New("Invalid key requested")
	}
	return value, nil
}

// GetString parse request to get int value.
func GetString(req *http.Request, key string) (value string, found bool) {
	vars := mux.Vars(req)
	value, found = vars[key]
	return
}

//Fail fails current request with API and redirect with Web; should set session error ;) but can't right now ...
func Fail(w http.ResponseWriter, req *http.Request, err string, backRoute string) {
	log.Printf("Web: Failed access to %s due to %s", req.URL.String(), err)
	if IsAPI(req) {
		GenerateAPIError(w, err)
	} else {
		Redirect(w, req, backRoute)
	}
}

//Redirect user to targeted page.
func Redirect(w http.ResponseWriter, req *http.Request, route string) {
	log.Printf("Web: Redirecting to %s", route)
	http.Redirect(w, req, route, http.StatusSeeOther)
}

// HasValue tell whether value is present or not.
func HasValue(req *http.Request, key string) bool {
	vars := mux.Vars(req)
	_, ok := vars[key]
	return ok
}

// GenerateAPIError generate a simple JSON reply with error message provided.
func GenerateAPIError(w http.ResponseWriter, message string) {
	var repm = make(map[string]string)
	repm["status"] = "error"
	repm["error"] = message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(repm)
}

// GenerateAPIOkAndSend generate a simple JSON reply with status: ok.
func GenerateAPIOkAndSend(w http.ResponseWriter) {
	var repm = make(map[string]string)
	repm["status"] = "ok"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(repm)
}

// GenerateAPIOk generate a simple JSON reply with status: ok.
func GenerateAPIOk(w http.ResponseWriter) map[string]string {
	var repm = make(map[string]string)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	repm["status"] = "ok"
	return repm
}