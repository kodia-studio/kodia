# Changelog - Kodia Frontend 🐨🎨

All notable changes to the Kodia Frontend (SvelteKit + Kodia UI).

## [1.4.0] - 2026-04-24

### Added
- **Framework Independence**: Complete refactor to build all UI components internally with zero third-party component dependencies.
- **Native Dark Mode**: Svelte 5 runes-based theme system (`theme.svelte.ts`) with localStorage persistence and system preference detection.
- **Chart System**: Rebuilt `AreaChart` and `BarChart` using pure SVG + d3-scale/d3-shape (no external charting libraries).
- **Custom Table Adapter**: Svelte 5-native `createSvelteTable()` adapter using @tanstack/table-core instead of unstable @tanstack/svelte-table.

### Removed
- `mode-watcher`: Replaced with internal theme store
- `@tanstack/svelte-table` (v8 alpha): Migrated to @tanstack/table-core with custom Svelte 5 adapter
- `layerchart`: Rebuilt charts with native SVG
- `zod`, `sveltekit-superforms` (partial), `@internationalized/date` (partial): Planned for Phase 2

### Improved
- **Build Quality**: 0 vulnerabilities, 0 warnings on fresh npm install
- **CSS**: Tailwind CSS v4 strict mode compliance - removed problematic custom class patterns
- **Performance**: Smaller bundle by removing alpha/unstable dependencies
- **Type Safety**: Full TypeScript support with Svelte 5 runes patterns

### Technical Details
- Dark mode respects `prefers-color-scheme` and persists selection to localStorage
- Charts feature interactive tooltips and hover states via Svelte 5 `$state`
- DataTable maintains full TanStack sorting/pagination without writable stores
- All components use Svelte 5 runes (`$state`, `$derived`, `$effect`)

## [1.3.0] - 2026-04-23

### Added
- **Premium Benchmark Page**: High-fidelity performance dashboard with interactive charts and animations.
- **Developer Experience (DX)**: Integrated DevTools panel (`Ctrl+K`) and institutional-grade Error Boundaries.
- **UI Components**: New `Skeleton` loaders, `Badge` variants, and `Card` components with glassmorphism.
- **i18n Engine**: Reactive internationalization powered by Svelte 5 runes with lazy-loading support.

### Fixed
- **Compiler Issues**: Resolved Svelte 5 `{@const}` placement constraints in complex loops.
- **Style Cleanup**: Removed redundant CSS rules to eliminate IDE "Unknown at rule" warnings.
- **Routing**: Fixed documentation pathing for improved SEO and user navigation.

### Improved
- **A11y**: Enhanced focus-visible management for WCAG 2.1 compliance across all interactive elements.
- **Build Speed**: Optimized Tailwind v4 integration for faster hot-module replacement (HMR).

## [1.2.0] - 2026-04-10

### Added
- **Svelte 5 Transition**: Fully migrated core components to use Runes (`$state`, `$derived`, `$effect`).
- **Infinite Scroll**: Reusable list component with virtualization support.
- **Optimistic Updates**: Built-in support for UI feedback before server confirmation.

### Fixed
- Resolved hydration issues in complex layout nesting.
- Fixed memory leaks in WebSocket store listeners.

## [1.1.0] - 2026-03-20

### Added
- **Kodia UI Library**: Initial set of institutional components (Button, Input, Table, Modal).
- **Authentication Store**: Reactive user session management with auto-refresh support.

---
© 2026 Kodia Studio. "Build like a user, code like a pro."
