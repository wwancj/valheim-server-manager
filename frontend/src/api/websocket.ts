/**
 * WebSocket Client for Valheim Server Manager
 * Replaces Wails Events system for real-time updates
 */

type EventHandler = (data: any) => void;

class WebSocketClient {
    private ws: WebSocket | null = null;
    private handlers: Map<string, Set<EventHandler>> = new Map();
    private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
    private reconnectAttempts = 0;
    private maxReconnectAttempts = 10;
    private baseDelay = 1000;

    connect() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const url = `${protocol}//${window.location.host}/api/ws/events`;

        try {
            this.ws = new WebSocket(url);

            this.ws.onopen = () => {
                console.log('[WebSocket] Connected');
                this.reconnectAttempts = 0;
            };

            this.ws.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    const { event: eventName, data } = message;
                    this.emit(eventName, data);
                } catch (e) {
                    console.error('[WebSocket] Failed to parse message:', e);
                }
            };

            this.ws.onclose = () => {
                console.log('[WebSocket] Disconnected');
                this.ws = null;
                this.scheduleReconnect();
            };

            this.ws.onerror = (error) => {
                console.error('[WebSocket] Error:', error);
            };
        } catch (e) {
            console.error('[WebSocket] Connection failed:', e);
            this.scheduleReconnect();
        }
    }

    private scheduleReconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            console.error('[WebSocket] Max reconnect attempts reached');
            return;
        }

        const delay = this.baseDelay * Math.pow(2, this.reconnectAttempts);
        this.reconnectAttempts++;

        console.log(`[WebSocket] Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`);

        this.reconnectTimer = setTimeout(() => {
            this.connect();
        }, delay);
    }

    private emit(eventName: string, data: any) {
        const handlers = this.handlers.get(eventName);
        if (handlers) {
            handlers.forEach(handler => {
                try {
                    handler(data);
                } catch (e) {
                    console.error(`[WebSocket] Handler error for ${eventName}:`, e);
                }
            });
        }

        // Also emit wildcard events
        const wildcardHandlers = this.handlers.get('*');
        if (wildcardHandlers) {
            wildcardHandlers.forEach(handler => {
                try {
                    handler({ event: eventName, data });
                } catch (e) {
                    console.error('[WebSocket] Wildcard handler error:', e);
                }
            });
        }
    }

    /**
     * Subscribe to an event
     * @param eventName Event name (e.g., 'log:new', 'server:status') or '*' for all events
     * @param handler Callback function
     * @returns Unsubscribe function
     */
    on(eventName: string, handler: EventHandler): () => void {
        if (!this.handlers.has(eventName)) {
            this.handlers.set(eventName, new Set());
        }
        this.handlers.get(eventName)!.add(handler);

        return () => {
            this.handlers.get(eventName)?.delete(handler);
        };
    }

    /**
     * Subscribe to an event once
     */
    once(eventName: string, handler: EventHandler): () => void {
        const wrappedHandler = (data: any) => {
            handler(data);
            this.handlers.get(eventName)?.delete(wrappedHandler);
        };
        return this.on(eventName, wrappedHandler);
    }

    disconnect() {
        if (this.reconnectTimer) {
            clearTimeout(this.reconnectTimer);
            this.reconnectTimer = null;
        }
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }
}

// Singleton instance
export const wsClient = new WebSocketClient();
