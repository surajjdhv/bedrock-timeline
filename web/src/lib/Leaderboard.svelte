<script>
	import { onMount } from 'svelte';

	export let players = [];
</script>

<div class="leaderboard">
	<h3>Leaderboard <span class="period">(30 Days)</span></h3>
	{#if players.length === 0}
		<p class="empty">No players yet</p>
	{:else}
		<ul class="list">
			{#each players.slice(0, 10) as player, i (player.name)}
				<li class="player-item" class:top-three={i < 3}>
					<span class="rank">{i + 1}</span>
					<div class="player-avatar" style="background: {getPlayerColor(i)}">{player.name[0].toUpperCase()}</div>
					<div class="player-info">
						<span class="player-name">{player.name}</span>
						{#if player.total_playtime}
							<span class="player-playtime">{formatPlaytime(player.total_playtime)}</span>
						{/if}
					</div>
					{#if i < 3}
						<span class="medal">{getMedal(i)}</span>
					{/if}
				</li>
			{/each}
		</ul>
	{/if}
</div>

<script context="module">
	function formatPlaytime(seconds) {
		if (!seconds) return '0m';
		const hours = Math.floor(seconds / 3600);
		const mins = Math.floor((seconds % 3600) / 60);
		if (hours > 0) return `${hours}h ${mins}m`;
		return `${mins}m`;
	}

	function getPlayerColor(index) {
		const colors = ['#ffd700', '#c0c0c0', '#cd7f32', '#1d9bf0', '#7856ff', '#00ba7c', '#f91880', '#ff7a00'];
		return colors[index % colors.length];
	}

	function getMedal(index) {
		const medals = ['🥇', '🥈', '🥉'];
		return medals[index];
	}
</script>

<style>
	.leaderboard {
		background: #16181c;
		border-radius: 12px;
		padding: 16px;
	}

	h3 {
		margin: 0 0 12px 0;
		font-size: 1rem;
		font-weight: 600;
		color: #e7e9ea;
	}

	.period {
		font-size: 0.75rem;
		font-weight: 400;
		color: #71767b;
	}

	.empty {
		color: #71767b;
		font-size: 0.875rem;
		text-align: center;
		padding: 20px 0;
	}

	.list {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.player-item {
		display: flex;
		align-items: center;
		gap: 10px;
		padding: 10px 12px;
		background: #191c20;
		border-radius: 8px;
		transition: background-color 0.2s;
	}

	.player-item:hover {
		background: #1f2227;
	}

	.player-item.top-three {
		background: linear-gradient(135deg, rgba(255, 215, 0, 0.1) 0%, rgba(255, 215, 0, 0.05) 100%);
	}

	.player-item.top-three:nth-child(2) {
		background: linear-gradient(135deg, rgba(192, 192, 192, 0.1) 0%, rgba(192, 192, 192, 0.05) 100%);
	}

	.player-item.top-three:nth-child(3) {
		background: linear-gradient(135deg, rgba(205, 127, 50, 0.1) 0%, rgba(205, 127, 50, 0.05) 100%);
	}

	.rank {
		min-width: 24px;
		font-size: 0.875rem;
		font-weight: 600;
		color: #71767b;
	}

	.player-item.top-three .rank {
		color: #e7e9ea;
	}

	.player-avatar {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: bold;
		font-size: 0.875rem;
		color: white;
		flex-shrink: 0;
	}

	.player-info {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 2px;
		min-width: 0;
	}

	.player-name {
		font-weight: 500;
		color: #e7e9ea;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.player-playtime {
		font-size: 0.75rem;
		color: #71767b;
	}

	.medal {
		font-size: 1.25rem;
	}
</style>