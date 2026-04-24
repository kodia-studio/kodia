<script lang="ts">
	import { Toaster as Sonner } from "svelte-sonner";
	import { themeStore } from "$lib/stores/theme.svelte";

	interface Props {
		position?: "top-left" | "top-center" | "top-right" | "bottom-left" | "bottom-center" | "bottom-right";
	}

	let { position = "bottom-right" }: Props = $props();

	const theme = $derived(themeStore.dark ? 'dark' : 'light');
</script>

<Sonner
	{theme}
	class="toaster group"
	{position}
	toastOptions={{
		classes: {
			toast: "group toast group-[.toaster]:bg-white group-[.toaster]:dark:bg-slate-900 group-[.toaster]:rounded-kodia group-[.toaster]:border group-[.toaster]:border-slate-200 dark:group-[.toaster]:border-slate-800 group-[.toaster]:shadow-kodia-lg",
			description: "group-[.toast]:text-slate-500",
			actionButton: "group-[.toast]:bg-primary group-[.toast]:text-white",
			cancelButton: "group-[.toast]:bg-slate-100 group-[.toast]:text-slate-500",
			title: "group-[.toast]:font-black group-[.toast]:tracking-tight group-[.toast]:text-slate-900 dark:group-[.toast]:text-white",
			success: "group-[.toast]:!text-emerald-500",
			error: "group-[.toast]:!text-red-500",
			info: "group-[.toast]:!text-primary"
		}
	}}
/>
