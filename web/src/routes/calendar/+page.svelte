<script>
	import { onMount } from 'svelte';
	import Calendar from '$lib/Calendar.svelte';

	let sessions = [];
	let currentWeekStart = getSundayOfCurrentWeek();
	let isLoading = true;

	function getSundayOfCurrentWeek() {
		const today = new Date();
		const dayOfWeek = today.getDay();
		const sunday = new Date(today);
		sunday.setDate(today.getDate() - dayOfWeek);
		return formatDateLocal(sunday);
	}

	function formatDateLocal(date) {
		const year = date.getFullYear();
		const month = String(date.getMonth() + 1).padStart(2, '0');
		const day = String(date.getDate()).padStart(2, '0');
		return `${year}-${month}-${day}`;
	}

	onMount(async () => {
		await loadSessions();
		isLoading = false;
	});

	async function loadSessions() {
		try {
			const res = await fetch(`/api/sessions?start=${currentWeekStart}`);
			sessions = await res.json();
		} catch (e) {
			console.error('Failed to load sessions', e);
		}
	}

	async function navigateWeek(newWeekStart) {
		currentWeekStart = newWeekStart;
		await loadSessions();
	}
</script>

<svelte:head>
	<title>Calendar - Bedrock Timeline</title>
</svelte:head>

<main>
	<header>
		<div class="header-left">
			<a href="/" class="back-link">← Back to Dashboard</a>
			<h1>All Player Sessions</h1>
			<p class="subtitle">Weekly view of all player activity</p>
		</div>
	</header>

	{#if isLoading}
		<div class="loading">Loading calendar...</div>
	{:else}
		<div class="calendar-wrapper">
			<Calendar {sessions} weekStart={currentWeekStart} onNavigate={navigateWeek} />
		</div>
	{/if}
</main>

<style>
	main {
		min-height: 100vh;
		background: #0f1419;
		color: #e7e9ea;
	}

	header {
		padding: 20px;
		border-bottom: 1px solid #2f3336;
	}

	.header-left {
		max-width: 1400px;
		margin: 0 auto;
	}

	.back-link {
		display: inline-block;
		color: #1d9bf0;
		text-decoration: none;
		font-size: 0.875rem;
		margin-bottom: 8px;
		transition: opacity 0.2s;
	}

	.back-link:hover {
		opacity: 0.8;
	}

	h1 {
		margin: 0;
		font-size: 1.5rem;
		color: #e7e9ea;
	}

	.subtitle {
		margin: 4px 0 0 0;
		font-size: 0.875rem;
		color: #71767b;
	}

	.loading {
		text-align: center;
		padding: 40px;
		color: #71767b;
	}

	.calendar-wrapper {
		max-width: 1400px;
		margin: 20px auto;
		padding: 0 20px;
	}
</style>