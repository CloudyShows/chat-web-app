<!-- +page.svelte -->
<script>
	import { messagesStore, usernameStore, connectedUsersStore, sendMessage, fetchUsername } from '$lib/messagesStore.js';
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

	onMount(() => {
		const savedUsername = localStorage.getItem('username');
		if (savedUsername) {
			username = savedUsername;
			usernameStore.set(username);
			hasUsername = true;
		}
	});
</script>

<div class="flex flex-col h-screen bg-gradient-to-b from-gray-100 to-gray-300">
	{#if hasUsername}
		<!-- Chat interface -->
		<div class="flex flex-col flex-1 p-6 space-y-4">
			<!-- Chat Messages -->
			<div class="flex flex-col h-full p-4 bg-white rounded-lg shadow-lg overflow-y-auto space-y-4">
				{#each $messagesStore || [] as message}
					<div class="flex items-start mb-4 {message.username === username ? 'justify-end' : ''}">
						{#if message.username !== username}
							<div class="text-xs text-gray-500 mr-2">{message.username}</div>
						{/if}
						<div class="flex flex-col space-y-2 text-sm max-w-xs mx-2 order-2 items-start">
							<div
								class="{message.username === username
									? 'bg-blue-600 text-white'
									: 'bg-gray-300 text-black'} px-5 py-3 rounded-3xl inline-block shadow-md">
								{message.text}
							</div>
							<span class="text-gray-500 text-xs self-end">{formatTimestamp(message.timestamp)}</span>
						</div>
					</div>
				{/each}
			</div>

			<div class="mt-4 flex items-center space-x-2 shadow-inner">
				<input
					class="flex-1 p-2 rounded-l-xl border bg-gray-100 focus:outline-none focus:border-blue-300"
					type="text"
					placeholder="Type a message..."
					bind:value={newMessage}
					on:keyup={(e) => e.key === 'Enter' && handleSend()} />
				<button
					class="transition-colors duration-200 px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:bg-blue-600"
					on:click={handleSend}>
					Send
				</button>
				<button
					class="transition-colors duration-200 px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:bg-blue-600"
					on:click={() => (showUsernameForm = true)}>
					Change Username
				</button>
			</div>
		</div>

		<!-- Connected Users List -->
		<div class="p-4 bg-white rounded-lg shadow-lg mt-4 border border-gray-200">
			<h2 class="text-lg font-semibold mb-2 text-blue-600">Connected Users</h2>
			<ul class="list-inside list-disc space-y-1">
				{#each $connectedUsersStore as user}
					<li class="text-sm">
						<span class="font-medium text-gray-700">
							{user}
							{@html user === username ? '<span class="text-green-500 ml-1">(me)</span>' : ''}
						</span>
					</li>
				{/each}
			</ul>
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
