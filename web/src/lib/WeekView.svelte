<script>
	import { onMount } from 'svelte';

	export let sessions = [];
	export let weekStart = null;

	const days = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
	const hours = Array.from({ length: 24 }, (_, i) => i);

	$: weekDays = getWeekDays(weekStart);

	function getWeekDays(start) {
		const days = [];
		const startDate = start ? new Date(start) : new Date();
		const dayOfWeek = startDate.getDay();
		const monday = new Date(startDate);
		monday.setDate(startDate.getDate() - dayOfWeek);

		for (let i = 0; i < 7; i++) {
			const date = new Date(monday);
			date.setDate(monday.getDate() + i);
			days.push(date);
		}
		return days;
	}

	function getSessionsForDay(date) {
		const dateStr = date.toISOString().split('T')[0];
		return sessions.filter(s => s.date === dateStr);
	}

	function formatHour(hour) {
		if (hour === 0) return '12a';
		if (hour < 12) return `${hour}a`;
		if (hour === 12) return '12p';
		return `${hour - 12}p`;
	}

	function getSessionPosition(session) {
		const start = new Date(session.start_time);
		const startHour = start.getHours() + start.getMinutes() / 60;
		const top = startHour * 24;

		let height = 24;
		if (session.end_time) {
			const end = new Date(session.end_time);
			const endHour = end.getHours() + end.getMinutes() / 60;
			height = Math.max((endHour - startHour) * 24, 12);
		}

		return { top: `${top}px`, height: `${height}px` };
	}

	function formatDuration(seconds) {
		if (!seconds) return '';
		const hours = Math.floor(seconds / 3600);
		const mins = Math.floor((seconds % 3600) / 60);
		if (hours > 0) return `${hours}h ${mins}m`;
		return `${mins}m`;
	}
</script>

<div class="week-view">
	<div class="header">
		<div class="time-column"></div>
		{#each weekDays as day}
			<div class="day-header" class:today={day.toDateString() === new Date().toDateString()}>
				<span class="day-name">{days[day.getDay()]}</span>
				<span class="day-date">{day.getDate()}</span>
			</div>
		{/each}
	</div>

	<div class="body">
		<div class="time-column">
			{#each hours as hour}
				<div class="hour-label">{formatHour(hour)}</div>
			{/each}
		</div>

		{#each weekDays as day}
			<div class="day-column" class:today={day.toDateString() === new Date().toDateString()}>
				{#each getSessionsForDay(day) as session}
					<div 
						class="session" 
						style="top: {getSessionPosition(session).top}; height: {getSessionPosition(session).height}"
						title="{formatDuration(session.duration_seconds)}"
					>
						<span class="session-duration">{formatDuration(session.duration_seconds)}</span>
					</div>
				{/each}

				<div class="hour-grid">
					{#each hours as hour}
						<div class="hour-line"></div>
					{/each}
				</div>
			</div>
		{/each}
	</div>
</div>

<style>
	.week-view {
		background: #16181c;
		border-radius: 12px;
		overflow: hidden;
	}

	.header {
		display: flex;
		border-bottom: 1px solid #2f3336;
	}

	.time-column {
		width: 40px;
		flex-shrink: 0;
	}

	.day-header {
		flex: 1;
		padding: 12px 4px;
		text-align: center;
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.day-header.today {
		background: rgba(29, 155, 240, 0.1);
	}

	.day-name {
		font-size: 0.75rem;
		color: #71767b;
		text-transform: uppercase;
	}

	.day-date {
		font-size: 1.25rem;
		font-weight: 600;
		color: #e7e9ea;
	}

	.day-header.today .day-date {
		color: #1d9bf0;
	}

	.body {
		display: flex;
		max-height: 400px;
		overflow-y: auto;
	}

	.time-column {
		width: 40px;
		flex-shrink: 0;
		display: flex;
		flex-direction: column;
	}

	.hour-label {
		height: 24px;
		font-size: 0.65rem;
		color: #71767b;
		text-align: right;
		padding-right: 4px;
		line-height: 24px;
	}

	.day-column {
		flex: 1;
		position: relative;
		min-height: 576px;
		background: #191c20;
		border-right: 1px solid #2f3336;
	}

	.day-column:last-child {
		border-right: none;
	}

	.day-column.today {
		background: rgba(29, 155, 240, 0.05);
	}

	.hour-grid {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		pointer-events: none;
	}

	.hour-line {
		height: 24px;
		border-bottom: 1px solid #2f3336;
	}

	.session {
		position: absolute;
		left: 2px;
		right: 2px;
		background: linear-gradient(135deg, rgba(29, 155, 240, 0.8), rgba(120, 86, 255, 0.8));
		border-radius: 4px;
		padding: 2px 4px;
		overflow: hidden;
		cursor: pointer;
		transition: transform 0.1s;
	}

	.session:hover {
		transform: scaleX(1.02);
		z-index: 10;
	}

	.session-duration {
		font-size: 0.7rem;
		color: white;
		font-weight: 600;
	}
</style>