# Kodia Framework

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Svelte Version](https://img.shields.io/badge/Svelte-5-FF3E00?style=flat&logo=svelte)](https://svelte.dev/)
[![Tailwind Version](https://img.shields.io/badge/Tailwind-4-38B2AC?style=flat&logo=tailwind-css)](https://tailwindcss.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

**Kodia** is a professional fullstack framework designed for development speed, security, and exceptional developer experience (DX). Built with the power of **Golang Gin** on the backend and the reactivity of **SvelteKit** on the frontend.

---

## ✨ Key Features

- 🐨 **Kodia CLI**: A powerful command-line tool for instant feature scaffolding.
- 🏗️ **Clean Architecture**: A modular, highly testable, and scalable backend following Clean Architecture principles.
- ⚡ **Modern Frontend**: Powered by Svelte 5 runes, Tailwind CSS v4, and Bits UI for a premium, responsive UI.
- 🔐 **Secure by Default**: JWT Authentication (Access & Refresh Tokens), CORS Middleware, and out-of-the-box Validation.
- 🗄️ **Multi-DB Support**: Supports PostgreSQL and MySQL via GORM.
- 🐳 **Docker Ready**: Deployment-ready with pre-configured Docker Compose.

---

## 🚀 Quick Start

### 1. Prerequisites
- Go 1.26+
- Node.js 25+
- Docker & Docker Compose

### 2. Installation
Clone the repository and run initialization:
```bash
git clone https://github.com/kodia/framework-kodia.git my-app
cd my-app
# Setup env files
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

### 3. Run with Docker
```bash
make docker-up
```

### 4. Start Development
```bash
make dev
```

---

## 🛠️ Project Structure

```text
.
├── backend/          # Golang Gin API
├── frontend/         # SvelteKit + Tailwind v4
├── kodia-cli/        # CLI tool for scaffolding
├── docker-compose.yml
└── Makefile
```

---

## 📚 Documentation

Please read our full documentation (Coming soon) to learn more about routing, middleware, scaffolding, and deployment.

---

## 📄 License

Kodia Framework is licensed under the [MIT License](LICENSE).
