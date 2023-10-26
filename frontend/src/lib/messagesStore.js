// messagesStore.js
import { writable } from 'svelte/store';

export const messagesStore = writable([]);
export const usernameStore = writable('');
export const connectedUsersStore = writable([]);

let ws;

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

		// Handle chat history (assuming it's an array)
		if (Array.isArray(receivedData)) {
			messagesStore.set(receivedData);
			return;
		}
	});
}


export function fetchUsername() {
    const savedUsername = localStorage.getItem('username');
    console.log('Fetching username with:', savedUsername);  // Debug line
    fetch(`http://10.10.0.2:3001/getUsername?username=${encodeURIComponent(savedUsername || '')}`)
        .then((response) => response.json())
        .then((data) => {
            console.log('Received username:', data.username);  // Debug line
            usernameStore.set(data.username);
        })
        .catch((error) => {
            console.error('Error fetching username:', error);
        });
}

export function sendMessage(text) {
	usernameStore.subscribe((username) => {
		console.log('Sending message with username:', username); // Debug line
		if (ws && ws.readyState === WebSocket.OPEN) {
			ws.send(JSON.stringify({ type: 'message', text, username }));
		} else {
			console.warn('WebSocket is not open.');
		}		
	})();
}
