<script>
	import { onMount } from 'svelte';

	export let players = [];
</script>

<div class="online-section">
	<h3>Online Now</h3>
	{#if players.length === 0}
		<p class="empty">No players online</p>
	{:else}
		<ul class="online-list">
			{#each players as player (player.name)}
				<li class="online-player">
					<div class="status-dot"></div>
					<span class="player-name">{player.name}</span>
				</li>
			{/each}
		</ul>
	{/if}
</div>

<style>
	.online-section {
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

	.empty {
		color: #71767b;
		font-size: 0.875rem;
		text-align: center;
		padding: 16px 0;
	}

	.online-list {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-wrap: wrap;
		gap: 8px;
	}

	.online-player {
		display: flex;
		align-items: center;
		gap: 8px;
		background: linear-gradient(135deg, rgba(0, 186, 124, 0.2) 0%, rgba(0, 186, 124, 0.1) 100%);
		border: 1px solid rgba(0, 186, 124, 0.3);
		padding: 8px 12px;
		border-radius: 20px;
	}

	.status-dot {
		width: 8px;
		height: 8px;
		border-radius: 50%;
		background: #00ba7c;
		animation: pulse 2s infinite;
	}

	@keyframes pulse {
		0%, 100% { opacity: 1; }
		50% { opacity: 0.5; }
	}

	.player-name {
		font-size: 0.875rem;
		font-weight: 500;
		color: #e7e9ea;
	}
</style>