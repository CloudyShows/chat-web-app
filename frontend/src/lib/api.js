// lib/api.js

/**
 * Fetches the username from the server.
 * The username is retrieved from local storage and sent as a parameter in the request.
 * If the fetch is successful, the username is returned.
 * If the fetch fails, an error is logged to the console.
 *
 * @returns {Promise<string>} A promise that resolves to the username.
 */
export async function getUsername() {
	try {
		const savedUsername = localStorage.getItem('username');
		const response = await fetch(
			`http://10.10.0.2:3001/getUsername?username=${encodeURIComponent(savedUsername || '')}`
		);
		const data = await response.json();
		return data.username;
	} catch (error) {
		console.error('Error fetching username:', error);
	}
}
