import { api } from "$lib/api/client";

export function createQuery<T>(fetcher: () => Promise<T>) {
  let data = $state<T | null>(null);
  let isLoading = $state(true);
  let error = $state<any>(null);
  let lastFetched = $state<number | null>(null);

  const execute = async (force = false) => {
    // Basic caching logic (5 minute TTL)
    if (!force && lastFetched && Date.now() - lastFetched < 300000) {
      return;
    }

    isLoading = true;
    error = null;
    try {
      data = await fetcher();
      lastFetched = Date.now();
    } catch (err) {
      error = err;
    } finally {
      isLoading = false;
    }
  };

  // Initial fetch
  execute();

  return {
    get data() { return data; },
    get isLoading() { return isLoading; },
    get error() { return error; },
    refetch: () => execute(true)
  };
}
