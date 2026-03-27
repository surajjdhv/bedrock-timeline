<script>
	import { onMount } from 'svelte';
	import Timeline from '$lib/Timeline.svelte';
	import PlayerList from '$lib/PlayerList.svelte';
	import Stats from '$lib/Stats.svelte';
	import Leaderboard from '$lib/Leaderboard.svelte';
	import PlayerDetail from '$lib/PlayerDetail.svelte';
	import { connectWebSocket } from '$lib/websocket';

	let events = [];
	let players = [];
	let stats = {};
	let connected = false;
	let wsInstance = null;
	let selectedPlayer = null;

	onMount(async () => {
		await loadData();
		wsInstance = connectWebSocket((event) => {
			events = [event, ...events].slice(0, 100);
			if (event.event_type === 'join' || event.event_type === 'leave') {
				loadData();
			}
		}, () => { connected = true; }, () => { connected = false; });
	});

	async function loadData() {
		try {
			const [eventsRes, playersRes, statsRes] = await Promise.all([
				fetch('/api/events?limit=50'),
				fetch('/api/players'),
				fetch('/api/stats')
			]);
			events = await eventsRes.json();
			players = await playersRes.json();
			stats = await statsRes.json();
		} catch (e) {
			console.error('Failed to load data', e);
		}
	}

	function selectPlayer(name) {
		selectedPlayer = name;
	}

	function closePlayer() {
		selectedPlayer = null;
	}
</script>

<main>
	<header>
		<div class="header-left">
			<h1>Bedrock Timeline</h1>
			<div class="connection-status" class:connected>
				{connected ? '● Connected' : '○ Disconnected'}
			</div>
		</div>
		<a href="/calendar" class="calendar-link">View Full Calendar →</a>
	</header>
	
	<div class="dashboard">
		<div class="sidebar">
			<Stats {stats} />
			<Leaderboard {players} />
			<PlayerList {players} onSelect={selectPlayer} />
		</div>
		
		<div class="content">
			<Timeline {events} />
		</div>
	</div>

	{#if selectedPlayer}
		<PlayerDetail playerName={selectedPlayer} onClose={closePlayer} />
	{/if}
</main>

<style>
	:global(body) {
		margin: 0;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
		background: #0f1419;
		color: #e7e9ea;
	}

	main {
		max-width: 1400px;
		margin: 0 auto;
		padding: 20px;
	}

	header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 24px;
		padding-bottom: 16px;
		border-bottom: 1px solid #2f3336;
	}

	.header-left {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	h1 {
		margin: 0;
		font-size: 1.5rem;
		color: #1d9bf0;
	}

	.connection-status {
		font-size: 0.875rem;
		color: #71767b;
	}

	.connection-status.connected {
		color: #00ba7c;
	}

	.calendar-link {
		background: transparent;
		border: 1px solid #1d9bf0;
		color: #1d9bf0;
		padding: 8px 16px;
		border-radius: 8px;
		text-decoration: none;
		font-size: 0.875rem;
		font-weight: 500;
		transition: background 0.2s;
	}

	.calendar-link:hover {
		background: rgba(29, 155, 240, 0.1);
	}

	.dashboard {
		display: grid;
		grid-template-columns: 300px 1fr;
		gap: 24px;
	}

	@media (max-width: 768px) {
		.dashboard {
			grid-template-columns: 1fr;
		}
	}

	.sidebar {
		display: flex;
		flex-direction: column;
		gap: 16px;
	}
</style>