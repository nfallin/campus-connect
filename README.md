# Campus Connect

Campus Connect is a web application designed to help Santa Clara University students discover campus events in one place. The goal is to make it easier for students to find events, explore campus activities, and get involved without having to search through multiple platforms.

This repository contains the React frontend and Go API/Database for Campus Connect. The backend must be running for the frontend to work properly.

## 1. Prerequisites
- Node.js 22 or newer (`node -v` to check)
    - Official Node installer: https://nodejs.org/en/download
- Go version 1.26.0 or newer (`go version` to check)
    - Official Go installer: https://go.dev/dl/


## 2. Install Backend Dependencies
From the root directory of the repository, navigate to the backend directory:
```bash
cd backend
```
Then run:
```bash
go mod download
```

## 3. Database Management
If you've previously ran a version of the backend software that contained outdated database table formats, delete the previous database file before running the backend to re-initialize it with the new format.
```bash
rm campus_connect.db
```
You can verify difference in database structures by looking for modifications to `Tables.go`. 

## 4. Launch the Backend
Open a dedicated terminal for the backend and run:
```bash
go run .
```
This will start the API on port `8080`. Keep this terminal open while using the API.

## 5. Running Backend Unit Tests
To run every backend unit test:
```bash
go test ./...
```

To run the unit tests for a specific package:
```bash
go test ./Model
```
Replace 'Model' with the package name of your choice. 

## 6. Install Frontend Dependencies
If you're still in the backend directory, navigate to the frontend directory using:
```bash
cd ../frontend
```

If you're in the root directory, navigate using:
```bash
cd frontend
```

## 7. Launch the Frontend
Open a dedicated terminal for the frontend and run:
```bash
npm run dev
```
Vite prints a local URL (by default it will be `http://localhost:5173`). Open it in a web browser.

## 8. Frontend Unit Tests
To run the frontend unit tests:
```bash
npm test
```

## CI
- `.github/workflows/backend-ci.yml` runs `go test ./...` and `go build ./...` on every push and pull request.
- `.github/workflows/frontend-ci.yml` runs `npm test`, `npm run lint`, and `npm run build` on every push and pull request.

## Information Directory
- `README.md`: Instructions for installing, running, and testing the project
- `API.md`: Developer communication contract for sending information between the backend and frontend
- `REQUIREMENTS.md` Current list of system user requirements
- `SPECIFICATIONS.md` Current list of system user specifications
- `MVP ER Diagram.png` MVP ER Diagram
- `MVP CRC Diagram.png` MVP CRC Diagram

## Tech Stack
- **React and TypeScript:** Frontend development
- **Vite:** Development server and build tool
- **Go:** Backend API
- **SQLite** Backend Database
- **Vitest:** Frontend unit testing
- **GitHub Actions:** Automated testing and build checks