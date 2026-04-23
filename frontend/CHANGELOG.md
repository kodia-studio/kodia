# Changelog - Kodia Frontend 🐨🎨

All notable changes to the Kodia Frontend (SvelteKit + Kodia UI).

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
