/**
 * Kodia DevStore 🐨🛠️
 * Reactive intelligence for developer introspection.
 */

interface ApiLog {
	id: string;
	method: string;
	path: string;
	status?: number;
	duration?: number;
	requestBody?: any;
	responseBody?: any;
	timestamp: Date;
}

interface SocketLog {
	type: 'send' | 'receive';
	data: any;
	timestamp: Date;
}

class DevStore {
	logs = $state<ApiLog[]>([]);
	socketLogs = $state<SocketLog[]>([]);
	isOpen = $state(false);

	/**
	 * Log an API request
	 */
	logRequest(log: Omit<ApiLog, 'timestamp'>) {
		this.logs = [{ ...log, timestamp: new Date() }, ...this.logs].slice(0, 50);
	}

	/**
	 * Update an existing API log with response data
	 */
	updateLog(id: string, update: Partial<ApiLog>) {
		const index = this.logs.findIndex(l => l.id === id);
		if (index !== -1) {
			this.logs[index] = { ...this.logs[index], ...update };
		}
	}

	/**
	 * Log a WebSocket message
	 */
	logSocket(type: 'send' | 'receive', data: any) {
		this.socketLogs = [{ type, data, timestamp: new Date() }, ...this.socketLogs].slice(0, 100);
	}

	toggle() {
		this.isOpen = !this.isOpen;
	}

	clear() {
		this.logs = [];
		this.socketLogs = [];
	}
}

export const devStore = new DevStore();
