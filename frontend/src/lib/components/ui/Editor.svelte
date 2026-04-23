<script lang="ts">
	import { onMount, onDestroy } from "svelte";
	import { Editor } from "@tiptap/core";
	import StarterKit from "@tiptap/starter-kit";
	import { cn } from "$lib/utils";
	import { Bold, Italic, List, ListOrdered, Quote, Heading1, Heading2 } from "lucide-svelte";

	interface Props {
		content?: string;
		label?: string;
		class?: string;
		onchange?: (content: string) => void;
	}

	let { content = "", label, class: className, onchange }: Props = $props();

	let element = $state<HTMLElement>();
	let editor = $state<Editor>();

	onMount(() => {
		editor = new Editor({
			element: element!,
			extensions: [StarterKit],
			content,
			onUpdate: ({ editor }) => {
				onchange?.(editor.getHTML());
			},
			editorProps: {
				attributes: {
					class: "prose prose-slate dark:prose-invert max-w-none focus:outline-none min-h-[150px] p-4 font-medium"
				}
			}
		});
	});

	onDestroy(() => {
		if (editor) {
			editor.destroy();
		}
	});

	const isActive = (type: string, options?: any) => editor?.isActive(type, options);
</script>

<div class={cn("space-y-2", className)}>
	{#if label}
		<label 
			for="tiptap-editor"
			class="text-xs font-black uppercase tracking-widest text-slate-500 ml-1"
		>
			{label}
		</label>
	{/if}

	<div 
		id="tiptap-editor"
		class="glass rounded-kodia-lg border border-slate-200 dark:border-slate-800 overflow-hidden focus-within:ring-2 focus-within:ring-primary focus-within:ring-offset-2 transition-all"
	>
		<!-- Toolbar -->
		<div class="flex items-center gap-1 p-1 border-b border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/80">
			<button
				class={cn("p-2 rounded-md hover:bg-primary/10 transition-colors", isActive('bold') && "bg-primary/20 text-primary")}
				onclick={() => editor?.chain().focus().toggleBold().run()}
			>
				<Bold class="w-4 h-4" />
			</button>
			<button
				class={cn("p-2 rounded-md hover:bg-primary/10 transition-colors", isActive('italic') && "bg-primary/20 text-primary")}
				onclick={() => editor?.chain().focus().toggleItalic().run()}
			>
				<Italic class="w-4 h-4" />
			</button>
			<div class="w-px h-4 bg-slate-200 dark:bg-slate-800 mx-1"></div>
			<button
				class={cn("p-2 rounded-md hover:bg-primary/10 transition-colors", isActive('heading', { level: 1 }) && "bg-primary/20 text-primary")}
				onclick={() => editor?.chain().focus().toggleHeading({ level: 1 }).run()}
			>
				<Heading1 class="w-4 h-4" />
			</button>
			<button
				class={cn("p-2 rounded-md hover:bg-primary/10 transition-colors", isActive('heading', { level: 2 }) && "bg-primary/20 text-primary")}
				onclick={() => editor?.chain().focus().toggleHeading({ level: 2 }).run()}
			>
				<Heading2 class="w-4 h-4" />
			</button>
			<div class="w-px h-4 bg-slate-200 dark:bg-slate-800 mx-1"></div>
			<button
				class={cn("p-2 rounded-md hover:bg-primary/10 transition-colors", isActive('bulletList') && "bg-primary/20 text-primary")}
				onclick={() => editor?.chain().focus().toggleBulletList().run()}
			>
				<List class="w-4 h-4" />
			</button>
			<button
				class={cn("p-2 rounded-md hover:bg-primary/10 transition-colors", isActive('orderedList') && "bg-primary/20 text-primary")}
				onclick={() => editor?.chain().focus().toggleOrderedList().run()}
			>
				<ListOrdered class="w-4 h-4" />
			</button>
			<button
				class={cn("p-2 rounded-md hover:bg-primary/10 transition-colors", isActive('blockquote') && "bg-primary/20 text-primary")}
				onclick={() => editor?.chain().focus().toggleBlockquote().run()}
			>
				<Quote class="w-4 h-4" />
			</button>
		</div>

		<!-- Editor -->
		<div bind:this={element}></div>
	</div>
</div>
