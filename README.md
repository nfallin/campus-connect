# Campus Connect

Campus Connect is a web application designed to help Santa Clara University students discover campus events in one place. The goal is to make it easier for students to find events, explore campus activities, and get involved without having to search through multiple platforms.

This repository contains the React frontend and Go API for Campus Connect. The backend must be running for the frontend to work properly.

## 1. Prerequisites
- Node.js 22 or newer (`node -v` to check)
    - Official Node installer: https://nodejs.org/en/download
- Go version 1.26.0 or newer (`go version` to check)
    - Official Go installer: https://go.dev/dl/

## 2. Install Dependencies
to install backend dependencies:
```bash
cd backend
go mod download
```
To install frontend dependencies:
```bash
cd frontend
npm install
```

## 3. Launch the backend
Open a dedicated terminal for the backend and run:
```bash
cd backend
go run .
```
This will start the API port `8080`. Keep this terminal open while using the
API.

## 4. Launch the frontend
Open a dedicated terminal for the frontend and run:
```bash
cd frontend
npm run dev
```
Vite prints a local URL (usually `http://localhost:5173`). Open it in a web browser.

## Running unit tests

To run the backend unit tests:
```bash
cd backend
go test ./...
```
To run a single backend package's unit tests:
```powershell
go test ./Model
```
To run the frontend unit tests:
```bash
cd frontend
npm test
```

## CI
- `.github/workflows/backend-ci.yml` runs `go test ./...` and `go build ./...` on every push and pull request.
- `.github/workflows/frontend-ci.yml` runs `npm test`, `npm run lint`, and `npm run build` on every push and pull request.

## Information Directory
- `README.md`: Instructions for installing, running, and testing the project
- `API.md`: Developer communication contract for sending information between the backend and frontend
- `requirements.md` Current list of system user requirements
- `specifications.md` Current list of system user specifications

## Tech Stack
- **React and TypeScript:** Frontend development
- **Vite:** Development server and build tool
- **Go:** Backend API
- **Vitest:** Frontend unit testing
- **GitHub Actions:** Automated testing and build checks