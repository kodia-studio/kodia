<script lang="ts">
	import { cn } from "$lib/utils";
	import { Upload, X, File as FileIcon, Image as ImageIcon } from "lucide-svelte";
	import { fade, slide } from "svelte/transition";

	interface FileStatus {
		file: File;
		progress: number;
		status: "uploading" | "completed" | "error";
		preview?: string;
	}

	interface Props {
		label?: string;
		accept?: string;
		multiple?: boolean;
		maxSize?: number; // in MB
		class?: string;
	}

	let { label, accept = "*", multiple = false, maxSize = 10, class: className }: Props = $props();

	let files = $state<FileStatus[]>([]);
	let dragging = $state(false);

	function handleFiles(newFiles: FileList | null) {
		if (!newFiles) return;

		const addedFiles: FileStatus[] = Array.from(newFiles).map(file => ({
			file,
			progress: 0,
			status: "uploading",
			preview: file.type.startsWith("image/") ? URL.createObjectURL(file) : undefined
		}));

		if (multiple) {
			files = [...files, ...addedFiles];
		} else {
			files = addedFiles;
		}

		// Simulate upload progress
		addedFiles.forEach(fs => {
			let p = 0;
			const interval = setInterval(() => {
				p += Math.random() * 30;
				if (p >= 100) {
					p = 100;
					fs.status = "completed";
					clearInterval(interval);
				}
				fs.progress = p;
			}, 300);
		});
	}

	function removeFile(index: number) {
		const fs = files[index];
		if (fs.preview) URL.revokeObjectURL(fs.preview);
		files = files.filter((_, i) => i !== index);
	}
</script>

<div class={cn("space-y-4", className)}>
	{#if label}
		<label 
			for="file-input"
			class="text-xs font-black uppercase tracking-widest text-slate-500 ml-1"
		>
			{label}
		</label>
	{/if}

	<div
		class={cn(
			"relative group cursor-pointer rounded-kodia-lg border-2 border-dashed transition-all duration-300 p-12 text-center outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2",
			dragging 
				? "border-primary bg-primary/5 scale-[0.99]" 
				: "border-slate-200 dark:border-slate-800 hover:border-primary/50 hover:bg-slate-50/50 dark:hover:bg-slate-900/50"
		)}
		role="button"
		tabindex="0"
		ondragover={(e) => { e.preventDefault(); dragging = true; }}
		ondragleave={() => dragging = false}
		ondrop={(e) => { e.preventDefault(); dragging = false; handleFiles(e.dataTransfer?.files || null); }}
		onclick={() => document.getElementById("file-input")?.click()}
		onkeydown={(e) => { if (e.key === "Enter" || e.key === " ") { e.preventDefault(); document.getElementById("file-input")?.click(); } }}
	>
		<input
			id="file-input"
			type="file"
			class="hidden"
			{accept}
			{multiple}
			onchange={(e) => handleFiles((e.target as HTMLInputElement).files)}
		/>

		<div class="flex flex-col items-center gap-3">
			<div class="p-4 rounded-full bg-primary/10 text-primary group-hover:scale-110 transition-transform duration-300">
				<Upload class="w-8 h-8" />
			</div>
			<div class="space-y-1">
				<p class="text-lg font-black tracking-tight text-slate-900 dark:text-white">
					Drop files here or click to upload
				</p>
				<p class="text-sm text-slate-500 font-medium">
					Maximum file size: {maxSize}MB
				</p>
			</div>
		</div>
	</div>

	{#if files.length > 0}
		<div class="grid grid-cols-1 sm:grid-cols-2 gap-4" transition:slide>
			{#each files as fs, i}
				<div class="glass p-3 rounded-kodia flex items-center gap-4 relative group overflow-hidden shadow-sm">
					<!-- Progress Bar Background -->
					<div 
						class="absolute bottom-0 left-0 h-1 bg-primary/20 transition-all duration-300"
						style="width: {fs.progress}%"
					></div>

					<div class="w-12 h-12 rounded-lg bg-slate-100 dark:bg-slate-800 flex items-center justify-center overflow-hidden shrink-0">
						{#if fs.preview}
							<img src={fs.preview} alt="preview" class="w-full h-full object-cover" />
						{:else}
							<FileIcon class="w-6 h-6 text-slate-400" />
						{/if}
					</div>

					<div class="flex-1 min-w-0">
						<p class="text-xs font-bold text-slate-900 dark:text-white truncate">
							{fs.file.name}
						</p>
						<p class="text-[10px] text-slate-500 font-bold uppercase tracking-tight">
							{(fs.file.size / 1024 / 1024).toFixed(2)} MB • {fs.status}
						</p>
					</div>

					<button
						class="p-1.5 rounded-full hover:bg-red-500/10 text-slate-400 hover:text-red-500 transition-all"
						onclick={(e) => { e.stopPropagation(); removeFile(i); }}
					>
						<X class="w-4 h-4" />
					</button>
				</div>
			{/each}
		</div>
	{/if}
</div>
