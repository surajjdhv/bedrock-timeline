<script>
	export let events;

	let page = 0;
	const pageSize = 12;

	$: totalPages = Math.ceil(events.length / pageSize);
	$: paginatedEvents = events.slice(page * pageSize, (page + 1) * pageSize);

	function formatTime(timestamp) {
		const date = new Date(timestamp);
		return date.toLocaleTimeString('en-US', { 
			hour: '2-digit', 
			minute: '2-digit',
			hour12: true 
		});
	}

	function formatDate(timestamp) {
		const date = new Date(timestamp);
		return date.toLocaleDateString('en-US', { 
			month: 'short', 
			day: 'numeric' 
		});
	}

	function getEventIcon(eventType) {
		return eventType === 'join' ? '→' : '←';
	}

	function getEventColor(eventType) {
		return eventType === 'join' ? '#00ba7c' : '#f4212b';
	}

	function prevPage() {
		if (page > 0) page--;
	}

	function nextPage() {
		if (page < totalPages - 1) page++;
	}
</script>

<div class="timeline">
	<h2>Recent Events</h2>
	
	{#if events.length === 0}
		<p class="empty">No events yet. Waiting for player activity...</p>
	{:else}
		<ul class="event-list">
			{#each paginatedEvents as event (event.timestamp + event.player_name + event.event_type)}
				<li class="event" style="--event-color: {getEventColor(event.event_type)}">
					<div class="event-icon">
						{getEventIcon(event.event_type)}
					</div>
					<div class="event-content">
						<span class="player-name">{event.player_name}</span>
						<span class="event-type">{event.event_type === 'join' ? 'joined' : 'left'}</span>
					</div>
					<div class="event-time">
						<span class="date">{formatDate(event.timestamp)}</span>
						<span class="time">{formatTime(event.timestamp)}</span>
					</div>
				</li>
			{/each}
		</ul>

		{#if totalPages > 1}
			<div class="pagination">
				<button class="page-btn" on:click={prevPage} disabled={page === 0}>
					← Prev
				</button>
				<span class="page-info">
					Page {page + 1} of {totalPages}
				</span>
				<button class="page-btn" on:click={nextPage} disabled={page >= totalPages - 1}>
					Next →
				</button>
			</div>
		{/if}
	{/if}
</div>

<style>
	.timeline {
		background: #16181c;
		border-radius: 12px;
		padding: 16px;
	}

	h2 {
		margin: 0 0 16px 0;
		font-size: 1.125rem;
		font-weight: 600;
	}

	.empty {
		color: #71767b;
		text-align: center;
		padding: 40px 20px;
	}

	.event-list {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.event {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 12px;
		background: #191c20;
		border-radius: 8px;
		transition: background-color 0.2s;
	}

	.event:hover {
		background: #1d2025;
	}

	.event-icon {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		background: color-mix(in srgb, var(--event-color) 15%, transparent);
		color: var(--event-color);
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1rem;
		font-weight: bold;
	}

	.event-content {
		flex: 1;
	}

	.player-name {
		font-weight: 600;
		color: #e7e9ea;
	}

	.event-type {
		color: #71767b;
		margin-left: 6px;
	}

	.event-time {
		text-align: right;
		font-size: 0.875rem;
		color: #71767b;
	}

	.date {
		display: block;
	}

	.time {
		display: block;
		opacity: 0.7;
	}

	.pagination {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 12px;
		margin-top: 16px;
		padding-top: 16px;
		border-top: 1px solid #2f3336;
	}

	.page-btn {
		background: transparent;
		border: 1px solid #2f3336;
		color: #e7e9ea;
		padding: 6px 12px;
		border-radius: 6px;
		font-size: 0.875rem;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.page-btn:hover:not(:disabled) {
		background: #2f3336;
	}

	.page-btn:disabled {
		opacity: 0.4;
		cursor: not-allowed;
	}

	.page-info {
		font-size: 0.875rem;
		color: #71767b;
	}
</style>