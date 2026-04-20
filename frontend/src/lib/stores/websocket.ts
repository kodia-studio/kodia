import { writable, derived } from 'svelte/store'

export type MessageHandler = (message: any) => void

interface WSMessage {
	type: string
	payload: any
	room_id?: string
	user_id?: string
	timestamp: number
}

interface WSStoreState {
	status: 'disconnected' | 'connecting' | 'connected' | 'error'
	lastMessage: WSMessage | null
	error: string | null
}

const defaultState: WSStoreState = {
	status: 'disconnected',
	lastMessage: null,
	error: null
}

// Create the main store
const store = writable<WSStoreState>(defaultState)

let ws: WebSocket | null = null
let reconnectAttempts = 0
const maxReconnectAttempts = 5
const baseReconnectDelay = 1000 // 1 second
let messageHandlers: Map<string, Set<MessageHandler>> = new Map()

/**
 * Calculates exponential backoff delay for reconnection.
 * Delay increases: 1s, 2s, 4s, 8s, 16s, etc.
 */
function getReconnectDelay(): number {
	return baseReconnectDelay * Math.pow(2, reconnectAttempts - 1)
}

/**
 * Connect to the WebSocket server.
 * @param token JWT authentication token
 * @param room Optional room ID for room-based connections
 */
export function connect(token: string, room?: string): void {
	if (ws) {
		console.warn('WebSocket already connected')
		return
	}

	store.update(s => ({ ...s, status: 'connecting' }))

	const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
	const host = window.location.host
	const roomPath = room ? `/room/${room}` : ''
	const url = `${protocol}//${host}/api/ws${roomPath}?token=${encodeURIComponent(token)}`

	try {
		ws = new WebSocket(url)

		ws.onopen = () => {
			console.log('WebSocket connected')
			store.update(s => ({
				...s,
				status: 'connected',
				error: null
			}))
			reconnectAttempts = 0
		}

		ws.onmessage = (event: MessageEvent) => {
			try {
				const message: WSMessage = JSON.parse(event.data)
				store.update(s => ({ ...s, lastMessage: message }))

				// Call registered handlers for this message type
				const handlers = messageHandlers.get(message.type)
				if (handlers) {
					handlers.forEach(handler => {
						try {
							handler(message)
						} catch (err) {
							console.error(`Error in message handler for ${message.type}:`, err)
						}
					})
				}
			} catch (err) {
				console.error('Failed to parse WebSocket message:', err)
			}
		}

		ws.onerror = (event: Event) => {
			console.error('WebSocket error:', event)
			store.update(s => ({
				...s,
				status: 'error',
				error: 'WebSocket connection error'
			}))
		}

		ws.onclose = () => {
			console.log('WebSocket disconnected')
			ws = null
			store.update(s => ({ ...s, status: 'disconnected' }))

			// Attempt to reconnect with exponential backoff
			if (reconnectAttempts < maxReconnectAttempts) {
				reconnectAttempts++
				const delay = getReconnectDelay()
				console.log(`Reconnecting in ${delay}ms (attempt ${reconnectAttempts}/${maxReconnectAttempts})`)
				setTimeout(() => connect(token, room), delay)
			} else {
				store.update(s => ({
					...s,
					error: 'Max reconnection attempts reached'
				}))
			}
		}
	} catch (err) {
		console.error('Failed to create WebSocket connection:', err)
		store.update(s => ({
			...s,
			status: 'error',
			error: err instanceof Error ? err.message : 'Failed to connect'
		}))
	}
}

/**
 * Disconnect from the WebSocket server.
 */
export function disconnect(): void {
	if (ws) {
		ws.close()
		ws = null
		messageHandlers.clear()
	}
	store.update(s => ({
		...s,
		status: 'disconnected',
		lastMessage: null
	}))
}

/**
 * Send a message to the WebSocket server.
 */
export function send(message: Partial<WSMessage>): void {
	if (!ws || ws.readyState !== WebSocket.OPEN) {
		console.error('WebSocket is not connected')
		return
	}

	const payload = {
		type: message.type || 'message',
		payload: message.payload || {},
		...message
	}

	try {
		ws.send(JSON.stringify(payload))
	} catch (err) {
		console.error('Failed to send WebSocket message:', err)
	}
}

/**
 * Register a handler for a specific message type.
 * @param type Message type to listen for
 * @param handler Function to call when message of this type is received
 * @returns Unsubscribe function
 */
export function onMessage(type: string, handler: MessageHandler): () => void {
	if (!messageHandlers.has(type)) {
		messageHandlers.set(type, new Set())
	}

	messageHandlers.get(type)!.add(handler)

	// Return unsubscribe function
	return () => {
		const handlers = messageHandlers.get(type)
		if (handlers) {
			handlers.delete(handler)
			if (handlers.size === 0) {
				messageHandlers.delete(type)
			}
		}
	}
}

/**
 * Register a handler for disconnection events.
 */
export function onDisconnect(handler: () => void): () => void {
	const unsubscribe = store.subscribe(state => {
		if (state.status === 'disconnected' && ws === null) {
			handler()
		}
	})
	return unsubscribe
}

/**
 * Send a ping message (for testing connection).
 */
export function ping(): void {
	send({ type: 'ping' })
}

/**
 * The main WebSocket store with reactive state.
 * Usage:
 *   $wsStore.status
 *   $wsStore.lastMessage
 *   $wsStore.error
 */
export const wsStore = derived(store, $store => ({
	...$store,

	connect: (token: string, room?: string) => connect(token, room),
	disconnect: () => disconnect(),
	send: (msg: Partial<WSMessage>) => send(msg),
	ping: () => ping(),
	onMessage: (type: string, handler: MessageHandler) => onMessage(type, handler),
	onDisconnect: (handler: () => void) => onDisconnect(handler)
}))

export default wsStore
