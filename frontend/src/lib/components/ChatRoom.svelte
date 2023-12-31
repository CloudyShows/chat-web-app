<!-- ChatRoom.svelte -->
<script>
	import { messagesStore } from '$lib/stores.js';
	import { formatTimestamp } from '$lib/utils.js';

	// Function to determine if a message should start a new group based on time
	function shouldStartNewGroup(message, index) {
		const previousMessage = $messagesStore[index - 1];
		if (!previousMessage) return true; // First message in the list
        
		const currentMessageTime = new Date(message.timestamp);
		const previousMessageTime = new Date(previousMessage.timestamp);
		const timeDifference = currentMessageTime - previousMessageTime;

		// Start new group if messages are more than 1 minute apart
        console.log(`Time difference for message at index ${index}: ${timeDifference}`);
		return timeDifference > 1 * 60 * 1000;
	}
</script>

<div class="chat-room bg-gray-800 text-white">
	<div class="header p-4 bg-gray-900 border-b border-gray-700">Chat Room</div>
	<div class="messages-display-area flex flex-col p-4 space-y-2">
		{#each $messagesStore as message, index (index)}
			{#if shouldStartNewGroup(message, index)}
				<div class="timestamp-group text-sm text-gray-500 mt-2 mb-1 text-center">
					{formatTimestamp(message.timestamp)}
				</div>
			{/if}
			<div class="message flex items-end space-x-2">
				<div class="text rounded-md px-3 py-1" style="background-color: rgba(255, 255, 255, 0.1);">
					{#if message.username !== 'You'}
						<div class="username text-blue-400 font-semibold">{message.username}</div>
					{/if}
					<div class="text-content text-gray-300">{message.text}</div>
				</div>
			</div>
		{/each}
	</div>
</div>

<style>
	.chat-room {
		position: relative;
		padding-bottom: 70px;
	}

	.message {
		max-width: 75%;
	}
	.text-content {
		word-wrap: break-word;
	}
	.messages-display-area {
		flex-grow: 1;
		overflow-y: auto;
		margin-bottom: 70px; /* Height of the input area to prevent overlap */
		max-height: calc(100% - 70px);
	}

	.message:last-child {
		margin-bottom: 1rem;
	}
</style>
