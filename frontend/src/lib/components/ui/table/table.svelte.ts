import {
    createTable,
    type TableOptions,
    type TableState,
    type RowData,
    type Updater,
} from '@tanstack/table-core';

/**
 * A Svelte 5 reactive adapter for TanStack Table Core.
 * This replaces the incompatible @tanstack/svelte-table v8.
 */
export function createSvelteTable<TData extends RowData>(options: TableOptions<TData>) {
    // 1. Initialize the core table instance with initial options.
    // We cast as any because createTable might be overly strict in certain versions.
    const table = createTable(options as any);

    // 2. Wrap state in Svelte 5 reactive $state
    let state = $state<TableState>(table.initialState);

    // 3. Update core options to use reactive state and a listener
    // We use a simpler setOptions call to avoid strict Resolved types issue
    table.setOptions((prev) => ({
        ...prev,
        ...options,
        state: {
            ...state,
            ...options.state,
        },
        // Hook into state changes
        onStateChange: (updater: Updater<TableState>) => {
            if (updater instanceof Function) {
                state = updater(state);
            } else {
                state = updater;
            }
            options.onStateChange?.(updater);
        },
    } as any)); // Using 'as any' here because TanStack Core's functional updater 
                // is known to have overly strict type requirements for Resolved options.

    return table;
}
