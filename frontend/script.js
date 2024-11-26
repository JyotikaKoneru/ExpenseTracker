const baseUrl = "http://localhost:8081/expenses";

function goHome() {
    document.getElementById("mainPage").style.display = "block";
    document.getElementById("addExpensePage").style.display = "none";
    document.getElementById("viewExpensesPage").style.display = "none";
    document.getElementById("filterExpensesPage").style.display = "none";
}

function showAddExpense() {
    document.getElementById("mainPage").style.display = "none";
    document.getElementById("addExpensePage").style.display = "block";
    // Clear the previous input fields when opening Add Expense page
    document.getElementById("expenseName").value = "";
    document.getElementById("expenseCategory").value = "";
    document.getElementById("expenseAmount").value = "";
}

function showViewExpenses() {
    document.getElementById("mainPage").style.display = "none";
    document.getElementById("viewExpensesPage").style.display = "block";
    fetch(baseUrl)
        .then(response => response.json())
        .then(data => {
            const expensesList = document.getElementById("expensesList");
            const totalAmount = document.getElementById("totalAmount");
            expensesList.innerHTML = "";
            let total = 0;
            if (data.length === 0) {
                expensesList.innerHTML = "<p>No expenses found!</p>";
            } else {
                data.forEach(expense => {
                    const expenseCard = document.createElement("div");
                    expenseCard.classList.add("expense-card");
                    expenseCard.innerHTML = `
                        <span>${expense.name} (${expense.category}) - $${expense.amount}</span>
                        <span>Date: ${expense.date}</span>
                        <button onclick="deleteExpense(${expense.id})">Delete</button>
                    `;
                    expensesList.appendChild(expenseCard);
                    total += expense.amount;
                });
            }
            totalAmount.innerText = `Total: $${total}`;
        });
}

function showFilterExpenses() {
    document.getElementById("mainPage").style.display = "none";
    document.getElementById("filterExpensesPage").style.display = "block";
}

function addExpense() {
    const name = document.getElementById("expenseName").value;
    const category = document.getElementById("expenseCategory").value;
    const amount = parseFloat(document.getElementById("expenseAmount").value);

    if (!name || !category || isNaN(amount)) {
        alert("Please fill all fields before adding an expense!");
        return;
    }

    // Send the expense data to the backend
    fetch(baseUrl, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, category, amount }),
    }).then(() => {
        alert("Expense added successfully!");
        goHome(); // Go back to the main page (if needed)
        showAddExpense(); // Show the Add Expense page again after adding an expense
    });
}

function deleteExpense(id) {
    fetch(`${baseUrl}/${id}`, { method: "DELETE" }).then(() => {
        showViewExpenses(); // Recalculate and refresh the list after deletion
    });
}

function filterExpenses() {
    const category = document.getElementById("filterCategory").value;

    if (!category) {
        alert("Please select a category to filter!");
        return;
    }

    fetch(`${baseUrl}?category=${category}`)
        .then(response => response.json())
        .then(data => {
            const filteredExpensesList = document.getElementById("filteredExpensesList");
            const filteredTotal = document.getElementById("filteredTotal");

            filteredExpensesList.innerHTML = "";
            let total = 0;

            // If no expenses found for the selected category
            if (!data || data.length === 0) {
                filteredExpensesList.innerHTML = `<p>No entries found in the "${category}" category.</p>`;
                filteredTotal.innerText = "";  // Clear the total display
            } else {
                data.forEach(expense => {
                    const expenseCard = document.createElement("div");
                    expenseCard.classList.add("expense-card");
                    expenseCard.innerHTML = `${expense.name} (${expense.category}) - $${expense.amount}`;
                    filteredExpensesList.appendChild(expenseCard);
                    total += expense.amount;
                });
                filteredTotal.innerText = `Total: $${total}`;
            }
        });
}
