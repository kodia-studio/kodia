import { writable } from 'svelte/store';
import { notificationQueue } from '$lib/utils/notification-queue';
import { notificationSoundManager } from '$lib/utils/notification-sound';

interface WebSocketState {
    isConnected: boolean;
    isConnecting: boolean;
    error: string | null;
    lastMessageTime: number | null;
}

const initialState: WebSocketState = {
    isConnected: false,
    isConnecting: false,
    error: null,
    lastMessageTime: null
};

function createWebSocketStore() {
    const { subscribe, set, update } = writable<WebSocketState>(initialState);

    let ws: WebSocket | null = null;
    let reconnectAttempts = 0;
    const maxReconnectAttempts = 5;
    const reconnectDelay = [1000, 2000, 4000, 8000, 15000]; // ms, exponential backoff

    const listeners = new Map<string, Function[]>();

    function on(event: string, callback: Function) {
        if (!listeners.has(event)) {
            listeners.set(event, []);
        }
        listeners.get(event)!.push(callback);

        return () => {
            const cbs = listeners.get(event)!;
            cbs.splice(cbs.indexOf(callback), 1);
        };
    }

    function emit(event: string, data: any) {
        const cbs = listeners.get(event);
        if (cbs) {
            cbs.forEach(cb => {
                try {
                    cb(data);
                } catch (err) {
                    console.error(`[WebSocket] Error in listener for "${event}":`, err);
                }
            });
        }
    }

    function connect() {
        if (ws) return;

        const token = localStorage.getItem('access_token');
        if (!token) {
            update(state => ({ ...state, isConnecting: false }));
            return;
        }

        update(state => ({ ...state, isConnecting: true }));

        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const hostname = window.location.hostname;
        const backendHost = hostname === 'localhost'
            ? 'localhost:8080'
            : window.location.host;
        const url = `${protocol}//${backendHost}/api/ws?token=${encodeURIComponent(token)}`;

        try {
            ws = new WebSocket(url);

            ws.onopen = () => {
                update(state => ({
                    ...state,
                    isConnected: true,
                    isConnecting: false,
                    error: null
                }));
                reconnectAttempts = 0;

                const queued = notificationQueue.getAll();
                if (queued.length > 0) {
                    queued.forEach(q => {
                        emit('notification', q.payload);
                        notificationQueue.remove(q.id);
                    });
                }

                emit('connect', null);
            };

            ws.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    update(state => ({ ...state, lastMessageTime: Date.now() }));
                    emit('message', message);

                    if (message.type) {
                        emit(message.type, message);
                    }
                } catch (err) {
                    console.error('[WebSocket] Failed to parse message:', err);
                }
            };

            ws.onerror = (err) => {
                if (reconnectAttempts > 0) {
                    console.error('[WebSocket] Connection error after retry:', err);
                }
                update(state => ({
                    ...state,
                    error: 'WebSocket connection error'
                }));
                emit('error', err);
            };

            ws.onclose = () => {
                ws = null;
                update(state => ({
                    ...state,
                    isConnected: false,
                    isConnecting: false
                }));
                emit('disconnect', null);

                if (reconnectAttempts < maxReconnectAttempts) {
                    const delay = reconnectDelay[reconnectAttempts] || 30000;
                    reconnectAttempts++;
                    if (reconnectAttempts > 1) {
                        console.warn(`[WebSocket] Reconnecting in ${delay}ms (attempt ${reconnectAttempts})`);
                    }
                    setTimeout(connect, delay);
                } else {
                    console.error('[WebSocket] Max reconnection attempts reached');
                    update(state => ({
                        ...state,
                        error: 'Failed to reconnect to WebSocket'
                    }));
                }
            };
        } catch (err) {
            console.error('[WebSocket] Failed to create connection:', err);
            update(state => ({
                ...state,
                isConnecting: false,
                error: err instanceof Error ? err.message : 'Unknown error'
            }));
        }
    }

    function disconnect() {
        if (ws) {
            ws.close();
            ws = null;
        }
    }

    return {
        subscribe,
        connect,
        disconnect,
        on,
        emit
    };
}

export const websocketStore = createWebSocketStore();
