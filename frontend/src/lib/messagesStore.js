// messagesStore.js
import { writable } from 'svelte/store';
import * as api from '$lib/api.js';

export const messagesStore = writable([]);
export const usernameStore = writable('');
export const connectedUsersStore = writable([]);

let ws;

/**
 * Initializes the WebSocket connection and sets up event listeners.
 * Fetches the username using the api.getUsername function and sets it in the usernameStore.
 */
export function initializeWebSocket() {
	ws = new WebSocket('ws://10.10.0.2:3001/ws');

	ws.addEventListener('message', function (event) {
		const receivedData = JSON.parse(event.data);

		// Handle user list updates
		if (receivedData.type === 'users') {
			connectedUsersStore.set(receivedData.users);
			return;
		}

		// Handle chat messages
		if (receivedData.type === 'message') {
			messagesStore.update((messages) => [...messages, receivedData]);
			return;
		}

		// Handle chat history
		if (Array.isArray(receivedData)) {
			messagesStore.set(receivedData);
			return;
		}
	});

	// Use api.getUsername to fetch the username
	api.getUsername().then((username) => {
		usernameStore.set(username);
	});
}

/**
 * Closes the WebSocket connection if it's open.
 */
export function disconnectWebSocket() {
	if (ws) {
		ws.close();
	}
}

/**
 * Sends a message over the WebSocket connection.
 * The message is sent as a JSON string with a type of 'message', the text of the message, and the username.
 * If the WebSocket connection is not open, a warning is logged to the console.
 *
 * @param {string} text - The text of the message to send.
 * @param {string} username - The username of the sender.
 */
export function sendMessage(text, username) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: 'message', text, username }));
    } else {
        console.warn('Cannot send message, WebSocket is not open');
    }
}