<script>
	import { onMount } from 'svelte';

	export let sessions = [];
	export let playerName = '';
	export let weekStart = null;
	export let onNavigate = null;

	$: weekDays = getWeekDays(weekStart);
	$: playerColors = getPlayerColors(sessions);
	$: sessionsByDay = groupSessionsByDay(sessions);

	const hourHeight = 40;
	const hours = Array.from({ length: 24 }, (_, i) => i);

	function parseLocalDate(dateStr) {
		if (!dateStr) return new Date();
		const [year, month, day] = dateStr.split('-').map(Number);
		return new Date(year, month - 1, day);
	}

	function getWeekDays(start) {
		const days = [];
		const startDate = parseLocalDate(start);
		const dayOfWeek = startDate.getDay();
		const sunday = new Date(startDate);
		sunday.setDate(startDate.getDate() - dayOfWeek);

		for (let i = 0; i < 7; i++) {
			const date = new Date(sunday);
			date.setDate(sunday.getDate() + i);
			days.push(date);
		}
		return days;
	}

	function groupSessionsByDay(sessionsData) {
		const grouped = {};
		if (!sessionsData || sessionsData.length === 0) {
			return grouped;
		}
		
		sessionsData.forEach(session => {
			if (!session.date) return;
			if (!grouped[session.date]) {
				grouped[session.date] = [];
			}
			grouped[session.date].push(session);
		});
		
		// Calculate positions for overlapping sessions
		Object.keys(grouped).forEach(date => {
			grouped[date] = calculateSessionPositions(grouped[date]);
		});
		
		return grouped;
	}

	function calculateSessionPositions(daySessions) {
		if (!daySessions || daySessions.length === 0) return [];
		
		// Sort by start time
		const sorted = [...daySessions].sort((a, b) => {
			const startA = new Date(a.start_time).getTime();
			const startB = new Date(b.start_time).getTime();
			const endA = a.end_time ? new Date(a.end_time).getTime() : startA + 3600000;
			const endB = b.end_time ? new Date(b.end_time).getTime() : startB + 3600000;
			if (startA !== startB) return startA - startB;
			return (endB - startB) - (endA - startA);
		});
		
		// Group overlapping sessions into columns
		const columns = [];
		
		sorted.forEach(session => {
			const start = new Date(session.start_time).getTime();
			const end = session.end_time ? new Date(session.end_time).getTime() : start + 3600000;
			
			// Find the first column where this session doesn't overlap
			let placed = false;
			for (let i = 0; i < columns.length; i++) {
				const lastInColumn = columns[i][columns[i].length - 1];
				const lastEnd = lastInColumn.end_time ? new Date(lastInColumn.end_time).getTime() : new Date(lastInColumn.start_time).getTime() + 3600000;
				
				if (start >= lastEnd) {
					// No overlap, add to this column
					columns[i].push(session);
					placed = true;
					break;
				}
			}
			
			if (!placed) {
				// Need a new column
				columns.push([session]);
			}
		});
		
		const totalColumns = columns.length;
		
		// Assign positions based on columns
		const result = [];
		const padding = 2;
		
		for (let colIndex = 0; colIndex < columns.length; colIndex++) {
			const column = columns[colIndex];
			const widthPercent = 100 / totalColumns;
			
			column.forEach(session => {
				result.push({
					...session,
					columnIndex: colIndex,
					totalColumns: totalColumns,
					left: `calc(${widthPercent * colIndex}% + ${padding}px)`,
					width: `calc(${widthPercent}% - ${padding * 2}px)`
				});
			});
		}
		
		return result;
	}

	function getPlayerColors(sessionsData) {
		const colors = [
			'#1d9bf0', '#7856ff', '#00ba7c', '#f91880', '#ff7a00', '#794bc4', '#1dcaff', '#17bf63', '#ffad1f', '#e0245e', '#189e4c', '#557cee'
		];
		const playerColorMap = {};
		let colorIndex = 0;
		
		sessionsData.forEach(session => {
			if (session.player_name && !playerColorMap[session.player_name]) {
				playerColorMap[session.player_name] = colors[colorIndex % colors.length];
				colorIndex++;
			}
		});
		
		return playerColorMap;
	}

	function getSessionPosition(session) {
		const start = new Date(session.start_time);
		const startHour = start.getHours() + start.getMinutes() / 60;
		const top = startHour * hourHeight;

		let height = hourHeight;
		if (session.end_time) {
			const end = new Date(session.end_time);
			const endHour = end.getHours() + end.getMinutes() / 60;
			const duration = endHour - startHour;
			height = Math.max(duration * hourHeight, 20);
		}

		return { top: `${top}px`, height: `${height}px` };
	}

	function formatHour(hour) {
		if (hour === 0) return '12 AM';
		if (hour < 12) return `${hour} AM`;
		if (hour === 12) return '12 PM';
		return `${hour - 12} PM`;
	}

	function formatDateShort(date) {
		return date.toLocaleDateString('en-US', { weekday: 'short' });
	}

	function formatDateNum(date) {
		return date.getDate();
	}

	function formatDuration(seconds) {
		if (!seconds) return '';
		const hours = Math.floor(seconds / 3600);
		const mins = Math.floor((seconds % 3600) / 60);
		if (hours > 0) return `${hours}h ${mins}m`;
		return `${mins}m`;
	}

	function goToPrevWeek() {
		if (onNavigate) {
			const currentStart = parseLocalDate(weekStart);
			const dayOfWeek = currentStart.getDay();
			const sunday = new Date(currentStart);
			sunday.setDate(currentStart.getDate() - dayOfWeek);
			const prevSunday = new Date(sunday);
			prevSunday.setDate(sunday.getDate() - 7);
			onNavigate(formatDateLocal(prevSunday));
		}
	}

	function goToNextWeek() {
		if (onNavigate) {
			const currentStart = parseLocalDate(weekStart);
			const dayOfWeek = currentStart.getDay();
			const sunday = new Date(currentStart);
			sunday.setDate(currentStart.getDate() - dayOfWeek);
			const nextSunday = new Date(sunday);
			nextSunday.setDate(sunday.getDate() + 7);
			onNavigate(formatDateLocal(nextSunday));
		}
	}

	function goToToday() {
		if (onNavigate) {
			const today = new Date();
			const dayOfWeek = today.getDay();
			const sunday = new Date(today);
			sunday.setDate(today.getDate() - dayOfWeek);
			onNavigate(formatDateLocal(sunday));
		}
	}

	function formatDateLocal(date) {
		const year = date.getFullYear();
		const month = String(date.getMonth() + 1).padStart(2, '0');
		const day = String(date.getDate()).padStart(2, '0');
		return `${year}-${month}-${day}`;
	}

	function isToday(date) {
		return date.toDateString() === new Date().toDateString();
	}
