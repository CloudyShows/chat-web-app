<script>
	//UsernameInput.svelte
	import { createEventDispatcher } from 'svelte';
	import * as api from '$lib/api.js';
	import { usernameStore } from '$lib/stores.js';

	export let currentUsername = '';
	let username = currentUsername;

	const dispatch = createEventDispatcher();

	function closeUsernameForm() {
		dispatch('close');
	}

	function handleOverlayClick(event) {
		if (event.target === event.currentTarget) {
			closeUsernameForm();
		}
	}

	async function changeUsername() {
		if (username.trim() !== '' && username.trim() !== currentUsername) {
			try {
				await api.changeUsername(username);

				usernameStore.set(username); // Update Svelte store
				currentUsername = username; // Update current username

				// Close the modal
				closeUsernameForm();
			} catch (error) {
				console.error('Failed to update username', error);
			}
		}
	}
</script>

<div
	class="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-50 p-4"
	on:click={handleOverlayClick}>
	<div class="bg-gray-800 p-4 rounded-md">
		<input
			class="p-2 w-full rounded-md border bg-gray-700 border-gray-600 text-white focus:outline-none focus:border-blue-500 mb-2"
			type="text"
			placeholder="New username..."
			bind:value={username} />
		<div class="flex justify-end space-x-2">
			<button
				class="transition-colors duration-200 px-4 py-2 bg-gray-500 text-white rounded-md hover:bg-gray-600 focus:outline-none focus:bg-gray-600"
				on:click={closeUsernameForm}>
				Cancel
			</button>
			<button
				class="transition-colors duration-200 px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:bg-blue-600"
				on:click={changeUsername}>
				Update
			</button>
		</div>
	</div>
</div>
