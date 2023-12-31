<!-- ChatRoom.svelte -->
<script>
	import { sendMessage } from '$lib/websocket.js';
	import { usernameStore } from '$lib/stores.js';

	let newMessage = '';
	let username = '';

	usernameStore.subscribe((value) => {
		username = value;
	});

	function handleSend() {
		if (newMessage.trim() !== '') {
			sendMessage(newMessage, username);
			newMessage = ''; // Clear the input after sending
		}
	}

	// Function to handle key press events
	function handleKeyPress(event) {
		if (event.key === 'Enter') {
			handleSend();
		}
	}
</script>

<div class="message-input-container p-4 bg-gray-800">
	<input
		class="message-input w-full p-2 rounded-lg bg-gray-700 text-white"
		type="text"
		placeholder="Type a message..."
		bind:value={newMessage}
		on:keypress={handleKeyPress} />
	<button class="send-button px-4 py-2 ml-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600" on:click={handleSend}>
		Send
	</button>
</div>

<style>
	.message-input-container {
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;
	}
</style>
