/**
 * Kodia Optimistic Update Helper 🐨⚡
 * High-performance UI synchronization with server rollback.
 */

export function createOptimistic<T>(
	initialValue: T,
	sync: (value: T) => Promise<any>
) {
	let current = $state<T>(initialValue);
	let previous = initialValue;
	let isSyncing = $state(false);
	let error = $state<any>(null);

	return {
		get value() { return current; },
		get isSyncing() { return isSyncing; },
		get error() { return error; },

		/**
		 * Update the value optimistically and sync with server
		 */
		async update(newValue: T | ((prev: T) => T)) {
			previous = current;
			error = null;

			// Step 1: Instant UI Update
			if (typeof newValue === 'function') {
				current = (newValue as (prev: T) => T)(current);
			} else {
				current = newValue;
			}

			// Step 2: Server Sync
			isSyncing = true;
			try {
				await sync(current);
			} catch (err) {
				// Step 3: Rollback on Error
				error = err;
				current = previous;
				throw err;
			} finally {
				isSyncing = false;
			}
		},

		/**
		 * Reset the state
		 */
		reset(value: T) {
			current = value;
			previous = value;
			error = null;
		}
	};
}
