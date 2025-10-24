# React + Go Todo App (TypeScript, React Query, MongoDB, Chakra UI)

![Demo App](https://i.ibb.co/JvRTWmW/Group-93.png)

Build a clean, responsive full‑stack Todo app with Go (Fiber) and React. It supports CRUD, filtering, sorting, validation, theming, and a production build served by Go.


## Highlights

- Tech Stack: Go (Fiber) · React 18 · TypeScript · MongoDB · TanStack Query · Chakra UI
- CRUD: Create, Read, Update, Delete todos
- Theming: Light/Dark mode with Chakra UI
- Responsive: Works great on mobile and desktop
- Data layer: React Query for fetching, caching, mutations, and revalidation
- New features:
	- Filtering: by status (all/active/completed) and by search text
	- Sorting: by created date, title, or status; order asc/desc
	- Validation: server- and client‑side checks (non‑empty, max length, duplicate prevention)
	- Metadata: todos store createdAt; PATCH supports toggling completed
	- CORS: enabled for local dev (http://localhost:5173)

## API Overview

Base URL: `/api`

Endpoints
- `GET /api/todos` — list todos with optional query params:
	- `status`: `all` | `active` | `completed` (default: `all`)
	- `sortBy`: `createdAt` | `body` | `completed` (default: `createdAt`)
	- `order`: `asc` | `desc` (default: `desc`)
	- `search`: string (case‑insensitive substring match on `body`)
- `POST /api/todos` — create todo
	- Body: `{ "body": string }`
	- Validation: trimmed, non‑empty, <= 200 chars, not duplicate (case‑insensitive)
- `PATCH /api/todos/:id` — update todo
	- Body (optional): `{ "completed": boolean }` (defaults to `true` if omitted)
- `DELETE /api/todos/:id` — delete todo

## Project Structure

```
./
	main.go                 # Fiber server + MongoDB + API + static serving in production
	client/                 # React app (Vite + TS + Chakra UI + React Query)
		src/components/       # TodoList, TodoItem, TodoForm, toolbar (filters + search)
```

## Prerequisites

- Go 1.20+ (tested with Fiber v2)
- Node.js 18+ and npm
- A MongoDB connection string (Atlas or local)

## Configuration (.env)

Create an `.env` file in the project root (same folder as `main.go`).

```dotenv
MONGODB_URI=<your_mongodb_uri>
PORT=5000
ENV=development
```

Notes
- `MONGODB_URI` name matches what the server reads in `main.go`.
- When `ENV=production`, Fiber serves the built React app from `client/dist`.

## Run Locally

Open two terminals.

1) Backend (run from the project root)

```powershell
cd D:\Semester-5\GO
go run main.go or npm run . 
```

2) Frontend (Vite dev server)

```powershell
cd D:\Semester-5\GO\client
npm install
npm run dev
```

App URLs
- React: http://localhost:5173
- API:   http://localhost:5000/api



## Troubleshooting

- Running `go run main.go` from `client/` fails with `GetFileAttributesEx main.go`: run the server from the project root or use `go run ..\main.go` from `client/`.
- Port already in use (5000): stop the previous process or kill it (Windows) `taskkill /f /pid <PID>`.
- CORS errors in dev: ensure the API runs on 5000 and the client on 5173; CORS is enabled for `http://localhost:5173`.
- Mongo connection: verify `MONGODB_URI` in `.env` and internet/firewall access to Atlas.

## App Screenshot
### The app link will be expiered on 25.10.2025  

![Screenshot from the app](/client/public/app.png)

https://courageous-simplicity-production-44c7.up.railway.app/