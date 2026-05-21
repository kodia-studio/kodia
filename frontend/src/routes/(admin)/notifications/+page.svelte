<script lang="ts">
  import { Bell, Trash2, Check, X, Loader } from "lucide-svelte";
  import { api } from "$lib/api/client.svelte";
  import { toast } from "svelte-sonner";
  import { cn } from "$lib/utils/styles";
  import { websocketStore } from "$lib/stores/websocket.store.svelte";
  import { authStore } from "$lib/stores/auth.store";

  let notifications = $state<any[]>([]);
  let isLoading = $state(true);
  let unreadCount = $state(0);
  let initialized = $state(false);
  let wsUnsubscribe: (() => void) | null = null;

  async function fetchNotifications() {
    isLoading = true;
    try {
      const res = await api.get<any>("/api/notifications");
      notifications = res?.notifications || [];
      const count = await api.get<{ count: number }>("/api/notifications/unread/count");
      unreadCount = count.count || 0;
    } catch (err: any) {
      toast.error(err.message || "Failed to load notifications");
    } finally {
      isLoading = false;
    }
  }

  async function markAsRead(id: string) {
    try {
      await api.post(`/api/notifications/${id}/read`, {});
      notifications = notifications.map(n =>
        n.id === id ? { ...n, is_read: true } : n
      );
      unreadCount = Math.max(0, unreadCount - 1);
      authStore.update(s => ({
        ...s,
        notificationUnreadCount: Math.max(0, (s.notificationUnreadCount ?? 0) - 1)
      }));
      toast.success("Marked as read");
    } catch (err: any) {
      toast.error(err.message || "Failed to mark as read");
    }
  }

  async function markAllAsRead() {
    try {
      await api.post("/api/notifications/read-all", {});
      notifications = notifications.map(n => ({ ...n, is_read: true }));
      unreadCount = 0;
      authStore.update(s => ({
        ...s,
        notificationUnreadCount: 0
      }));
      toast.success("All marked as read");
    } catch (err: any) {
      toast.error(err.message || "Failed to mark all as read");
    }
  }

  async function deleteNotification(id: string) {
    try {
      await api.delete(`/api/notifications/${id}`);
      const deleted = notifications.find(n => n.id === id);
      notifications = notifications.filter(n => n.id !== id);

      // Update unread count if deleted notification was unread
      if (deleted && !deleted.is_read) {
        unreadCount = Math.max(0, unreadCount - 1);
        authStore.update(s => ({
          ...s,
          notificationUnreadCount: Math.max(0, (s.notificationUnreadCount ?? 0) - 1)
        }));
      }

      toast.success("Notification deleted");
    } catch (err: any) {
      toast.error(err.message || "Failed to delete notification");
    }
  }

  function handleNewNotification(msg: any) {
    try {
      // Extract payload
      const payload = msg.payload || msg;
      const newNotif = {
        id: payload.id || `notif-${Date.now()}`,
        title: payload.title || 'Notification',
        message: payload.message || '',
        type: payload.type || 'info',
        is_read: false,
        created_at: new Date().toISOString(),
        data: payload.data || {}
      };

      // Prepend to list with animation
      notifications = [newNotif, ...notifications];

      // Update unread count
      unreadCount = (unreadCount || 0) + 1;
    } catch (err) {
      console.error('[Notifications Page] Error handling new notification:', err);
    }
  }

  $effect(() => {
    if (!initialized) {
      fetchNotifications();

      // Setup WebSocket listener for real-time updates
      wsUnsubscribe = websocketStore.on('notification', handleNewNotification);
      initialized = true;

      return () => {
        if (wsUnsubscribe) {
          wsUnsubscribe();
        }
      };
    }
  });

  function getNotificationIcon(type: string) {
    switch (type) {
      case "success":
        return { color: "bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400", icon: "✓" };
      case "error":
        return { color: "bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400", icon: "!" };
      case "warning":
        return { color: "bg-amber-100 dark:bg-amber-900/30 text-amber-600 dark:text-amber-400", icon: "⚠" };
      default:
        return { color: "bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400", icon: "ℹ" };
    }
  }

  function formatDate(dateString: string): string {
    try {
      const date = new Date(dateString);
      if (isNaN(date.getTime())) return "—";
      const now = new Date();
      const diff = now.getTime() - date.getTime();
      const minutes = Math.floor(diff / 60000);
      const hours = Math.floor(diff / 3600000);
      const days = Math.floor(diff / 86400000);

      if (minutes < 1) return "Just now";
      if (minutes < 60) return `${minutes}m ago`;
      if (hours < 24) return `${hours}h ago`;
      if (days < 7) return `${days}d ago`;
      return date.toLocaleDateString("en-US", { month: "short", day: "numeric", year: "numeric" });
    } catch {
      return "—";
    }
  }
</script>

<svelte:head>
  <title>Notifications | Kodia Console</title>
</svelte:head>

