package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// Expense represents a single expense entry
type Expense struct {
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
}

// ExpensesData holds all the expenses
type ExpensesData struct {
	Expenses []Expense `json:"expenses"`
}

// LoadExpenses loads expenses from a file
func LoadExpenses(filename string) (ExpensesData, error) {
	var data ExpensesData

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil // Return empty if file does not exist
		}
		return data, err
	}

	err = json.Unmarshal(file, &data)
	return data, err
}

// SaveExpenses saves expenses to a file
func SaveExpenses(filename string, data ExpensesData) error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, file, 0644)
}

// AddExpense adds a new expense to the data
func AddExpense(data *ExpensesData, description, category string, amount float64) {
	expense := Expense{
		Description: description,
		Category:    category,
		Amount:      amount,
		Date:        time.Now(),
	}
	data.Expenses = append(data.Expenses, expense)
}

// ViewExpenses prints all expenses with optional filtering by category
func ViewExpenses(data ExpensesData, category string) {
	var total float64
	fmt.Println("\nExpenses:")
	for _, expense := range data.Expenses {
		if category == "" || expense.Category == category {
			fmt.Printf("Description: %s | Category: %s | Amount: %.2f | Date: %s\n",
				expense.Description, expense.Category, expense.Amount, expense.Date.Format("2006-01-02"))
			total += expense.Amount
		}
	}
	fmt.Printf("Total Spending: %.2f\n", total)
}

// Main menu for CLI
func main() {
	const filename = "expenses.json"

	data, err := LoadExpenses(filename)
	if err != nil {
		fmt.Println("Error loading expenses:", err)
		return
	}

	for {
		fmt.Println("\nExpense Tracker Menu:")
		fmt.Println("1. Add Expense")
		fmt.Println("2. View All Expenses")
		fmt.Println("3. View Expenses by Category")
		fmt.Println("4. Exit")
		fmt.Print("Enter your choice: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			var description, category string
			var amount float64

			fmt.Print("Enter description: ")
			fmt.Scan(&description)
			fmt.Print("Enter category (e.g., food, transportation): ")
			fmt.Scan(&category)
			fmt.Print("Enter amount: ")
			fmt.Scan(&amount)

			AddExpense(&data, description, category, amount)
			err := SaveExpenses(filename, data)
			if err != nil {
				fmt.Println("Error saving expenses:", err)
			} else {
				fmt.Println("Expense added successfully.")
			}

		case 2:
			ViewExpenses(data, "")

		case 3:
			var category string
			fmt.Print("Enter category to filter by: ")
			fmt.Scan(&category)
			ViewExpenses(data, category)

		case 4:
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
