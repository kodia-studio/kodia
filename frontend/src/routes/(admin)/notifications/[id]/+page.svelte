<script lang="ts">
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { ArrowLeft, Check, Trash2 } from 'lucide-svelte';
  import { api } from '$lib/api/client.svelte';
  import { toast } from 'svelte-sonner';
  import { authStore } from '$lib/stores/auth.store';

  let notification = $state<any>(null);
  let isLoading = $state(true);
  let isDeleting = $state(false);

  async function loadNotification() {
    isLoading = true;
    try {
      const res = await api.get<any>(`/api/notifications`);
      const notifs = res?.notifications || [];
      const found = notifs.find((n: any) => n.id === page.params.id);

      if (!found) {
        toast.error('Notification not found');
        goto('/notifications');
        return;
      }

      notification = found;

      // Mark as read if not already read
      if (!found.is_read) {
        try {
          await api.post(`/api/notifications/${found.id}/read`, {});
          authStore.update(s => ({
            ...s,
            notificationUnreadCount: Math.max(0, (s.notificationUnreadCount ?? 0) - 1)
          }));
          notification.is_read = true;
        } catch (err) {
          console.error('Failed to mark as read:', err);
        }
      }
    } catch (err: any) {
      toast.error(err.message || 'Failed to load notification');
      goto('/notifications');
    } finally {
      isLoading = false;
    }
  }

  async function handleDelete() {
    if (!notification) return;

    isDeleting = true;
    try {
      await api.delete(`/api/notifications/${notification.id}`);

      if (!notification.is_read) {
        authStore.update(s => ({
          ...s,
          notificationUnreadCount: Math.max(0, (s.notificationUnreadCount ?? 0) - 1)
        }));
      }

      toast.success('Notification deleted');
      goto('/notifications');
    } catch (err: any) {
      toast.error(err.message || 'Failed to delete notification');
      isDeleting = false;
    }
  }

  loadNotification();
</script>

<svelte:head>
  <title>Notification | Kodia Console</title>
</svelte:head>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex items-center gap-4">
    <button
      onclick={() => goto('/notifications')}
      class="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
      title="Back to notifications"
    >
      <ArrowLeft class="w-5 h-5 text-slate-600 dark:text-slate-400" />
    </button>
    <div>
      <h1 class="text-4xl font-black font-heading tracking-tight text-slate-900 dark:text-white leading-none">
        Notification
      </h1>
      <p class="text-xs font-medium text-slate-400 mt-3 leading-none">View notification details</p>
    </div>
  </div>

  <!-- Content -->
  {#if isLoading}
    <div class="flex items-center justify-center py-16">
      <div class="text-center">
        <div class="w-10 h-10 rounded-full border-2 border-slate-200 dark:border-slate-700 border-t-blue-500 animate-spin mx-auto mb-3"></div>
        <p class="text-sm text-slate-600 dark:text-slate-400">Loading notification...</p>
      </div>
    </div>
  {:else if notification}
    <div class="space-y-6">
      <!-- Card -->
      <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-white/10 rounded-2xl p-8 shadow-sm">
        <!-- Status Badge -->
        <div class="mb-6 flex items-center gap-2">
          {#if notification.is_read}
            <span class="inline-flex items-center gap-2 px-3 py-1 bg-slate-100 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-xs font-semibold text-slate-700 dark:text-slate-300">
              <Check class="w-3 h-3" />
              Read
            </span>
          {:else}
            <span class="inline-flex items-center gap-2 px-3 py-1 bg-blue-100 dark:bg-blue-900/30 border border-blue-200 dark:border-blue-800 rounded-lg text-xs font-semibold text-blue-700 dark:text-blue-300">
              <div class="w-2 h-2 bg-blue-500 rounded-full"></div>
              Unread
            </span>
          {/if}
        </div>

        <!-- Title -->
        <h2 class="text-3xl font-black text-slate-900 dark:text-white mb-2">{notification.title}</h2>

        <!-- Timestamp -->
        <p class="text-xs font-medium text-slate-500 dark:text-slate-400 mb-6">
          {new Date(notification.created_at).toLocaleString('en-US', {
            weekday: 'long',
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
          })}
        </p>

        <!-- Divider -->
        <div class="h-px bg-slate-200 dark:bg-white/5 my-6"></div>

        <!-- Message -->
        <p class="text-base leading-relaxed text-slate-700 dark:text-slate-300 whitespace-pre-wrap">
          {notification.message}
        </p>

        <!-- Data if present -->
        {#if notification.data && Object.keys(notification.data).length > 0}
          <div class="mt-8 pt-6 border-t border-slate-200 dark:border-white/5">
            <p class="text-xs font-black uppercase tracking-[0.2em] text-slate-400 mb-4">Additional Data</p>
            <div class="space-y-3">
              {#each Object.entries(notification.data) as [key, value]}
                <div class="flex items-start gap-3 bg-slate-50 dark:bg-slate-800/50 p-3 rounded-lg">
                  <p class="text-xs font-semibold text-slate-600 dark:text-slate-400 shrink-0 min-w-25">{key}:</p>
                  <p class="text-sm text-slate-900 dark:text-white break-all">{JSON.stringify(value)}</p>
                </div>
              {/each}
            </div>
          </div>
        {/if}
      </div>

      <!-- Actions -->
      <div class="flex items-center gap-3">
        <button
          onclick={() => goto('/notifications')}
          class="px-6 py-2.5 bg-slate-100 dark:bg-slate-800 text-slate-900 dark:text-white rounded-xl text-sm font-semibold hover:bg-slate-200 dark:hover:bg-slate-700 transition-colors"
        >
          Back to list
        </button>

        <button
          onclick={handleDelete}
          disabled={isDeleting}
          class="px-6 py-2.5 bg-red-500/10 text-red-600 dark:text-red-400 border border-red-200 dark:border-red-900/50 rounded-xl text-sm font-semibold hover:bg-red-500/20 dark:hover:bg-red-900/30 transition-colors disabled:opacity-50 flex items-center gap-2"
        >
          <Trash2 class="w-4 h-4" />
          {isDeleting ? 'Deleting...' : 'Delete'}
        </button>
      </div>
    </div>
  {/if}
</div>
