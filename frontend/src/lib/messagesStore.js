// messagesStore.js
import { writable } from 'svelte/store';

export const messagesStore = writable([]);
export const usernameStore = writable('');

let ws;

export function initializeWebSocket() {
	ws = new WebSocket('ws://10.10.0.2:3001/ws');

	ws.addEventListener('message', function (event) {
		const receivedData = JSON.parse(event.data);

		// Check if the received data is an array (chat history) or an object (new message)
		if (Array.isArray(receivedData)) {
			messagesStore.set(receivedData); // Set the entire chat history
		} else {
			messagesStore.update((messages) => [...messages, receivedData]); // Add the new message
		}
	});
}

export function fetchUsername() {
	const savedUsername = localStorage.getItem('username');
	fetch(`http://10.10.0.2:3001/getUsername?username=${savedUsername || ''}`)
		.then((response) => response.json())
		.then((data) => {
			usernameStore.set(data.username);
		})
		.catch((error) => {
			console.error('Error fetching username:', error);
		});
}

export function sendMessage(text) {
	usernameStore.subscribe((username) => {
		console.log('Sending message with username:', username); // Debug line
		if (ws) {
			ws.send(JSON.stringify({ type: 'message', text, username }));
		} else {
			console.warn('WebSocket is not initialized.');
		}
	})();
}
