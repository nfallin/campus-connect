# Campus Connect Backend
An API that manages creating, storing, and organizing events for the Campus Connect web application.

# Setup, building, running, and testing

## 1. Install Go

Use the official installer at https://go.dev/dl/.

to verify:
```powershell
go version
```

## 2. Download the project

Clone the repository:

```powershell
git clone https://github.com/nfallin/campus-connect-backend.git
cd campus-connect-backend
```
(All commands from now on are performed from the campus-connect-backend directory)

## 3. Download the Go dependencies

From the project folder, run:

```powershell
go mod download
```

## 4. Start the API

Run the application directly with:

```powershell
go run .
```

The API will start on port `8080`. Keep this terminal open while using the
API.

On the same computer, the API is available at:

```text
http://localhost:8080
```

To access it from another device on the same local network, replace
`localhost` with the computer's local IP address, for example:

```text
http://192.168.1.42:8080
```

## 5. Optional: build and run an executable

Instead of using `go run`, create an executable:

```powershell
go build -o campus-connect.exe .
```

Start the built application at http://localhost:8080 with:

```powershell
.\campus-connect.exe
```

## 6. Run unit tests
Run every unit test with
```powershell
go test ./...
```

Run a single package's unit tests with
```powershell
go test ./Model
```