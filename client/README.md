# Kaiyo AI - Travel Planning Assistant

A modern travel planning application built with React, TypeScript, and Tailwind CSS.

## Features

- **Landing Page**: Choose between trip planning and community features
- **Authentication**: JWT-based login/signup system
- **Chat Interface**: AI-powered travel planning assistant
- **Interactive Maps**: React Leaflet integration for location visualization
- **Community Forum**: Connect with fellow travelers
- **Comic-style UI**: Unique border design inspired by Redis documentation

## Tech Stack

- **Frontend**: React 18 + TypeScript
- **Styling**: Tailwind CSS + Custom Comic Borders
- **Routing**: React Router DOM
- **State Management**: TanStack Query
- **Maps**: React Leaflet
- **Icons**: Lucide React
- **Build Tool**: Vite

## Getting Started

1. Install dependencies:

```bash
npm install
```

2. Copy environment variables:

```bash
cp .env.example .env
```

3. Start the development server:

```bash
npm run dev
```

4. Open [http://localhost:5174](http://localhost:5174) in your browser

## Project Structure

```
src/
├── components/
│   ├── ui/           # Reusable UI components
│   └── chat/         # Chat-specific components
├── pages/            # Page components
├── lib/              # Utilities and helpers
└── index.css         # Global styles and comic borders
```

## Features Overview

### Landing Page

- Two main options: "Plan Your Trip" and "Talk to Community"
- Redirects to authentication if user is not logged in

### Chat Interface

- Three-panel layout: Sidebar, Chat Area, Map Panel
- Real-time chat with AI assistant
- Interactive map showing travel locations
- Responsive design with mobile-friendly sidebar

### Community Forum

- Browse travel posts and experiences
- Search functionality
- Popular topics and travel statistics
- Social features (likes, comments, sharing)

### Authentication

- Login/Signup forms
- JWT token management
- Protected routes with automatic redirects

## Styling

The app uses a unique "comic-style" design with:

- Bold black borders (3px)
- Box shadows for depth
- Rounded corners
- Purple primary color scheme
- Hover animations and transitions
