<script lang="ts">
	import { authStore } from '$lib/stores/auth.store';

	const stats = [
		{ name: 'Total Revenue', value: '$45,231.89', change: '+20.1%', icon: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z' },
		{ name: 'Active Users', value: '+2,350', change: '+180.1%', icon: 'M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z' },
		{ name: 'New Customers', value: '+12,234', change: '+19%', icon: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z' },
		{ name: 'Engagement Rate', value: '64.5%', change: '+12.5%', icon: 'M13 7h8m0 0v8m0-8l-8 8-4-4-6 6' }
	];
</script>

<div class="space-y-8">
	<div>
		<h1 class="text-3xl font-bold tracking-tight">Dashboard Overview</h1>
		<p class="text-slate-500 dark:text-slate-400 mt-1">Welcome back, {$authStore.user?.name}. Here's what's happening today.</p>
	</div>

	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
		{#each stats as stat}
			<div class="card glass flex flex-col gap-4">
				<div class="flex items-center justify-between">
					<div class="w-10 h-10 rounded-lg bg-primary-600/10 flex items-center justify-center text-primary-600">
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={stat.icon} />
						</svg>
					</div>
					<span class="text-xs font-semibold px-2 py-1 rounded-full bg-green-100 dark:bg-green-900/30 text-green-600">
						{stat.change}
					</span>
				</div>
				<div>
					<p class="text-sm text-slate-500 font-medium">{stat.name}</p>
					<h3 class="text-2xl font-bold mt-1 tracking-tight">{stat.value}</h3>
				</div>
			</div>
		{/each}
	</div>

	<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
		<div class="lg:col-span-2 card glass">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-lg font-bold">Recent Activity</h3>
				<button class="text-xs text-primary-600 font-medium hover:underline">View all</button>
			</div>
			<div class="space-y-6">
				{#each Array(5) as _, i}
					<div class="flex items-center gap-4">
						<div class="w-10 h-10 rounded-full bg-slate-100 dark:bg-slate-800 shrink-0"></div>
						<div class="flex-1 min-w-0">
							<p class="text-sm font-medium">New project "Kodia Framework" created</p>
							<p class="text-xs text-slate-500 mt-0.5">2 hours ago</p>
						</div>
						<span class="text-xs px-2 py-1 rounded bg-blue-100 dark:bg-blue-900/30 text-blue-600 font-medium">Update</span>
					</div>
				{/each}
			</div>
		</div>

		<div class="card glass">
			<h3 class="text-lg font-bold mb-6">Quick Actions</h3>
			<div class="grid grid-cols-2 gap-4">
				<button class="flex flex-col items-center gap-2 p-4 rounded-xl hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors border border-transparent hover:border-slate-200 dark:hover:border-slate-700">
					<div class="w-12 h-12 rounded-full bg-purple-100 dark:bg-purple-900/30 text-purple-600 flex items-center justify-center">
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
						</svg>
					</div>
					<span class="text-xs font-medium">New Feature</span>
				</button>
				<button class="flex flex-col items-center gap-2 p-4 rounded-xl hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors border border-transparent hover:border-slate-200 dark:hover:border-slate-700">
					<div class="w-12 h-12 rounded-full bg-blue-100 dark:bg-blue-900/30 text-blue-600 flex items-center justify-center">
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
						</svg>
					</div>
					<span class="text-xs font-medium">Reports</span>
				</button>
				<button class="flex flex-col items-center gap-2 p-4 rounded-xl hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors border border-transparent hover:border-slate-200 dark:hover:border-slate-700">
					<div class="w-12 h-12 rounded-full bg-orange-100 dark:bg-orange-900/30 text-orange-600 flex items-center justify-center">
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					</div>
					<span class="text-xs font-medium">Schedule</span>
				</button>
				<button class="flex flex-col items-center gap-2 p-4 rounded-xl hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors border border-transparent hover:border-slate-200 dark:hover:border-slate-700">
					<div class="w-12 h-12 rounded-full bg-green-100 dark:bg-green-900/30 text-green-600 flex items-center justify-center">
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065z" />
						</svg>
					</div>
					<span class="text-xs font-medium">System</span>
				</button>
			</div>
		</div>
	</div>
</div>
