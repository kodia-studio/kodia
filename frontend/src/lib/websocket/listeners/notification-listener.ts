import { websocketStore } from '$lib/stores/websocket.store.svelte';
import { authStore } from '$lib/stores/auth.store';
import { toast } from 'svelte-sonner';
import { notificationQueue } from '$lib/utils/notification-queue';
import { notificationSoundManager } from '$lib/utils/notification-sound';

export interface NotificationMessage {
    type: string;
    event: string;
    payload: {
        title: string;
        message: string;
        type?: 'info' | 'success' | 'warning' | 'error';
        data?: Record<string, any>;
    };
    timestamp: number;
}

// Duplicate detection: track recently processed notification IDs
const processedNotifications = new Set<string>();
// Store for 5 minutes to prevent duplicates from polling + WebSocket
const DEDUP_TTL = 5 * 60 * 1000;

let unsubscribe: (() => void) | null = null;

export function setupNotificationListener() {
    if (typeof window !== 'undefined') {
        const initAudio = () => {
            notificationSoundManager.initAudioContext();
            window.removeEventListener('click', initAudio);
        };
        window.addEventListener('click', initAudio);
    }

    unsubscribe = websocketStore.on('notification', handleNotification);
}

export function teardownNotificationListener() {
    if (unsubscribe) {
        unsubscribe();
        unsubscribe = null;
    }
}

function handleNotification(msg: any) {
    try {
        console.log('[Notification Listener] Received notification:', msg);

        // Extract payload
        const payload = msg.payload || msg;
        const title = payload.title || 'Notification';
        const message = payload.message || '';
        const notificationType = payload.type || 'info';

        // Deduplication: check if we've already processed this notification
        // Use timestamp as unique identifier if ID not available
        const notificationId = payload.id || `${title}-${message}-${msg.timestamp}`;

        if (processedNotifications.has(notificationId)) {
            console.log('[Notification Listener] Duplicate notification ignored:', notificationId);
            return;
        }

        // Mark as processed
        processedNotifications.add(notificationId);

        // Auto-cleanup after TTL
        setTimeout(() => {
            processedNotifications.delete(notificationId);
        }, DEDUP_TTL);

        // Update unread count
        authStore.update(state => {
            const currentCount = state.notificationUnreadCount ?? 0;
            return {
                ...state,
                notificationUnreadCount: currentCount + 1
            };
        });

        // Play notification sound
        notificationSoundManager.play(notificationType as any);

        // Show toast notification
        const toastConfig = {
            info: () => toast.info(message, { description: title }),
            success: () => toast.success(message, { description: title }),
            warning: () => toast.warning(message, { description: title }),
            error: () => toast.error(message, { description: title })
        };

        (toastConfig[notificationType as keyof typeof toastConfig] || toastConfig.info)();

        // Emit custom event untuk components yang perlu handle notification
        if (typeof window !== 'undefined') {
            window.dispatchEvent(
                new CustomEvent('notification:received', {
                    detail: { id: notificationId, title, message, type: notificationType, payload }
                })
            );
        }
    } catch (err) {
        console.error('[Notification Listener] Error handling notification:', err);
    }
}
