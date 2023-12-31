// api.js
import axios from 'axios';

export async function getUsername() {
    try {
        const savedUsername = localStorage.getItem('username');
        const response = await axios.get(
            `http://10.10.0.2:3001/getUsername?username=${encodeURIComponent(savedUsername || '')}`
        );

        return response.data.username;
    } catch (error) {
        console.error('Error fetching username:', error);
        throw error;
    }
}

export async function changeUsername(newUsername) {
    try {
        await axios.post(`http://10.10.0.2:3001/changeUsername`, {
            newUsername
        }, {
            headers: { 'Content-Type': 'application/json' }
        });
    } catch (error) {
        console.error('Error changing username:', error);
        throw error;
    }
}