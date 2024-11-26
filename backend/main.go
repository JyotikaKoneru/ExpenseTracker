package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Expense represents an expense entry
type Expense struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Date     string  `json:"date"`
}

var (
	expenses  = []Expense{}
	idCounter = 1
	mu        sync.Mutex
)

func main() {
	mux := http.NewServeMux()

	// Define routes
	mux.HandleFunc("/expenses", handleExpenses)
	mux.HandleFunc("/expenses/", handleDeleteExpense)

	// Enable CORS middleware
	handlerWithCORS := enableCORS(mux)

	// Start the server
	fmt.Println("Server running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", handlerWithCORS))
}

// handleExpenses handles GET and POST requests for expenses
func handleExpenses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Filter by category if query parameter exists
		category := r.URL.Query().Get("category")
		if category != "" {
			filteredExpenses := filterExpensesByCategory(category)
			writeJSON(w, filteredExpenses)
			return
		}

		// Return all expenses
		writeJSON(w, expenses)

	case http.MethodPost:
		// Parse the expense from the request body
		var expense Expense
		if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Add a date to the expense (current date in US order)
		expense.Date = time.Now().Format("2006/02/01")

		// Add the expense to the slice
		mu.Lock()
		expense.ID = idCounter
		idCounter++
		expenses = append(expenses, expense)
		mu.Unlock()

		writeJSON(w, map[string]string{"message": "Expense added successfully"})
	}
}

// handleDeleteExpense handles DELETE requests to delete an expense
func handleDeleteExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the expense ID from the URL
	idStr := strings.TrimPrefix(r.URL.Path, "/expenses/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	// Delete the expense with the given ID
	mu.Lock()
	defer mu.Unlock()
	for i, expense := range expenses {
		if expense.ID == id {
			expenses = append(expenses[:i], expenses[i+1:]...)
			writeJSON(w, map[string]string{"message": "Expense deleted successfully"})
			return
		}
	}

	http.Error(w, "Expense not found", http.StatusNotFound)
}

// filterExpensesByCategory filters expenses by category
func filterExpensesByCategory(category string) []Expense {
	mu.Lock()
	defer mu.Unlock()
	var filtered []Expense
	for _, expense := range expenses {
		if strings.EqualFold(expense.Category, category) {
			filtered = append(filtered, expense)
		}
	}
	return filtered
}

// enableCORS is a middleware to handle CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
