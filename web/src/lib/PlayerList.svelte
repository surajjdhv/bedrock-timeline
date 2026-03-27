<script>
	export let players;
	export let onSelect;

	function formatLastSeen(timestamp) {
		if (!timestamp) return 'Never';
		const date = new Date(timestamp);
		const now = new Date();
		const diff = now - date;
		const minutes = Math.floor(diff / 60000);
		const hours = Math.floor(diff / 3600000);
		const days = Math.floor(diff / 86400000);

		if (days > 0) return `${days}d ago`;
		if (hours > 0) return `${hours}h ago`;
		if (minutes > 0) return `${minutes}m ago`;
		return 'Just now';
	}
</script>

<div class="player-list">
	<h2>Players ({players.length})</h2>
	
	{#if players.length === 0}
		<p class="empty">No players recorded yet.</p>
	{:else}
		<ul class="list">
			{#each players as player (player.name)}
				<li class="player" on:click={() => onSelect?.(player.name)} role="button" tabindex="0">
					<div class="player-avatar">{player.name[0].toUpperCase()}</div>
					<div class="player-info">
						<span class="player-name">{player.name}</span>
						<span class="last-seen">Last seen: {formatLastSeen(player.last_join)}</span>
					</div>
					<div class="player-action">→</div>
				</li>
			{/each}
		</ul>
	{/if}
</div>

<style>
	.player-list {
		background: #16181c;
		border-radius: 12px;
		padding: 16px;
	}

	h2 {
		margin: 0 0 12px 0;
		font-size: 1rem;
		font-weight: 600;
	}

	.empty {
		color: #71767b;
		text-align: center;
		padding: 20px;
		font-size: 0.875rem;
	}

	.list {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 4px;
		max-height: 300px;
		overflow-y: auto;
	}

	.player {
		display: flex;
		align-items: center;
		gap: 10px;
		padding: 8px;
		border-radius: 6px;
		transition: background-color 0.2s;
		cursor: pointer;
	}

	.player:hover {
		background: #191c20;
	}

	.player:focus {
		outline: 2px solid #1d9bf0;
		outline-offset: -2px;
	}

	.player-avatar {
		width: 28px;
		height: 28px;
		border-radius: 50%;
		background: linear-gradient(135deg, #1d9bf0, #7856ff);
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 0.75rem;
		font-weight: bold;
		color: white;
		flex-shrink: 0;
	}

	.player-info {
		display: flex;
		flex-direction: column;
		min-width: 0;
		flex: 1;
	}

	.player-name {
		font-weight: 600;
		font-size: 0.875rem;
		color: #e7e9ea;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.last-seen {
		font-size: 0.75rem;
		color: #71767b;
	}

	.player-action {
		color: #71767b;
		font-size: 0.875rem;
		transition: color 0.2s;
	}

	.player:hover .player-action {
		color: #1d9bf0;
	}
</style>