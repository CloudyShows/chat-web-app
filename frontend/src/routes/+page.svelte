<!-- +page.svelte -->
<script>
	import { onMount, onDestroy } from 'svelte';
	import { browser } from '$app/environment';
	import {
		initializeWebSocket,
		disconnectWebSocket,
		messagesStore,
		usernameStore,
		connectedUsersStore,
		sendMessage,
		fetchUsername
	} from '$lib/messagesStore.js';

	let state = {
		isLoading: true,
		isError: false,
		errorMessage: 'Error loading chat.'
	};

	let newMessage = '';
	let username = '';
	let hasUsername = false;
	let showUsernameForm = false;

	function handleSend() {
		if (newMessage.trim() !== '') {
			sendMessage(newMessage);
			newMessage = '';
		}
	}
	function setUsername() {
		if (username.trim() !== '') {
			localStorage.setItem('username', username); // Save to local storage
			usernameStore.set(username);
			fetchUsername(); // Update the backend
			hasUsername = true;
		}
	}

	function changeUsername() {
		if (username.trim() !== '') {
			localStorage.setItem('username', username); // Save to local storage
			usernameStore.set(username);
			fetchUsername(); // Update the backend
			showUsernameForm = false; // Hide the username form
		}
	}

	function formatTimestamp(timestampStr) {
		const messageDate = new Date(timestampStr);
		const today = new Date();

		if (messageDate.toDateString() === today.toDateString()) {
			// If the message is from today, show the time
			return messageDate.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
		} else {
			// Otherwise, show the date
			return messageDate.toLocaleDateString();
		}
	}

	onMount(async () => {
		if (browser) {
			try {
				await initializeWebSocket(); // Initialize WebSocket connection
				const savedUsername = localStorage.getItem('username');
				if (savedUsername) {
					username = savedUsername;
					usernameStore.set(username);
					hasUsername = true;
				}
				await fetchUsername(); // Fetch username
				state.isLoading = false; // Set loading to false after initialization
			} catch (error) {
				state.isError = true;
				state.errorMessage = error.message;
			}
		}
	});
	onDestroy(() => {
		if (browser) {
			disconnectWebSocket(); // Disconnect WebSocket connection
		}
	});
</script>

<div class="flex flex-col h-screen bg-white">
	{#if state.isLoading}
		<div class="flex items-center justify-center h-full">
			<svg
				class="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
				<path
					class="opacity-75"
					fill="currentColor"
					d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
			</svg>
		</div>
	{:else if hasUsername}
		<!-- Connected Users List -->
		<div class="p-4 bg-white border-b border-gray-200">
			<div class="flex justify-between items-center">
				<h2 class="text-lg font-semibold mb-2 text-gray-800">Connected Users</h2>
				<button
					class="transition-colors duration-200 px-4 py-2 bg-gray-300 text-gray-700 rounded-md hover:bg-gray-400 focus:outline-none focus:bg-gray-400"
					on:click={() => (showUsernameForm = true)}>
					Change Username
				</button>
			</div>
			<ul class="list-none space-y-1">
				{#each $connectedUsersStore as user}
					<li class="text-sm">
						<span class="font-medium text-gray-700">
							{user}
							{#if user === username}
								<span class="text-blue-500 ml-1">(me)</span>
							{/if}
						</span>
					</li>
				{/each}
			</ul>
		</div>

		<!-- Chat interface -->
		<div class="flex flex-col flex-1">
			<!-- Chat Messages -->
			<div class="flex-1 overflow-y-auto p-4 space-y-4">
				{#each $messagesStore || [] as message (message.id)}
					<div class={`flex ${message.username === username ? 'justify-end' : 'justify-start'} items-end space-x-2`}>
						<div
							class={`max-w-xs md:max-w-md lg:max-w-lg xl:max-w-xl px-4 py-2 my-1 rounded-full shadow ${
								message.username === username ? 'bg-blue-500 text-white' : 'bg-gray-200 text-gray-700'
							}`}>
							{message.text}
							<span class="block text-xs text-gray-200 mt-1">{formatTimestamp(message.timestamp)}</span>
						</div>
					</div>
				{/each}
			</div>

			<!-- Message input with adjusted spacing -->
			<div class="p-6 bg-gray-50">
				<div class="relative flex items-center">
					<input
						class="flex-grow p-4 rounded-full border border-gray-300 bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 text-lg"
						type="text"
						placeholder="Type a message..."
						bind:value={newMessage}
						on:keyup={(e) => e.key === 'Enter' && handleSend()} />
					<button
						class="absolute right-4 px-6 py-3 bg-blue-500 text-white rounded-full hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50 active:bg-blue-700"
						on:click={handleSend}>
						Send
					</button>
				</div>
			</div>
		</div>

		{#if showUsernameForm}
			<!-- Username form overlay -->
			<div class="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-50 p-4">
				<div class="bg-white p-4 rounded-md">
					<input
						class="p-2 w-full rounded-md border bg-gray-100 focus:outline-none focus:border-blue-300 mb-2"
						type="text"
						placeholder="New username..."
						bind:value={username} />
					<button
						class="transition-colors duration-200 px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:bg-blue-600"
						on:click={changeUsername}>
						Update
					</button>
				</div>
			</div>
		{/if}
	{:else}
		<!-- Username input -->
		<div class="flex flex-col justify-center items-center h-full space-y-2">
			<input
				class="p-2 w-1/2 rounded-md border bg-gray-100 focus:outline-none focus:border-blue-300"
				type="text"
				placeholder="Enter your username..."
				bind:value={username}
				on:keyup={(e) => e.key === 'Enter' && setUsername()} />
			<button
				class="transition-colors duration-200 px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:bg-blue-600"
				on:click={setUsername}>
				Set Username
			</button>
		</div>
	{/if}
</div>
