<!-- UserList.svelte -->
<script>
    import { connectedUsersStore, usernameStore } from '$lib/stores.js';
	import Spinner from '$lib/components/Spinner.svelte';
	
    $: username = $usernameStore;
    $: isLoadingUsername = username === ''; // Loading state when username is not fetched yet
</script>

{#if isLoadingUsername}
	<!-- Loading indicator -->
	<div class="flex items-center justify-center h-full">
		<Spinner size="5" />
	</div>
{:else}
	<aside class="user-list bg-gray-800 text-white p-4">
		<h2 class="text-lg font-semibold mb-4">Connected Users</h2>
		<ul class="list-none space-y-1">
			{#each $connectedUsersStore as user}
				<li class="text-sm flex items-center justify-between">
					<span>{user}</span>
					{#if user === username}
						<span class="text-blue-500 ml-1">(me)</span>
					{/if}
				</li>
			{/each}
		</ul>
	</aside>
{/if}

<style>
	.user-list {
		background-color: rgba(255, 255, 255, 0.05); /* Semi-transparent background for the user list */
		border-left: 1px solid rgba(255, 255, 255, 0.1); /* A subtle border to the left */
	}
</style>
