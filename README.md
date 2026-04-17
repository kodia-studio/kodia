# Kodia Framework

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Svelte Version](https://img.shields.io/badge/Svelte-5-FF3E00?style=flat&logo=svelte)](https://svelte.dev/)
[![Tailwind Version](https://img.shields.io/badge/Tailwind-4-38B2AC?style=flat&logo=tailwind-css)](https://tailwindcss.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

**Kodia** adalah framework fullstack profesional yang dirancang untuk kecepatan pengembangan, keamanan, dan pengalaman pengembang (DX) yang luar biasa. Dibangun dengan kekuatan **Golang Gin** di backend dan reaktivitas **SvelteKit** di frontend.

---

## ✨ Fitur Utama

- 🐨 **Kodia CLI**: Tool baris perintah yang kuat untuk scaffolding fitur secara instan.
- 🏗️ **Clean Architecture**: Backend yang modular, mudah diuji, dan skalabel dengan prinsip Clean Architecture.
- ⚡ **Modern Frontend**: Menggunakan Svelte 5 runes, Tailwind CSS v4, dan Bits UI untuk UI yang premium dan responsif.
- 🔐 **Secure by Default**: Autentikasi JWT (Access & Refresh Tokens), Middleware CORS, dan Validation out-of-the-box.
- 🗄️ **Multi-DB Support**: Mendukung PostgreSQL dan MySQL melalui GORM.
- 🐳 **Docker Ready**: Siap untuk deployment dengan Docker Compose yang sudah terkonfigurasi.

---

## 🚀 Quick Start

### 1. Prasyarat
- Go 1.26+
- Node.js 25+
- Docker & Docker Compose

### 2. Instalasi
Clone repository dan jalankan inisialisasi:
```bash
git clone https://github.com/kodia/framework-kodia.git my-app
cd my-app
# Setup env files
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

### 3. Jalankan menggunakan Docker
```bash
make docker-up
```

### 4. Mulai Development
```bash
make dev
```

---

## 🛠️ Struktur Proyek

```text
.
├── backend/          # Golang Gin API
├── frontend/         # SvelteKit + Tailwind v4
├── kodia-cli/        # Tool CLI untuk scaffolding
├── docker-compose.yml
└── Makefile
```

---

## 📚 Dokumentasi

Silakan baca dokumentasi lengkap kami (Segera hadir) untuk mempelajari lebih lanjut tentang routing, middleware, scaffolding, dan deployment.

---

## 📄 Lisensi

Kodia Framework dilisensikan di bawah [Lisensi MIT](LICENSE).
