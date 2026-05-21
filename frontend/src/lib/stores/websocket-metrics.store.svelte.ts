import { writable, derived } from 'svelte/store';
import { websocketStore } from './websocket.store.svelte';
import { notificationQueue } from '$lib/utils/notification-queue';

export interface WebSocketMetrics {
    isConnected: boolean;
    isConnecting: boolean;
    connectionAttempts: number;
    failedAttempts: number;
    lastConnectedAt: number | null;
    lastDisconnectedAt: number | null;
    uptime: number; // seconds
    totalMessages: number;
    queuedNotifications: number;
}

const initialMetrics: WebSocketMetrics = {
    isConnected: false,
    isConnecting: false,
    connectionAttempts: 0,
    failedAttempts: 0,
    lastConnectedAt: null,
    lastDisconnectedAt: null,
    uptime: 0,
    totalMessages: 0,
    queuedNotifications: 0
};

function createWebSocketMetricsStore() {
    const { subscribe, set, update } = writable<WebSocketMetrics>(initialMetrics);

    let connectionStartTime: number | null = null;
    let uptimeInterval: NodeJS.Timeout | null = null;

    // Subscribe to WebSocket store changes
    const unsubscribeWS = websocketStore.subscribe(wsState => {
        update(metrics => ({
            ...metrics,
            isConnected: wsState.isConnected,
            isConnecting: wsState.isConnecting,
            queuedNotifications: notificationQueue.size()
        }));

        // Track connection attempts
        if (wsState.isConnecting && !connectionStartTime) {
            connectionStartTime = Date.now();
            update(metrics => ({
                ...metrics,
                connectionAttempts: metrics.connectionAttempts + 1
            }));
        }

        // Track successful connection
        if (wsState.isConnected && connectionStartTime) {
            connectionStartTime = null;
            update(metrics => ({
                ...metrics,
                lastConnectedAt: Date.now(),
                uptime: 0
            }));

            // Start uptime tracking
            if (uptimeInterval) clearInterval(uptimeInterval);
            uptimeInterval = setInterval(() => {
                update(metrics => ({
                    ...metrics,
                    uptime: metrics.uptime + 1
                }));
            }, 1000);
        }

        // Track disconnection
        if (!wsState.isConnected && !wsState.isConnecting && connectionStartTime === null) {
            update(metrics => ({
                ...metrics,
                lastDisconnectedAt: Date.now(),
                uptime: 0
            }));

            if (uptimeInterval) {
                clearInterval(uptimeInterval);
                uptimeInterval = null;
            }
        }

        // Track failed attempts
        if (wsState.error && !wsState.isConnected && !wsState.isConnecting) {
            update(metrics => ({
                ...metrics,
                failedAttempts: metrics.failedAttempts + 1
            }));
        }
    });

    // Track total messages
    const trackMessage = () => {
        update(metrics => ({
            ...metrics,
            totalMessages: metrics.totalMessages + 1
        }));
    };

    const unsubscribeMessages = websocketStore.on('message', trackMessage);

    return {
        subscribe,
        cleanup: () => {
            unsubscribeWS();
            unsubscribeMessages?.();
            if (uptimeInterval) clearInterval(uptimeInterval);
        }
    };
}

export const websocketMetricsStore = createWebSocketMetricsStore();

// Derived store for connection status text
export const connectionStatusText = derived(websocketStore, $ws => {
    if ($ws.isConnected) return 'Connected';
    if ($ws.isConnecting) return 'Connecting...';
    if ($ws.error) return 'Connection failed';
    return 'Disconnected';
});

// Derived store for formatted uptime
export const formattedUptime = derived(websocketMetricsStore, $metrics => {
    const seconds = $metrics.uptime;
    if (seconds < 60) return `${seconds}s`;
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes}m`;
    const hours = Math.floor(minutes / 60);
    return `${hours}h`;
});
