# Enterprise Frontend System

The Kodia Framework features a premium, enterprise-grade frontend system built on **Svelte 5**, **Tailwind 4**, and **Bits UI**. This document outlines the core components and architectural patterns used to build professional, high-performance web applications.

## Core Principles

1.  **Svelte 5 Runes**: Leveraging `$state`, `$derived`, and `$effect` for highly performant and reactive state management.
2.  **Tailwind 4 Design System**: A modernized design system focused on HSL variables for flexible theming and dark mode support.
3.  **Bits UI Foundation**: All interactive components (Selects, Modals, DatePickers) are built on top of the accessible Bits UI primitives.
4.  **TanStack Table Core**: Advanced data grids with sorting, filtering, and pagination using a custom Svelte 5 adapter.
5.  **Framework Independence**: All UI components and utilities are built internally with zero third-party component dependencies.
6.  **Zero Third-Party Dependencies**: Custom implementations for forms, validation, notifications, and error handling — no svelte-sonner, no sveltekit-superforms.

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

### Form System (`$lib/components/forms` & `$lib/forms`)

A **custom, zero-dependency form validation system** that works seamlessly with SvelteKit's `use:enhance` actions.

#### Form Factory (`createForm`)

The `createForm` function creates reactive forms with built-in validation, dirty tracking, and error handling:

```typescript
import { createForm } from '$lib/forms/createForm.svelte';

const { values, errors, touched, isSubmitting, isDirty, setFieldValue, validate, handleSubmit } = 
  createForm({
    initialValues: { email: '', password: '' },
    validators: {
      email: (val) => !val ? 'Email required' : !val.includes('@') ? 'Invalid email' : undefined,
      password: (val) => !val ? 'Password required' : val.length < 8 ? 'Min 8 chars' : undefined,
    },
  });

const onSubmit = handleSubmit(async (values) => {
  const response = await fetch('/api/login', {
    method: 'POST',
    body: JSON.stringify(values),
  });
  return response.json();
});
```

**Features:**
- Reactive state with Svelte 5 `$state` runes
- Field-level validation with custom validators
- Dirty tracking (detects if form has changed)
- Form-level error handling via `setErrors()`
- Compatible with SvelteKit's `use:enhance` for progressive enhancement
- No external dependencies (replaces sveltekit-superforms)

#### Form Components

-   **Form**: Base form component for building custom forms with validation integration.
-   **Input**: A standardized text input with error display and validation state.
-   **Button**: Reusable button component with variant/size control and disabled state.
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

### Toast Notification System (`$lib/stores/toast.svelte.ts`)

A custom, dependency-free notification system with FIFO eviction and auto-dismiss:

```typescript
import { toast } from '$lib/stores/toast.svelte';

// Show success notification
toast.success('Profile updated!');

// Show error notification
toast.error('Failed to save changes');

// Show info notification
toast.info('Processing your request...');

// Show warning notification
toast.warning('This action cannot be undone');

// Custom duration (default 4000ms)
toast.success('Quick message', { duration: 2000 });
```

**Features:**
- Max 3 toasts visible at once (older ones evicted)
- Auto-dismiss after configurable duration
- FIFO queue management
- Colored borders by type (success=green, error=red, info=blue, warning=yellow)
- Slide-in animation from bottom-right
- Responsive layout

### Error Boundary (`$lib/stores/error.store.svelte.ts`)

Global error state management for API errors:

```typescript
import { errorStore } from '$lib/stores/error.store.svelte';

// Errors automatically captured from:
// - window.onerror (global JS errors)
// - window.onunhandledrejection (unhandled promises)
// - API responses with 500+ status codes

// Display error in component:
{#if $errorStore}
  <ErrorBoundary error={$errorStore} />
{/if}

// Clear error
errorStore.clear();
```

**Features:**
- Captures global JavaScript errors
- Captures unhandled promise rejections
- Captures API server errors (500+)
- Displays error card with request ID for debugging
- Type-safe with ApiError interface

### ApiClient (`$lib/api/client.svelte.ts`)

The `api` singleton provides a strictly typed interface for all backend communications:

```typescript
import { api } from '$lib/api/client';

// GET request
const users = await api.get<User[]>('/users');

// POST request
const newUser = await api.post<User>('/users', { name: 'John' });

// PUT request
await api.put<User>(`/users/${id}`, { name: 'Jane' });

// PATCH request
await api.patch<User>(`/users/${id}`, { status: 'active' });

// DELETE request
await api.delete(`/users/${id}`);
```

**Features:**
- Automatic `Content-Type` header management
- Automatic `Authorization` header (Bearer token from auth store)
- Automatic error handling (500+ → errorStore)
- Form data support (multipart/form-data for file uploads)
- Environment variable configuration via `PUBLIC_API_URL`
- TypeScript generic response types

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

## Loading & Skeleton Components

### Skeleton Loaders (`$lib/components/ui/Skeleton.svelte`)

Flexible skeleton placeholder components for loading states:

```svelte
<!-- Simple text skeleton -->
<Skeleton variant="text" count={3} />

<!-- Avatar skeleton -->
<Skeleton variant="avatar" />

<!-- Card skeleton with multiple elements -->
<SkeletonCard />

<!-- List skeleton -->
<SkeletonList rows={5} />

<!-- Table skeleton -->
<SkeletonTable rows={10} columns={5} />
```

### LoadingBoundary

Conditional rendering of skeletons vs content:

```svelte
<LoadingBoundary isLoading={$query.isLoading}>
  {#if $query.data}
    <UserCard user={$query.data} />
  {/if}
</LoadingBoundary>
```

## Error Handling

### API Error Display (`$lib/components/ui/ApiErrorDisplay.svelte`)

Displays API error responses with field-level errors:

```svelte
<ApiErrorDisplay 
  error={apiError}
  onRetry={() => refetchData()}
/>
```

Shows:
- Main error message
- Field-specific validation errors
- Request ID for debugging
- Optional retry button

### EmptyState

Empty state UI for no results:

```svelte
<EmptyState
  title="No users found"
  description="Try adjusting your filters"
  icon="users"
/>
```

---

## Best Practices

-   **Styling**: Use the HSL-based Tailwind variables (e.g., `text-primary`, `bg-muted`) instead of hardcoded colors.
-   **Transitions**: Use Svelte's built-in `fade` and `slide` transitions for a premium, non-static user experience.
-   **Accessibility**: Always provide descriptive labels and follow the patterns established in the Bits UI-based components.
-   **Form Validation**: Use `createForm` for all forms instead of managing form state manually.
-   **Error Handling**: Let `errorStore` handle global errors; display field errors via `ApiErrorDisplay`.
-   **Loading States**: Always show skeleton loaders while data is fetching for better UX.
-   **Notifications**: Use `toast` for user feedback (success, error, warning, info).
