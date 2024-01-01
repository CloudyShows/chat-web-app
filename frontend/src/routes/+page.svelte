<!-- +page.svelte -->
<script>
	import ChatRoom from '$lib/components/ChatRoom.svelte';
	import UserList from '$lib/components/UserList.svelte';
	import MessageInput from '$lib/components/MessageInput.svelte';
	import UsernameInput from '$lib/components/UsernameInput.svelte';
	import { darkMode, usernameStore } from '$lib/stores.js';
	import { getUsername } from '$lib/api.js';
	import Spinner from '$lib/components/Spinner.svelte';
	import { initializeWebSocket, disconnectWebSocket } from '$lib/websocket.js';
	import { onMount, onDestroy } from 'svelte';
	import { browser } from '$app/environment';

	let state = {
		isLoading: true,
		isError: false,
		errorMessage: 'Error loading chat.'
	};

	let hasUsername = false;
	let showUsernameForm = false;

	// Subscribe to the dark mode store and update the theme
	if (browser) {
		darkMode.subscribe((value) => {
			const theme = value ? 'light' : 'dark';
			document.documentElement.classList.toggle('dark', value);
			localStorage.setItem('theme', theme);
		});
	}

	// Immediately fetch the username
	getUsername()
		.then((name) => {
			usernameStore.set(name);
		})
		.catch((error) => {
			console.error('Error fetching username:', error);
		});

	// Lifecycle hooks
	onMount(async () => {
		if (browser) {
			try {
				await initializeWebSocket(); // Corrected here

				// Subscribe to the username store
				usernameStore.subscribe((value) => {
					if (value) {
						hasUsername = true;
					} else {
						showUsernameForm = true;
					}
				});

				const savedTheme = localStorage.getItem('theme');
				darkMode.set(savedTheme === 'dark');
				state.isLoading = false;
			} catch (error) {
				console.log(error);
				state.isError = true;
				state.errorMessage = error.message;
			}
		}
	});

	onDestroy(() => {
		if (browser) {
			disconnectWebSocket();
		}
	});
</script>

<div class="flex flex-col h-screen bg-white dark:bg-gray-900">
	{#if state.isLoading}
		<!-- Loading indicator -->
		<div class="flex items-center justify-center h-full">
			<Spinner size="5" />
		</div>
	{:else if state.isError}
		<p class="text-red-500">{state.errorMessage}</p>
	{:else if hasUsername}
		<div class="flex h-full">
			<div class="flex flex-col flex-grow">
				<ChatRoom />
				<MessageInput />
			</div>
			<UserList />
		</div>
	{:else}
		<!-- Show the form to enter or change the username -->
		<UsernameInput on:close={() => (showUsernameForm = false)} />
	{/if}
</div>

<style>
	.flex .h-full {
		max-height: 100vh; /* This ensures that the chat does not exceed the viewport */
	}
</style>
