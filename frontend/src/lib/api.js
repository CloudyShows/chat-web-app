// api.js
import axios from 'axios';
import { SERVER_IP } from '$lib/constants';

export async function getUsername() {
    try {
        const response = await axios.get(
            `http://${SERVER_IP}/getUsername`
        );

        return response.data.username;
    } catch (error) {
        console.error('Error fetching username:', error);
        throw error;
    }
}

export async function changeUsername(newUsername) {
    try {
        await axios.post(`http://${SERVER_IP}/changeUsername`, {
            newUsername
        }, {
            headers: { 'Content-Type': 'application/json' }
        });
    } catch (error) {
        console.error('Error changing username:', error);
        throw error;
    }
}