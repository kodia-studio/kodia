# Enterprise Frontend System

The Kodia Framework features a premium, enterprise-grade frontend system built on **Svelte 5**, **Tailwind 4**, and **Bits UI**. This document outlines the core components and architectural patterns used to build professional, high-performance web applications.

## Core Principles

1.  **Svelte 5 Runes**: Leveraging `$state`, `$derived`, and `$effect` for highly performant and reactive state management.
2.  **Tailwind 4 Design System**: A modernized design system focused on HSL variables for flexible theming and dark mode support.
3.  **Bits UI Foundation**: All interactive components (Selects, Modals, DatePickers) are built on top of the accessible Bits UI primitives.
4.  **TanStack Table Core**: Advanced data grids with sorting, filtering, and pagination using a custom Svelte 5 adapter.
5.  **Framework Independence**: All UI components and utilities are built internally with zero third-party component dependencies.

---

## Dark Mode & Theme System

The Kodia Framework includes a native, built-in dark mode system powered by **Svelte 5 runes**:

### `themeStore` (`$lib/stores/theme.svelte.ts`)

A singleton reactive store that manages theme state:

```typescript
import { themeStore, initTheme } from '$lib/stores/theme.svelte';

// Initialize on app startup (called in +layout.svelte)
initTheme();

// Reactive access to dark mode state
const isDark = $derived(themeStore.dark);

// Toggle theme
themeStore.toggle();
```

**Features:**
- Automatically detects system preference (`prefers-color-scheme`)
- Persists user selection to `localStorage` as `theme: 'dark' | 'light'`
- Applies `.dark` class to `<html>` element for Tailwind CSS dark mode
- Fully reactive with Svelte 5 `$state` and `$derived`
- No external dependencies (replaces deprecated `mode-watcher`)

---

## Layout System

Layouts are located in `$lib/components/layouts`. They provide the structural framework for different parts of the application.

-   **AppLayout**: The standard wrapper for public-facing or general app pages, featuring a fixed responsive `Navbar`.
-   **AdminLayout**: A sidebar-driven layout designed for dashboards and management consoles. Features a collapsible sidebar, search header, and notifications.
-   **AuthLayout**: A centered, glassmorphic container designed for login, registration, and password reset flows with premium background effects.

---

## Component Library

### Form System (`$lib/components/forms`)

Our form system is designed around **Bits UI** for accessibility primitives and seamlessly integrates with **SvelteKit Superforms** for server-side validation.

-   **KForm**: A wrapper component around Superforms for simplified form handling with error display.
-   **Form**: Base form component for building custom forms with validation integration.
-   **Input**: A standardized text input with support for icons and error states.
-   **Select**: An accessible, premium selection component built with Bits UI.
-   **Checkbox**: A fully styled, accessible checkbox.
-   **DatePicker**: An advanced date selection component using `@internationalized/date`.
-   **FileUpload**: A premium drag-and-drop file uploader with previews and progress visualizers.

### Data & Analytics (`$lib/components/data` & `charts`)

-   **DataTable**: A high-performance table component using TanStack Table Core with custom Svelte 5 adapter. Features sorting, pagination, and fully reactive state via `$state` runes.
    
-   **AreaChart** & **BarChart**: Built-in chart components using pure SVG and D3 scale/shape functions. Features include:
    - Interactive tooltips on hover
    - Responsive scaling
    - Dark mode support via Tailwind CSS
    - No external charting dependencies
    
    ```typescript
    <AreaChart
      data={[{ date: new Date(), value: 100 }, ...]}
      height={300}
      color="hsl(var(--primary))"
    />
    ```
    
-   **ChartCard**: A container for analytics featuring titles, subtitles, and action areas.

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
