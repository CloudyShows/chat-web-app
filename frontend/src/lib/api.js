export async function getUsername() {
	try {
		const savedUsername = localStorage.getItem('username');
		const response = await fetch(
			`http://10.10.0.2:3001/getUsername?username=${encodeURIComponent(savedUsername || '')}`
		);

		if (!response.ok) {
			throw new Error('Failed to fetch username from the server.');
		}

		const data = await response.json();
		return data.username;
	} catch (error) {
		console.error('Error fetching username:', error);
		throw error;
	}
}

export async function changeUsername(newUsername) {
    try {
        const response = await fetch(`http://10.10.0.2:3001/changeUsername`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ newUsername })
        });

        if (!response.ok) {
            throw new Error('Failed to change username on the server.');
        }
    } catch (error) {
        console.error('Error changing username:', error);
        throw error;
    }
}

