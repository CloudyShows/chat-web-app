import { writable } from 'svelte/store';
import * as api from '$lib/api.js';

export const messagesStore = writable([]);
export const usernameStore = writable('');
export const connectedUsersStore = writable([]);

let ws;

function handleWebSocketMessage(event) {
    const data = JSON.parse(event.data);

    switch (data.type) {
        case 'users':
            connectedUsersStore.set(data.users);
            break;
        case 'message':
            messagesStore.update(messages => [...messages, data]);
            break;
        default:
            console.warn('Unhandled message type:', data.type);
    }
}

export function initializeWebSocket() {
    ws = new WebSocket('ws://10.10.0.2:3001/ws');
    ws.addEventListener('message', handleWebSocketMessage);

    api.getUsername().then(username => {
        if (username) {
            usernameStore.set(username);
        }
    });
}

export function disconnectWebSocket() {
    if (ws) {
        ws.removeEventListener('message', handleWebSocketMessage);
        ws.close();
    }
}

export function sendMessage(text, username) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: 'message', text, username }));
    } else {
        console.warn('WebSocket is not open. Cannot send message.');
    }
}
