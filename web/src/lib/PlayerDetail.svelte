<script>
	import { onMount } from 'svelte';
	import Calendar from './Calendar.svelte';

	export let playerName;
	export let onClose;

	let sessions = [];
	let totalPlaytime = 0;
	let recentEvents = [];
	let currentWeekStart = null;

	onMount(async () => {
		const today = new Date();
		const dayOfWeek = today.getDay();
		const sunday = new Date(today);
		sunday.setDate(today.getDate() - dayOfWeek);
		currentWeekStart = sunday.toISOString().split('T')[0];
		await loadData();
	});

	async function loadData() {
		try {
			const weekStart = currentWeekStart || getSundayOfCurrentWeek();
			const [sessionsRes, eventsRes] = await Promise.all([
				fetch(`/api/sessions?player=${encodeURIComponent(playerName)}&start=${weekStart}`),
				fetch(`/api/events?player=${encodeURIComponent(playerName)}&limit=50`)
			]);
			sessions = await sessionsRes.json();
			recentEvents = await eventsRes.json();
			totalPlaytime = sessions.reduce((sum, s) => sum + (s.duration_seconds || 0), 0);
		} catch (e) {
			console.error('Failed to load player data', e);
		}
	}

	function getSundayOfCurrentWeek() {
		const today = new Date();
		const dayOfWeek = today.getDay();
		const sunday = new Date(today);
		sunday.setDate(today.getDate() - dayOfWeek);
		return sunday.toISOString().split('T')[0];
	}

	async function navigateWeek(newWeekStart) {
		currentWeekStart = newWeekStart;
		try {
			const res = await fetch(`/api/sessions?player=${encodeURIComponent(playerName)}&start=${newWeekStart}`);
			sessions = await res.json();
		} catch (e) {
			console.error('Failed to load sessions', e);
		}
	}

	function formatDuration(seconds) {
		const hours = Math.floor(seconds / 3600);
		const mins = Math.floor((seconds % 3600) / 60);
		if (hours > 0) return `${hours}h ${mins}m`;
		return `${mins}m`;
	}

	function formatDate(timestamp) {
		return new Date(timestamp).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function getEventIcon(type) {
		return type === 'join' ? '→' : '←';
	}

	function getEventColor(type) {
		return type === 'join' ? '#00ba7c' : '#f4212b';
	}
</script>

<svelte:window on:keydown={(e) => e.key === 'Escape' && onClose?.()}/>

<div class="modal-backdrop" on:click={onClose} on:keydown={(e) => e.key === 'Escape' && onClose?.()}>
	<div class="modal" on:click|stopPropagation>
		<div class="modal-header">
			<h2>{playerName}</h2>
			<button class="close-btn" on:click={onClose}>×</button>
		</div>

		<div class="modal-body">
			<div class="total-playtime">
				<span class="label">Total Playtime</span>
				<span class="value">{formatDuration(totalPlaytime)}</span>
			</div>

			<Calendar {sessions} {playerName} weekStart={currentWeekStart} onNavigate={navigateWeek} />

			<div class="recent-sessions">
				<h3>Recent Sessions</h3>
				{#if recentEvents.length === 0}
					<p class="empty">No sessions recorded yet</p>
				{:else}
					<ul class="session-list">
						{#each recentEvents.slice(0, 10) as session}
							<li class="session-item" style="--event-color: {getEventColor(session.event_type)}">
								<div class="session-icon">{getEventIcon(session.event_type)}</div>
								<div class="session-info">
									<span class="session-type">{session.event_type === 'join' ? 'Joined' : 'Left'}</span>
									<span class="session-time">{formatDate(session.timestamp)}</span>
								</div>
							</li>
						{/each}
					</ul>
				{/if}
			</div>
		</div>
	</div>
</div>

<style>
	.modal-backdrop {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.7);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
	}

	.modal {
		background: #16181c;
		border-radius: 16px;
		width: 90%;
		max-width: 700px;
		max-height: 90vh;
		overflow-y: auto;
		display: flex;
		flex-direction: column;
	}

	.modal-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 16px 20px;
		border-bottom: 1px solid #2f3336;
	}

	h2 {
		margin: 0;
		font-size: 1.25rem;
		color: #e7e9ea;
	}

	.close-btn {
		background: none;
		border: none;
		font-size: 1.5rem;
		color: #71767b;
		cursor: pointer;
		padding: 0;
		line-height: 1;
	}

	.close-btn:hover {
		color: #e7e9ea;
	}

	.modal-body {
		padding: 20px;
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.total-playtime {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 16px;
		background: linear-gradient(135deg, #1d9bf0 0%, #7856ff 100%);
		border-radius: 12px;
	}

	.total-playtime .label {
		font-size: 0.75rem;
		text-transform: uppercase;
		opacity: 0.9;
	}

	.total-playtime .value {
		font-size: 2rem;
		font-weight: 700;
		margin-top: 4px;
	}

	.recent-sessions h3 {
		margin: 0 0 12px 0;
		font-size: 1rem;
		font-weight: 600;
	}

	.empty {
		color: #71767b;
		text-align: center;
		padding: 20px;
	}

	.session-list {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.session-item {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 10px;
		background: #191c20;
		border-radius: 8px;
	}

	.session-icon {
		width: 28px;
		height: 28px;
		border-radius: 50%;
		background: color-mix(in srgb, var(--event-color) 15%, transparent);
		color: var(--event-color);
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: bold;
	}

	.session-info {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.session-type {
		font-weight: 600;
		font-size: 0.875rem;
	}

	.session-time {
		font-size: 0.75rem;
		color: #71767b;
	}
</style>