<!-- +page.svelte -->
<script>
	// Components
	import UsernameModal from '$lib/UsernameModal.svelte';

	// Stores
	import { darkMode } from '$lib/stores.js';
	import {
		initializeWebSocket,
		disconnectWebSocket,
		messagesStore,
		usernameStore,
		connectedUsersStore,
		sendMessage
	} from '$lib/messagesStore.js';
	import * as api from '$lib/api.js';
	import { onMount, onDestroy } from 'svelte';
	import { browser } from '$app/environment';

	// Variables
	let state = {
		isLoading: true,
		isError: false,
		errorMessage: 'Error loading chat.'
	};

	let newMessage = '';
	let username = '';
	let hasUsername = false;
	let showUsernameForm = false;

	// Stores
	
	// Subscribe to the usernameStore
	usernameStore.subscribe((value) => {
		username = value;
	});

	// Functions
	function toggleTheme() {
		darkMode.update((value) => {
			const theme = value ? 'light' : 'dark';
			document.documentElement.classList.toggle('dark', !value);
			localStorage.setItem('theme', theme);
			return !value;
		});
	}

	function closeModal() {
		showUsernameForm = false;
	}

	function setUsername() {
		if (username.trim() !== '') {
			localStorage.setItem('username', username); // Save to local storage
			usernameStore.set(username);
			api.getUsername(); // Update the backend
			hasUsername = true;
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

	function handleSend() {
		if (newMessage.trim() !== '') {
			sendMessage(newMessage, username);
			newMessage = '';
		}
	}

	// Lifecycle hooks
	onMount(async () => {
		if (browser) {
			try {
				// Initialize WebSocket connection
				await initializeWebSocket();

				// Initialize dark mode from localStorage
				const savedTheme = localStorage.getItem('theme');
				darkMode.set(savedTheme === 'dark');

				// Apply the dark class if necessary
				darkMode.subscribe((value) => {
					document.documentElement.classList.toggle('dark', value);
				});

				// Check if the user has a username saved
				const savedUsername = localStorage.getItem('username');
				if (savedUsername) {
					username = savedUsername;
					usernameStore.set(username);
					hasUsername = true;
				}
				await api.getUsername(); // Fetch username
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

<div class="flex flex-col h-screen bg-white dark:bg-gray-900">
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
		<div class="p-4 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
			<div class="flex justify-between items-center">
				<h2 class="text-lg font-semibold mb-2 text-gray-800 dark:text-gray-200">Connected Users</h2>
				<button
					class="transition-colors duration-200 px-4 py-2 bg-gray-300 text-gray-700 rounded-md hover:bg-gray-400 focus:outline-none focus:bg-gray-400 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600 dark:focus:bg-gray-600"
					on:click={() => (showUsernameForm = true)}>
					Change Username
				</button>
				<!-- Dark mode toggle button -->
				<div class="theme-toggle">
					<button
						class="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:bg-blue-600 dark:bg-blue-400 dark:hover:bg-blue-500 dark:focus:bg-blue-500"
						on:click={toggleTheme}>
						Toggle Dark Mode
					</button>
				</div>
			</div>
			<ul class="list-none space-y-1">
				{#each $connectedUsersStore as user}
					<li class="text-sm">
						<span class="font-medium text-gray-700 dark:text-gray-300">
							{user}
							{#if user === username}
								<span class="text-blue-500 ml-1">(me)</span>
							{/if}
						</span>
					</li>
				{/each}
			</ul>
		</div>

		<div class="flex flex-col flex-1">
			<!-- Chat Messages -->
			<div class="flex-1 overflow-y-auto p-4 space-y-4">
				{#each $messagesStore || [] as message, index (index)}
					<div class={`flex ${message.username === username ? 'justify-end' : 'justify-start'} items-end space-x-2`}>
						{#if message.username !== username}
							<div class="text-sm text-gray-500">
								{message.username}
							</div>
						{/if}
						<div
							class={`max-w-xs md:max-w-md lg:max-w-lg xl:max-w-xl px-4 py-2 my-1 rounded-full shadow ${
								message.username === username
									? 'bg-blue-500 text-white dark:bg-blue-600'
									: 'bg-gray-200 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
							}`}>
							{message.text}
							<span class="block text-xs mt-1 dark:text-gray-200 text-gray-800"
								>{formatTimestamp(message.timestamp)}</span>
						</div>
					</div>
				{/each}
			</div>
		</div>
		<!-- Message input -->
		<div class="p-6 bg-gray-50 dark:bg-gray-800">
			<div class="relative flex items-center">
				<input
					class="flex-grow p-4 rounded-full border border-gray-300 bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 text-lg dark:bg-gray-700 dark:border-gray-600 dark:text-white dark:focus:ring-blue-400"
					type="text"
					placeholder="Type a message..."
					bind:value={newMessage}
					on:keyup={(e) => e.key === 'Enter' && handleSend()} />
				<button
					class="absolute right-4 px-6 py-3 bg-blue-500 text-white rounded-full hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50 active:bg-blue-700 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-400"
					on:click={handleSend}>
					Send
				</button>
			</div>
		</div>
		{#if showUsernameForm}
			<UsernameModal on:close={closeModal} />
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
