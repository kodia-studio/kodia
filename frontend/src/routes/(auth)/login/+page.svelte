<script lang="ts">
	import { authApi } from '$lib/api/auth';
	import { toastStore } from '$lib/stores/toast.store';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth.store';
	import { onMount } from 'svelte';

	let email = $state('');
	let password = $state('');
	let isLoading = $state(false);

	onMount(() => {
		if ($authStore.isAuthenticated) {
			goto('/dashboard');
		}
	});

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		isLoading = true;

		try {
			const res = await authApi.login({ email, password });
			if (res.success) {
				toastStore.add('Welcome back!', 'success');
				goto('/dashboard');
			}
		} catch (err: any) {
			toastStore.add(err.message || 'Login failed', 'error');
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="min-h-screen flex items-center justify-center p-4 bg-slate-50 dark:bg-slate-950">
	<div class="w-full max-w-md">
		<div class="text-center mb-8">
			<div class="w-16 h-16 rounded-2xl bg-primary-600 mx-auto flex items-center justify-center text-white text-3xl font-bold mb-4 shadow-lg shadow-primary-500/20">K</div>
			<h1 class="text-3xl font-bold tracking-tight mb-2">Welcome to Kodia</h1>
			<p class="text-slate-500 dark:text-slate-400">Sign in to your account to continue</p>
		</div>

		<div class="card glass">
			<form onsubmit={handleSubmit} class="space-y-4">
				<div>
					<label for="email" class="block text-sm font-medium mb-1.5 ml-1">Email Address</label>
					<input
						type="email"
						id="email"
						placeholder="name@example.com"
						bind:value={email}
						class="input-standard"
						required
					/>
				</div>

				<div>
					<div class="flex items-center justify-between mb-1.5 ml-1">
						<label for="password" class="block text-sm font-medium">Password</label>
						<a href="/forgot-password" class="text-xs text-primary-600 hover:text-primary-700 font-medium">Forgot password?</a>
					</div>
					<input
						type="password"
						id="password"
						placeholder="••••••••"
						bind:value={password}
						class="input-standard"
						required
					/>
				</div>

				<button
					type="submit"
					class="w-full btn-primary mt-2 flex items-center justify-center gap-2"
					disabled={isLoading}
				>
					{#if isLoading}
						<svg class="animate-spin h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						Signing in...
					{:else}
						Sign in
					{/if}
				</button>
			</form>

			<div class="mt-6 text-center">
				<p class="text-sm text-slate-500">
					Don't have an account?
					<a href="/register" class="text-primary-600 hover:text-primary-700 font-medium ml-1">Create an account</a>
				</p>
			</div>
		</div>

		<div class="mt-8 text-center flex items-center justify-center gap-4">
			<button class="text-xs text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 transition-colors">Privacy Policy</button>
			<span class="w-1 h-1 rounded-full bg-slate-300"></span>
			<button class="text-xs text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 transition-colors">Terms of Service</button>
		</div>
	</div>
</div>
