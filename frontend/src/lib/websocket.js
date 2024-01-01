// websocket.js
import * as api from '$lib/api.js';
import { SERVER_IP } from '$lib/constants';
import { connectedUsersStore, messagesStore, usernameStore } from '$lib/stores.js';

let ws;

export function initializeWebSocket(attempt = 0) {
	ws = new WebSocket(`ws://${SERVER_IP}/ws`);

	ws.addEventListener('open', handleWebSocketOpen);
	ws.addEventListener('close', () => handleWebSocketClose(attempt));
	ws.addEventListener('error', handleWebSocketError);
	ws.addEventListener('message', handleWebSocketMessage);
}

function handleWebSocketOpen(event) {
    // // Fetch the username from the backend when the WebSocket connection is opened
    // api.getUsername().then((username) => {
    //     if (username) {
    //         usernameStore.set(username); // Set the username in the store
    //     }
    // }).catch((error) => {
    //     console.error('Error fetching username:', error);
    // });
}

function handleWebSocketClose(attempt = 0) {
	console.warn('WebSocket closed:', event);
	if (attempt < 3) {
		setTimeout(() => initializeWebSocket(attempt + 1), 1000 * (attempt + 1)); // exponential backoff
	}
}

function handleWebSocketError(event) {
	console.error('WebSocket error:', event);
}

function handleWebSocketMessage(event) {
    const data = event.data;
    // console.log('Received message:', data);
    if (typeof data === 'string') {
        try {
            const message = JSON.parse(data);
            // console.log('Received message object:', message);

            switch (message.type) {
                case 'users':
                    // console.info('Received users:', message.users);
                    connectedUsersStore.set(message.users);
                    break;
                case 'message':
                    messagesStore.update((messages) => [...messages, message]);
                    break;
                case 'error':
                    console.error('Server error:', message.error);
                    break;
                case 'success':
                    // console.log('Server success:', message.message);
                    break;
                default:
                    console.warn('Unhandled message type:', message.type);
            }
        } catch (error) {
            console.error('Error parsing WebSocket message:', error);
        }
    } else {
        console.warn('Received non-string message:', data);
    }
}

export function disconnectWebSocket() {
	if (ws) {
		ws.removeEventListener('message', handleWebSocketMessage);
		ws.removeEventListener('error', handleWebSocketError);
		ws.removeEventListener('close', handleWebSocketClose);
		ws.close();
	}
}

export function sendMessage(text, username) {
	if (ws && ws.readyState === WebSocket.OPEN) {
		ws.send(JSON.stringify({ type: 'message', text, username }));
	} else {
		console.warn('WebSocket is not open. Cannot send message.');
		console.warn('WebSocket readyState:', ws.readyState);
	}
}
