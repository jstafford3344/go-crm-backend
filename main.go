package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Customer struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Contacted bool   `json:"contacted"`
}

var customers = []Customer{
	{1, "Shawn Bradley", "GOAT", "shawnbradley@gmail.com", "215-610-7812", true},
	{2, "Tyrone Hill", "WOAT", "thill@gmail.com", "856-523-1357", true},
	{3, "John Thompson", "Coach", "jthompson@georgetown.edu", "856-229-7171", false},
	{4, "Allen Iverson", "Basketball Hooper", "ai@76ers.com", "215-281-7817", false},
	{5, "Dikembe Mutombo", "Shot Blocker", "dmutombo@76ers.com", "215-271-8289", false},
	{6, "Glenn Robinson", "PF", "bennett@76ers.com", "819-010-7287", true},
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customers)
	log.Println("Received GET customers request")
}

func getCustomerById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid customerId!", http.StatusBadRequest)
		log.Println("Invalid customerId!")
		return
	}

	for _, customer := range customers {
		if customer.Id == id {
			json.NewEncoder(w).Encode(customer)
			log.Println("Received GET customers request for id:", id)
			return
		}
	}

	http.Error(w, "Customer not found", http.StatusNotFound)
	log.Println("Customer not found!")
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newEntry Customer
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &newEntry)

	for _, customer := range customers {
		if customer.Email == newEntry.Email {
			http.Error(w, "Email already exists!", http.StatusConflict)
			log.Printf("Received POST request for email %s, which already exists!", newEntry.Email)
			return
		}
	}

	maxID := 0
	for _, customer := range customers {
		if customer.Id > maxID {
			maxID = customer.Id
		}
	}
	newEntry.Id = maxID + 1

	customers = append(customers, newEntry)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customers)
	log.Println("Received POST customers request. New id:", newEntry.Id)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid customerId!", http.StatusBadRequest)
		log.Println("Invalid customerId!")
		return
	}

	var updatedCustomer Customer
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Failed to read request body:", err)
		return
	}
	err = json.Unmarshal(reqBody, &updatedCustomer)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Println("Failed to unmarshal JSON:", err)
		return
	}

	index := -1
	for i, customer := range customers {
		if customer.Id == id {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(w, "Customer not found", http.StatusNotFound)
		log.Println("Customer not found!")
		return
	}

	customers[index].Name = updatedCustomer.Name
	customers[index].Role = updatedCustomer.Role
	customers[index].Email = updatedCustomer.Email
	customers[index].Phone = updatedCustomer.Phone
	customers[index].Contacted = updatedCustomer.Contacted

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customers[index])
	log.Println("Received PUT customers request for id:", id)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid customerId!", http.StatusBadRequest)
		log.Println("Invalid customerId!")
		return
	}

	index := -1
	for i, customer := range customers {
		if customer.Id == id {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(w, "Customer not found", http.StatusNotFound)
		log.Println("Customer not found!")
		return
	}

	customers = append(customers[:index], customers[index+1:]...)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customers)
	log.Println("Received DELETE customers request for id:", id)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/customers", getCustomers).Methods("GET")
	r.HandleFunc("/customers/{id}", getCustomerById).Methods("GET")
	r.HandleFunc("/customers", addCustomer).Methods("POST")
	r.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")
	r.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	log.Println("Starting server on port 8000...")
	http.ListenAndServe(":8000", r)
}
