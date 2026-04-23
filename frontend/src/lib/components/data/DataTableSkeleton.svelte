<script lang="ts">
	import Skeleton from '../ui/Skeleton.svelte';

	interface Props {
		rows?: number;
		columns?: number;
	}

	let { rows = 5, columns = 4 }: Props = $props();
</script>

<div class="space-y-6">
	<!-- Header / Search Area -->
	<div class="flex items-center justify-between gap-4">
		<Skeleton width="300px" height="40px" class="rounded-xl" />
		<div class="flex gap-2">
			<Skeleton width="100px" height="40px" class="rounded-xl" />
			<Skeleton width="40px" height="40px" class="rounded-xl" />
		</div>
	</div>

	<!-- Table Area -->
	<div class="glass border border-slate-200 dark:border-slate-800 rounded-3xl overflow-hidden">
		<!-- Table Header -->
		<div class="bg-slate-50 dark:bg-slate-900/50 p-4 flex gap-4">
			{#each Array(columns) as _}
				<div class="flex-1">
					<Skeleton variant="text" width="60%" />
				</div>
			{/each}
		</div>

		<!-- Table Rows -->
		<div class="divide-y divide-slate-100 dark:divide-slate-800">
			{#each Array(rows) as _}
				<div class="p-4 flex gap-4 items-center">
					{#each Array(columns) as _}
						<div class="flex-1">
							<Skeleton variant="text" width="80%" />
						</div>
					{/each}
				</div>
			{/each}
		</div>
	</div>

	<!-- Pagination -->
	<div class="flex items-center justify-between">
		<Skeleton width="150px" height="20px" />
		<div class="flex gap-2">
			<Skeleton width="32px" height="32px" />
			<Skeleton width="32px" height="32px" />
			<Skeleton width="32px" height="32px" />
		</div>
	</div>
</div>
