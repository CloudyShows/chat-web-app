<!-- +page.svelte -->
<script>
	import ChatRoom from '$lib/components/ChatRoom.svelte';
	import UserList from '$lib/components/UserList.svelte';
	import MessageInput from '$lib/components/MessageInput.svelte';
	import UsernameInput from '$lib/components/UsernameInput.svelte';
	import Spinner from '$lib/components/Spinner.svelte';
	import { initializeWebSocket, disconnectWebSocket } from '$lib/websocket.js';
	import { usernameStore } from '$lib/stores.js';
	import { onMount, onDestroy } from 'svelte';
	import { browser } from '$app/environment';

	let state = {
		isLoading: true,
		isError: false,
		errorMessage: 'Error loading chat.'
	};

	let hasUsername = false;
	let showUsernameForm = false;

	function toggleUsernameForm() {
		showUsernameForm = !showUsernameForm;
	}

	// Lifecycle hooks
	onMount(async () => {
		if (browser) {
			try {
				console.log('Initializing websocket...');
				await initializeWebSocket();
				console.log('Websocket initialized.');
				usernameStore.subscribe((value) => {
					hasUsername = !!value;
					showUsernameForm = !value;
				});

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

<div class="flex flex-col h-screen bg-gray-900">
	<!-- Set the dark mode background color -->
	{#if state.isLoading}
		<!-- Loading indicator -->
		<div class="flex items-center justify-center h-full">
			<Spinner size="5" />
		</div>
	{:else if state.isError}
		<p class="text-red-500">{state.errorMessage}</p>
	{/if}
	{#if hasUsername}
		<div class="flex h-full">
			<div class="flex flex-col flex-grow">
				<ChatRoom />
				<MessageInput />
			</div>
			<UserList />
			<button class="change-username-btn" on:click={toggleUsernameForm}> Change Username </button>
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
	.change-username-btn {
		position: fixed;
		bottom: 20px;
		right: 20px;
		padding: 10px 15px;
		background-color: #5865f2;
		color: white;
		border: none;
		border-radius: 5px;
		cursor: pointer;
		transition: background-color 0.3s;
	}

	.change-username-btn:hover {
		background-color: #4e5dcf;
	}
</style>
