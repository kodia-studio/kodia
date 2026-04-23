# Kodia UI Library 🐨💎

**Kodia UI** is an elite, institutional-grade Svelte 5 component library designed for high-velocity, premium web development. It is meticulously engineered to work seamlessly with the Kodia Framework backend, utilizing modern Svelte 5 Runes for maximum reactivity and performance.

## ✨ Features
- **Svelte 5 First**: Built from the ground up using `$state`, `$derived`, and `$props`.
- **Premium Aesthetics**: Glassmorphism, institutional-grade typography, and smooth transitions.
- **Dark Mode Native**: Support for dark mode out of the box with tailored color palettes.
- **Accessible (a11y)**: Built on top of `bits-ui` and standard ARIA practices.
- **Tailwind v4 Optimized**: Utilizes the latest Tailwind CSS features for styling.
- **Data Intelligence**: Integrated API Codegen, Optimistic Updates, and Data Virtualization.
- **Real-time Sync**: Reactive WebSocket stores for instant data synchronization.
- **Offline Ready**: Built-in PWA support with intelligent caching strategies.
- **Smart Components**: Pre-wired for Kodia Framework's backend error structures and pagination.

## 📦 Component Suite

### 🧩 Core UI
- **Button**: Premium action components with multiple variants (`premium`, `primary`, `danger`, etc.) and built-in loading states.
- **Input**: Smart input fields with integrated label and error handling.
- **Badge**: Distinctive status indicators for Enterprise ecosystems.
- **Avatar**: Professional user representation with fallback support.
- **Modal**: Glassmorphic dialogs with smooth scaling transitions.
- **Dropdown**: Elegant menu systems with fly-in animations.
- **Toast**: Sophisticated notifications using `svelte-sonner`.

### 📝 Form System
- **SuperFormField**: A declarative field component that automatically binds to Superforms proxies.
- **Form (KForm)**: A smart container that handles backend validation error mapping (`422 Unprocessable Entity`) automatically.

### 📊 Data & Media
- **DataTable (KDataTable)**: Advanced table system with TanStack Table integration.
- **InfiniteList**: High-performance scrolling for continuous data streams.
- **VirtualList**: Elite virtualization for rendering 10,000+ items with zero lag.
- **Uploader**: Professional drag-and-drop file uploader with progress tracking.
- **Editor**: Elite WYSIWYG editor powered by Tiptap.
- **BarChart**: High-fidelity data visualization built on LayerChart (D3).

### 📡 Data Layer
- **API Codegen**: Run `npm run codegen` to sync types instantly with the backend.
- **Optimistic Helper**: Use `createOptimistic()` for instant UI feedback with auto-rollback.
- **Socket Store**: Reactive `socketStore` for institutional-grade WebSocket communication.

## 🚀 Getting Started

### Installation
All components are available via the `$lib/components/ui` and `$lib/components/forms` directories.

### Basic Usage
```svelte
<script lang="ts">
  import { Button, Input, Modal } from '$lib/components/ui';
  let showModal = $state(false);
</script>

<Button variant="premium" onclick={() => showModal = true}>
  Unlock Potential
</Button>

<Modal bind:open={showModal} title="System Activation">
  <p>Ready to deploy the Kodia ecosystem?</p>
</Modal>
```

## 🛠️ Utils
- **cn()**: A professional utility for merging Tailwind classes safely using `clsx` and `tailwind-merge`.

## 📜 Principles
> "Build like a user, code like a pro."

Every component is designed to be responsive, secure, and visually stunning across all devices.

---
© 2026 Kodia Studio. All rights reserved.
