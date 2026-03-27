export function connectWebSocket(onMessage, onOpen, onClose) {
	const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
	const wsUrl = `${protocol}//${window.location.host}/ws`;

	const ws = new WebSocket(wsUrl);

	ws.onopen = () => {
		console.log('WebSocket connected');
		if (onOpen) onOpen();
	};

	ws.onmessage = (event) => {
		try {
			const data = JSON.parse(event.data);
			if (onMessage) onMessage(data);
		} catch (e) {
			console.error('Failed to parse WebSocket message', e);
		}
	};

	ws.onclose = () => {
		console.log('WebSocket disconnected');
		if (onClose) onClose();
		setTimeout(() => connectWebSocket(onMessage, onOpen, onClose), 5000);
	};

	ws.onerror = (error) => {
		console.error('WebSocket error', error);
	};

	return ws;
}