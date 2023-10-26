<!-- +layout.svelte -->
<script>
	import '../app.postcss';
	import { browser } from '$app/environment';
	import { initializeWebSocket, fetchUsername } from '$lib/messagesStore.js';

	if (browser) {
		initializeWebSocket();
		fetchUsername(); // New line
	}
</script>

<main>
	<slot />
</main>
