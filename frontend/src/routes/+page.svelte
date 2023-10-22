<script>
	import { messagesStore, usernameStore, sendMessage } from '$lib/messagesStore.js';
	import { onMount } from 'svelte';

	let newMessage = '';
	let username = '';
	let hasUsername = false;
	let showUsernameForm = false;

	function handleSend() {
		if (newMessage.trim() !== '') {
			sendMessage(newMessage); // Update this line
			newMessage = '';
		}
	}
	function setUsername() {
		if (username.trim() !== '') {
			localStorage.setItem('username', username); // Save to local storage
			usernameStore.set(username);
			hasUsername = true;
		}
	}

	function changeUsername() {
		if (username.trim() !== '') {
			localStorage.setItem('username', username); // Save to local storage
			usernameStore.set(username);
			showUsernameForm = false; // Hide the username form
		}
	}

	onMount(() => {
		const savedUsername = localStorage.getItem('username');
		if (savedUsername) {
			username = savedUsername;
			usernameStore.set(username);
			hasUsername = true;
		}
	});
</script>

<div class="flex flex-col h-screen bg-gray-100">
	{#if hasUsername}
		<!-- Chat interface -->
		<div class="flex flex-col flex-1 p-6">
			<div class="flex flex-col h-full p-4 bg-white rounded-lg shadow-lg overflow-y-auto">
				{#each $messagesStore || [] as message}
					<div class="flex items-start mb-4">
						<div class="flex items-end">
							<span class="text-sm font-medium text-gray-800 mr-2">{message.username}:</span>
						</div>
						<div class="flex flex-col space-y-2 text-xs max-w-xs mx-2 order-2 items-start">
							<div>
								<span class="px-4 py-2 rounded-lg inline-block rounded-bl-none bg-blue-600 text-white">
									{message.text}
								</span>
							</div>
						</div>
					</div>
				{/each}
			</div>

			<div class="mt-4 flex items-center">
				<input
					class="flex-1 p-2 rounded-l border"
					type="text"
					placeholder="Type a message..."
					bind:value={newMessage}
					on:keyup={(e) => e.key === 'Enter' && handleSend()} />
				<button class="px-4 py-2 bg-blue-500 text-white rounded-l-none rounded-r" on:click={handleSend}> Send </button>
				<button class="ml-2 px-4 py-2 bg-gray-400 text-white rounded" on:click={() => (showUsernameForm = true)}>
					Change Username
				</button>
			</div>
		</div>

		{#if showUsernameForm}
			<!-- Username form overlay -->
			<div class="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-50">
				<div class="bg-white p-4 rounded">
					<input class="p-2 rounded border mb-2" type="text" placeholder="New username..." bind:value={username} />
					<button class="px-4 py-2 bg-blue-500 text-white rounded" on:click={changeUsername}> Update </button>
				</div>
			</div>
		{/if}
	{:else}
		<!-- Username input -->
		<div class="flex flex-col justify-center items-center h-full">
			<input
				class="p-2 rounded border mb-2"
				type="text"
				placeholder="Enter your username..."
				bind:value={username}
				on:keyup={(e) => e.key === 'Enter' && setUsername()} />
			<button class="px-4 py-2 bg-blue-500 text-white rounded" on:click={setUsername}> Set Username </button>
		</div>
	{/if}
</div>
