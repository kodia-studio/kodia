<script lang="ts">
	import { toastStore } from '$lib/stores/toast.store';
	
	let email = $state('');
	let isLoading = $state(false);
	let isSubmitted = $state(false);

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		isLoading = true;

		try {
			// Mocking API call for forgot password since not implemented in backend yet
			await new Promise(resolve => setTimeout(resolve, 1000));
			isSubmitted = true;
			toastStore.success('Password reset link sent');
		} catch (err: any) {
			toastStore.error(err.message || 'Failed to send reset link');
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="min-h-screen flex items-center justify-center p-4 bg-slate-50 dark:bg-slate-950">
	<div class="w-full max-w-md">
		<div class="text-center mb-8">
			<div class="w-16 h-16 rounded-2xl bg-primary-600 mx-auto flex items-center justify-center text-white text-3xl font-bold mb-4 shadow-lg shadow-primary-500/20">K</div>
			<h1 class="text-3xl font-bold tracking-tight mb-2">Reset password</h1>
			<p class="text-slate-500 dark:text-slate-400">Enter your email to receive a reset link</p>
		</div>

		<div class="card glass">
			{#if isSubmitted}
				<div class="text-center py-6">
					<div class="w-16 h-16 rounded-full bg-green-100 dark:bg-green-900/30 text-green-600 mx-auto flex items-center justify-center mb-4">
						<svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
						</svg>
					</div>
					<h3 class="text-xl font-bold mb-2">Check your email</h3>
					<p class="text-sm text-slate-500 dark:text-slate-400 mb-6">
						We've sent a password reset link to <span class="font-medium text-slate-900 dark:text-white">{email}</span>
					</p>
					<button onclick={() => window.location.href = '/login'} class="btn-primary w-full inline-block">
						Back to login
					</button>
				</div>
			{:else}
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
							Sending link...
						{:else}
							Send reset link
						{/if}
					</button>
				</form>
			{/if}

			{#if !isSubmitted}
				<div class="mt-6 text-center">
					<p class="text-sm text-slate-500">
						Remember your password?
						<a href="/login" class="text-primary-600 hover:text-primary-700 font-medium ml-1">Sign in</a>
					</p>
				</div>
			{/if}
		</div>
	</div>
</div>
