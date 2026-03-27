<script>
	export let stats;
</script>

<div class="stats">
	<div class="stat">
		<span class="stat-value">{stats.total_players || 0}</span>
		<span class="stat-label">Total Players</span>
	</div>
	<div class="stat">
		<span class="stat-value">{stats.total_joins || 0}</span>
		<span class="stat-label">Joins</span>
	</div>
	<div class="stat">
		<span class="stat-value">{stats.total_leaves || 0}</span>
		<span class="stat-label">Leaves</span>
	</div>
</div>

<style>
	.stats {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 12px;
		background: #16181c;
		border-radius: 12px;
		padding: 16px;
	}

	.stat {
		text-align: center;
	}

	.stat-value {
		display: block;
		font-size: 1.25rem;
		font-weight: 700;
		color: #e7e9ea;
	}

	.stat-label {
		font-size: 0.75rem;
		color: #71767b;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	@media (max-width: 768px) {
		.stats {
			grid-template-columns: repeat(3, 1fr);
		}
	}
</style>