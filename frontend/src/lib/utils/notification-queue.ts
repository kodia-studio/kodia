/**
 * Offline Message Queue
 * Stores messages received while offline and processes them when back online
 */

export interface QueuedNotification {
    id: string;
    payload: any;
    timestamp: number;
    retry_count: number;
}

class NotificationQueue {
    private queue: QueuedNotification[] = [];
    private readonly MAX_QUEUE_SIZE = 100;
    private readonly STORAGE_KEY = 'notification_queue';

    constructor() {
        this.loadFromStorage();
    }

    /**
     * Add notification to queue
     */
    add(notification: any): QueuedNotification {
        if (this.queue.length >= this.MAX_QUEUE_SIZE) {
            // Remove oldest if queue is full
            this.queue.shift();
        }

        const queued: QueuedNotification = {
            id: notification.id || `offline-${Date.now()}`,
            payload: notification,
            timestamp: Date.now(),
            retry_count: 0
        };

        this.queue.push(queued);
        this.saveToStorage();

        console.log(`[Notification Queue] Added to queue (size: ${this.queue.length})`);
        return queued;
    }

    /**
     * Get all queued notifications
     */
    getAll(): QueuedNotification[] {
        return [...this.queue];
    }

    /**
     * Remove notification from queue
     */
    remove(id: string): void {
        this.queue = this.queue.filter(n => n.id !== id);
        this.saveToStorage();
    }

    /**
     * Clear all queued notifications
     */
    clear(): void {
        this.queue = [];
        this.saveToStorage();
    }

    /**
     * Get queue size
     */
    size(): number {
        return this.queue.length;
    }

    /**
     * Increment retry count
     */
    incrementRetry(id: string): void {
        const notif = this.queue.find(n => n.id === id);
        if (notif) {
            notif.retry_count++;
        }
    }

    /**
     * Check if notification has exceeded max retries
     */
    hasExceededMaxRetries(id: string, maxRetries: number = 3): boolean {
        const notif = this.queue.find(n => n.id === id);
        return notif ? notif.retry_count >= maxRetries : false;
    }

    /**
     * Save queue to localStorage
     */
    private saveToStorage(): void {
        if (typeof window !== 'undefined') {
            try {
                localStorage.setItem(this.STORAGE_KEY, JSON.stringify(this.queue));
            } catch (err) {
                console.error('[Notification Queue] Failed to save to storage:', err);
            }
        }
    }

    /**
     * Load queue from localStorage
     */
    private loadFromStorage(): void {
        if (typeof window !== 'undefined') {
            try {
                const stored = localStorage.getItem(this.STORAGE_KEY);
                if (stored) {
                    this.queue = JSON.parse(stored);
                    console.log(`[Notification Queue] Loaded ${this.queue.length} queued notifications from storage`);
                }
            } catch (err) {
                console.error('[Notification Queue] Failed to load from storage:', err);
                this.queue = [];
            }
        }
    }
}

export const notificationQueue = new NotificationQueue();
