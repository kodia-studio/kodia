import { PUBLIC_API_URL } from '$env/static/public';
import { authStore } from './auth.store';
import { devStore } from './dev.svelte';
import { get } from 'svelte/store';

/**
 * Kodia WebSocket Store 🐨🔌
 * Reactive Svelte 5 store for real-time synchronization.
 */

class SocketStore {
	private ws: WebSocket | null = null;
	private reconnectAttempts = 0;
	private maxReconnectAttempts = 5;
	
	status = $state<'connecting' | 'open' | 'closed' | 'error'>('closed');
	messages = $state<any[]>([]);
	subscriptions = $state<Set<string>>(new Set());

	constructor() {
		// Auto-connect if token exists
		$effect(() => {
			const auth = get(authStore);
			if (auth.isAuthenticated && this.status === 'closed') {
				this.connect();
			} else if (!auth.isAuthenticated && this.status === 'open') {
				this.disconnect();
			}
		});
	}

	connect() {
		if (this.ws) return;

		const wsUrl = PUBLIC_API_URL.replace('http', 'ws') + '/ws';
		const token = get(authStore).accessToken;
		
		this.status = 'connecting';
		this.ws = new WebSocket(`${wsUrl}?token=${token}`);

		this.ws.onopen = () => {
			this.status = 'open';
			this.reconnectAttempts = 0;
			console.log('[Socket] Connected to Kodia Hub');
			
			// Re-subscribe to active channels
			this.subscriptions.forEach(channel => this.send('subscribe', { channel }));
		};

		this.ws.onmessage = (event) => {
			const data = JSON.parse(event.data);
			devStore.logSocket('receive', data);
			this.messages = [...this.messages, data].slice(-100); // Keep last 100
		};

		this.ws.onclose = () => {
			this.status = 'closed';
			this.ws = null;
			this.attemptReconnect();
		};

		this.ws.onerror = () => {
			this.status = 'error';
		};
	}

	disconnect() {
		if (this.ws) {
			this.ws.close();
			this.ws = null;
			this.status = 'closed';
		}
	}

	private attemptReconnect() {
		if (this.reconnectAttempts < this.maxReconnectAttempts) {
			this.reconnectAttempts++;
			const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
			console.log(`[Socket] Reconnecting in ${delay}ms...`);
			setTimeout(() => this.connect(), delay);
		}
	}

	send(type: string, payload: any) {
		if (this.status === 'open' && this.ws) {
			const data = { type, payload };
			devStore.logSocket('send', data);
			this.ws.send(JSON.stringify(data));
		}
	}

	subscribe(channel: string) {
		this.subscriptions.add(channel);
		this.send('subscribe', { channel });
	}

	unsubscribe(channel: string) {
		this.subscriptions.delete(channel);
		this.send('unsubscribe', { channel });
	}
}

export const socketStore = new SocketStore();
