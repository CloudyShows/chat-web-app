export function formatTimestamp(timestampStr) {
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