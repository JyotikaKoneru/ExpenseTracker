# ExpenseTracker
1. Unzip the files. Make sure the folders are as follows:
expense_tracker
    -backend
        -main.go
    -frontend
        -index.html
        -script.js
2. Open a new terminal and run the commands:
    > cd expense_tracker/backend
    > go clean -modcache
    > go mod init expense_track
    > go mod tidy
    > go run main.go

    Allow the app to run and open index.html in the browser. In the terminal, one should see "Server running on http://localhost:8081".
    In the browser, the application opens and the home page is displayed.
3. For stopping the server, use ctrl+c on windows.
4. For all other runs later on, simply use "go run main.go"

