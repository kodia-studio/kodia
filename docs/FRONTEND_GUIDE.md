# Enterprise Frontend System

The Kodia Framework features a premium, enterprise-grade frontend system built on **Svelte 5**, **Tailwind 4**, and **Bits UI**. This document outlines the core components and architectural patterns used to build professional, high-performance web applications.

## Core Principles

1.  **Svelte 5 Runes**: Leveraging `$state`, `$derived`, and `$effect` for highly performant and reactive state management.
2.  **Tailwind 4 Design System**: A modernized design system focused on HSL variables for flexible theming and dark mode support.
3.  **Bits UI Foundation**: All interactive components (Selects, Modals, DatePickers) are built on top of the accessible Bits UI primitives.
4.  **TanStack Table**: Advanced data grids with sorting, filtering, and pagination out of the box.

---

## Layout System

Layouts are located in `$lib/components/layouts`. They provide the structural framework for different parts of the application.

-   **AppLayout**: The standard wrapper for public-facing or general app pages, featuring a fixed responsive `Navbar`.
-   **AdminLayout**: A sidebar-driven layout designed for dashboards and management consoles. Features a collapsible sidebar, search header, and notifications.
-   **AuthLayout**: A centered, glassmorphic container designed for login, registration, and password reset flows with premium background effects.

---

## Component Library

### Form System (`$lib/components/forms`)

Our form system is designed to work seamlessly with **SvelteKit Superforms** and **Zod** validation.

-   **FormField**: Wraps inputs with labels, error messages, and descriptions.
-   **Input**: A standardized text input with support for icons and error states.
-   **Select**: An accessible, premium selection component built with Bits UI.
-   **Checkbox**: A fully styled, accessible checkbox.
-   **DatePicker**: An advanced date selection component using `@internationalized/date`.
-   **FileUpload**: A premium drag-and-drop file uploader with previews and progress visualizers.

### Data & Analytics (`$lib/components/data` & `charts`)

-   **DataTable**: A high-performance table component using TanStack Table.
-   **ChartCard**: A container for analytics featuring titles, subtitles, and action areas.
-   **AreaChart**: A responsive area chart visualization built on **LayerChart** and **D3**.

---

## State Management

### ApiClient (`$lib/api/client.ts`)

The new `api` singleton provides a strictly typed interface for all backend communications:

-   **Global Loading/Error**: Global `isLoading` and `error` states are handled via Svelte 5 runes.
-   **Automatic Headers**: Automatically handles `Content-Type` and `Authorization` headers.

### Query Rune (`$lib/stores/query.svelte.ts`)

A reactive data-fetching primitive that handles:
-   Initial data fetching.
-   Reactive loading/error states.
-   Basic TTL-based caching (5 minutes by default).
-   Manual refetching.

```typescript
const users = createQuery(() => api.get<User[]>("/users"));
// Reactive usage:
// {#if users.isLoading}...{:else}{users.data}{/if}
```

---

## Best Practices

-   **Styling**: Use the HSL-based Tailwind variables (e.g., `text-primary`, `bg-muted`) instead of hardcoded colors.
-   **Transitions**: Use Svelte's built-in `fade` and `slide` transitions for a premium, non-static user experience.
-   **Accessibility**: Always provide descriptive labels and follow the patterns established in the Bits UI-based components.
