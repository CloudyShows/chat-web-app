import { writable } from 'svelte/store';

export const darkMode = writable(false);
export const messagesStore = writable([]);
export const usernameStore = writable('');
export const connectedUsersStore = writable([]);