<div class="space-y-6 md:space-y-8 pb-4">
  <!-- Header -->
  <div class="flex flex-col gap-4 md:gap-6">
    <div>
      <h1 class="text-2xl md:text-4xl font-black font-heading tracking-tight text-slate-900 dark:text-white leading-none">Notifications</h1>
      <p class="text-xs md:text-sm font-medium text-slate-400 mt-2 md:mt-3 leading-none">Stay updated with your messages and alerts</p>
    </div>

    <div class="flex flex-col sm:flex-row sm:items-center gap-3 sm:justify-between">
      <div class="flex items-center gap-2 sm:gap-3">
        {#if unreadCount > 0}
          <div class="flex items-center gap-2 px-3 md:px-4 py-2 md:py-2.5 bg-blue-500/10 rounded-lg md:rounded-xl border border-blue-200 dark:border-blue-900/50">
            <div class="w-2 h-2 rounded-full bg-blue-500 shadow-[0_0_8px_rgba(59,130,246,0.5)]"></div>
            <span class="text-xs font-bold text-blue-600 dark:text-blue-400">{unreadCount} unread</span>
          </div>
        {/if}
      </div>

      {#if unreadCount > 0}
        <button
          onclick={markAllAsRead}
          class="w-full sm:w-auto px-4 md:px-6 py-2 md:py-3 bg-slate-100 dark:bg-slate-800 text-slate-900 dark:text-white rounded-lg md:rounded-xl text-xs md:text-sm font-semibold hover:bg-slate-200 dark:hover:bg-slate-700 transition-colors flex items-center justify-center sm:justify-start gap-2"
        >
          <Check class="w-4 h-4" />
          <span>Mark all as read</span>
        </button>
      {/if}
    </div>
  </div>

  <!-- Notifications -->
  <div class="space-y-4">
    {#if isLoading}
      <div class="h-64 flex items-center justify-center">
        <div class="text-center">
          <div class="w-10 h-10 rounded-full border-2 border-slate-200 dark:border-slate-700 border-t-blue-500 animate-spin mx-auto mb-3"></div>
          <p class="text-sm text-slate-600 dark:text-slate-400">Loading notifications...</p>
        </div>
      </div>
    {:else if notifications.length === 0}
      <div class="flex flex-col items-center justify-center py-12 md:py-16 px-4 rounded-xl md:rounded-2xl border-2 border-dashed border-slate-200 dark:border-slate-800">
        <div class="w-12 h-12 md:w-16 md:h-16 rounded-lg md:rounded-2xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center mb-3 md:mb-4">
          <Bell class="w-6 h-6 md:w-8 md:h-8 text-slate-400" />
        </div>
        <p class="text-base md:text-lg font-semibold text-slate-900 dark:text-white mb-2 text-center">No notifications yet</p>
        <p class="text-xs md:text-sm text-slate-500 dark:text-slate-400 text-center max-w-xs">You're all caught up! Come back later for updates.</p>
      </div>
    {:else}
      {#each notifications as notification (notification.id)}
        <div
          class={cn(
            "p-4 md:p-6 rounded-lg md:rounded-2xl border transition-all duration-300 animate-in slide-in-from-top-2",
            notification.is_read
              ? "bg-slate-50 dark:bg-slate-900/30 border-slate-200 dark:border-slate-800"
              : "bg-white dark:bg-slate-900 border-blue-200 dark:border-blue-900/50 shadow-[0_0_20px_rgba(59,130,246,0.1)]"
          )}
        >
          <div class="flex flex-col sm:flex-row gap-3 md:gap-4">
            <!-- Icon -->
            <div
              class={cn(
                "w-10 h-10 md:w-12 md:h-12 rounded-lg md:rounded-xl flex items-center justify-center text-base md:text-lg font-bold shrink-0",
                getNotificationIcon(notification.type).color
              )}
            >
              {getNotificationIcon(notification.type).icon}
            </div>

            <!-- Content -->
            <div class="flex-1 min-w-0">
              <div class="flex flex-col sm:flex-row sm:items-start justify-between gap-2 sm:gap-4 mb-1 md:mb-2">
                <h3 class={cn(
                  "text-xs md:text-sm font-bold leading-tight",
                  notification.is_read
                    ? "text-slate-700 dark:text-slate-300"
                    : "text-slate-900 dark:text-white"
                )}>
                  {notification.title}
                  {#if !notification.is_read}
                    <span class="ml-2 inline-block w-2 h-2 bg-blue-500 rounded-full"></span>
                  {/if}
                </h3>
                <span class={cn(
                  "text-[11px] md:text-xs font-medium whitespace-nowrap shrink-0",
                  notification.is_read
                    ? "text-slate-500 dark:text-slate-500"
                    : "text-blue-600 dark:text-blue-400"
                )}>
                  {formatDate(notification.created_at)}
                </span>
              </div>

              <p class={cn(
                "text-xs md:text-sm leading-relaxed",
                notification.is_read
                  ? "text-slate-600 dark:text-slate-400"
                  : "text-slate-700 dark:text-slate-300"
              )}>
                {notification.message}
              </p>
            </div>

            <!-- Actions -->
            <div class="flex gap-1 md:gap-2 shrink-0 self-start sm:self-center">
              {#if !notification.is_read}
                <button
                  onclick={() => markAsRead(notification.id)}
                  class="p-1.5 md:p-2 hover:bg-blue-100 dark:hover:bg-blue-900/30 rounded-lg transition-colors text-blue-600 dark:text-blue-400"
                  title="Mark as read"
                >
                  <Check class="w-4 h-4" />
                </button>
              {/if}
              <button
                onclick={() => deleteNotification(notification.id)}
                class="p-1.5 md:p-2 hover:bg-red-100 dark:hover:bg-red-900/30 rounded-lg transition-colors text-red-600 dark:text-red-400"
                title="Delete notification"
              >
                <Trash2 class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      {/each}
    {/if}
  </div>
</div>
