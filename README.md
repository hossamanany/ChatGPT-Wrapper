# ChatGPT Wrapper

A full-stack ChatGPT wrapper application with a Vue.js frontend and Go backend, designed to provide a seamless chat experience with OpenAI's ChatGPT models.

Created by Hossam Elanany

## Features

- Modern Vue.js 3 frontend with TypeScript support
- Go backend with Gin framework
- Real-time chat streaming capabilities
- Content filtering and moderation
- Markdown support for messages
- Code highlighting
- Local storage for chat history
- Responsive design with Tailwind CSS
- Direct integration with OpenAI's ChatGPT API

## Project Structure

```
ai-chat/
├── frontend/          # Vue.js frontend application
│   ├── src/           # Source files
│   ├── public/        # Static assets
│   └── package.json   # Frontend dependencies
└── backend/           # Go backend server
    ├── handlers/      # Request handlers
    ├── config/        # Configuration files
    ├── filters/       # Content filtering
    ├── openai/        # OpenAI integration
    └── models/        # Data models
```

## Prerequisites

- Node.js (v16 or higher)
- Go (v1.16 or higher)
- OpenAI API key with access to ChatGPT models

## Setup

### Frontend

```sh
cd frontend
npm install
```

### Backend

```sh
cd backend
go mod download
```

## Development

### Frontend Development Server

```sh
cd frontend
npm run dev
```

### Backend Server

```sh
cd backend
go run main.go
```

The frontend will be available at `http://localhost:5173` and the backend at `http://localhost:8080`.

## Building for Production

### Frontend

```sh
cd frontend
npm run build
```

### Backend

```sh
cd backend
go build
```

## Environment Variables

Create a `.env` file in the backend directory with the following variables:

```
OPENAI_API_KEY=your_api_key_here
PORT=8080
```

## API Endpoints

- `POST /api/chat` - Standard ChatGPT completion
- `POST /api/chat/stream` - Streaming ChatGPT completion

## Technologies Used

- Frontend:

  - Vue.js 3
  - TypeScript
  - Tailwind CSS
  - Vite
  - Pinia
  - Markdown-it
  - Highlight.js

- Backend:
  - Go
  - Gin
  - OpenAI API (ChatGPT)
  - Godotenv