</script>

<div class="calendar-container">
	<div class="calendar-header">
		<div class="title-section">
			<h3>{playerName ? `${playerName}'s Sessions` : 'All Player Sessions'}</h3>
			<p class="week-range">
				{weekDays[0].toLocaleDateString('en-US', { month: 'short', day: 'numeric' })} - {weekDays[6].toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}
			</p>
		</div>
		{#if onNavigate}
			<div class="nav-buttons">
				<button class="nav-btn" on:click={goToToday}>Today</button>
				<button class="nav-btn" on:click={goToPrevWeek}>←</button>
				<button class="nav-btn" on:click={goToNextWeek}>→</button>
			</div>
		{/if}
	</div>

	<div class="calendar-body">
		<div class="time-column">
			<div class="time-spacer"></div>
			{#each hours as hour}
				<div class="hour-label">{formatHour(hour)}</div>
			{/each}
		</div>

		<div class="days-container">
			{#each weekDays as day}
				{@const dateStr = formatDateLocal(day)}
				{@const daySessions = sessionsByDay[dateStr] || []}
				<div class="day-column" class:today={isToday(day)}>
					<div class="day-header">
						<span class="day-name">{formatDateShort(day)}</span>
						<span class="day-num" class:today-num={isToday(day)}>{formatDateNum(day)}</span>
					</div>
					<div class="day-sessions">
						{#each daySessions as session}
							<div 
								class="session" 
								style="top: {getSessionPosition(session).top}; height: {getSessionPosition(session).height}; left: {session.left}; width: {session.width}; background: {playerColors[session.player_name] || '#1d9bf0'}"
								title="{session.player_name}: {formatDuration(session.duration_seconds)}"
							>
								<span class="session-name">{session.player_name}</span>
								{#if session.duration_seconds}
									<span class="session-duration">{formatDuration(session.duration_seconds)}</span>
								{/if}
							</div>
						{/each}
						{#each hours as hour}
							<div class="hour-line"></div>
						{/each}
					</div>
				</div>
			{/each}
		</div>
	</div>

	{#if Object.keys(playerColors).length > 1}
		<div class="legend">
			{#each Object.entries(playerColors) as [player, color]}
				<div class="legend-item">
					<div class="legend-dot" style="background: {color}"></div>
					<span>{player}</span>
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	.calendar-container {
		background: #16181c;
		border-radius: 12px;
		overflow: hidden;
	}

	.calendar-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		padding: 16px 20px;
		border-bottom: 1px solid #2f3336;
	}

	.title-section {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	h3 {
		margin: 0;
		font-size: 1.1rem;
		font-weight: 600;
		color: #e7e9ea;
	}

	.week-range {
		margin: 0;
		font-size: 0.8rem;
		color: #71767b;
	}

	.nav-buttons {
		display: flex;
		gap: 8px;
	}

	.nav-btn {
		background: transparent;
		border: 1px solid #2f3336;
		color: #e7e9ea;
		padding: 6px 12px;
		border-radius: 6px;
		font-size: 0.85rem;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.nav-btn:hover {
		background: #2f3336;
	}

	.calendar-body {
		display: flex;
		max-height: 600px;
		overflow-y: auto;
	}

	.time-column {
		flex-shrink: 0;
		width: 60px;
		display: flex;
		flex-direction: column;
		border-right: 1px solid #2f3336;
	}

	.time-spacer {
		height: 50px;
		border-bottom: 1px solid #2f3336;
	}

	.hour-label {
		height: 40px;
		font-size: 0.7rem;
		color: #71767b;
		text-align: right;
		padding-right: 8px;
		line-height: 40px;
		border-bottom: 1px solid #191c20;
	}

	.days-container {
		display: flex;
		flex: 1;
	}

	.day-column {
		flex: 1;
		min-width: 100px;
		border-right: 1px solid #2f3336;
	}

	.day-column:last-child {
		border-right: none;
	}

	.day-column.today {
		background: rgba(29, 155, 240, 0.05);
	}

	.day-header {
		height: 50px;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		border-bottom: 1px solid #2f3336;
		background: #191c20;
	}

	.day-name {
		font-size: 0.75rem;
		color: #71767b;
		text-transform: uppercase;
	}

	.day-num {
		font-size: 1.25rem;
		font-weight: 600;
		color: #e7e9ea;
		margin-top: 2px;
	}

	.day-num.today-num {
		color: #1d9bf0;
	}

	.day-sessions {
		position: relative;
		min-height: 960px;
	}

	.hour-line {
		height: 40px;
		border-bottom: 1px solid #191c20;
	}

	.session {
		position: absolute;
		border-radius: 4px;
		padding: 4px 6px;
		overflow: hidden;
		cursor: pointer;
		transition: transform 0.1s, z-index 0s;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
	}

	.session:hover {
		transform: translateY(-1px);
		box-shadow: 0 2px 6px rgba(0, 0, 0, 0.4);
		z-index: 10 !important;
	}

	.session-name {
		display: block;
		font-size: 0.7rem;
		font-weight: 600;
		color: white;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.session-duration {
		display: block;
		font-size: 0.6rem;
		color: rgba(255, 255, 255, 0.8);
		margin-top: 1px;
	}

	.legend {
		display: flex;
		flex-wrap: wrap;
		gap: 12px;
		padding: 12px 20px;
		border-top: 1px solid #2f3336;
	}

	.legend-item {
		display: flex;
		align-items: center;
		gap: 6px;
	}

	.legend-dot {
		width: 12px;
		height: 12px;
		border-radius: 3px;
	}

	.legend-item span {
		font-size: 0.75rem;
		color: #e7e9ea;
	}
</style>