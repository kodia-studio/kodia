<script lang="ts" module>
	import {
		getCoreRowModel,
		getPaginationRowModel,
		getSortedRowModel,
		type ColumnDef,
		type SortingState,
		type Updater,
	} from "@tanstack/table-core";

	export function flexRender(content: any, context: any) {
		if (typeof content === "function") {
			return content(context);
		}
		return content;
	}
</script>

<script lang="ts">
	import { createSvelteTable } from "$lib/components/ui/table/table.svelte";
	import { cn } from "$lib/utils";
	import Button from "../ui/Button.svelte";
	import { ChevronLeft, ChevronRight, ChevronsUpDown, ChevronUp, ChevronDown } from "lucide-svelte";

	interface Props<TData, TValue> {
		columns: ColumnDef<TData, TValue>[];
		data: TData[];
		class?: string;
	}

	let { data, columns, class: className }: Props<any, any> = $props();

	let sorting = $state<SortingState>([]);

	const handleSortingChange = (updater: Updater<SortingState>) => {
		if (typeof updater === "function") {
			sorting = updater(sorting);
		} else {
			sorting = updater;
		}
	};

	const table = $derived.by(() => {
		return createSvelteTable({
			data,
			columns,
			state: { sorting },
			onSortingChange: handleSortingChange,
			getCoreRowModel: getCoreRowModel(),
			getPaginationRowModel: getPaginationRowModel(),
			getSortedRowModel: getSortedRowModel(),
		});
	});
</script>

<div class={cn("flex flex-col gap-4", className)}>
	<div class="rounded-kodia-lg border border-slate-200 dark:border-slate-800 bg-white/50 dark:bg-slate-900/50 backdrop-blur-sm overflow-hidden shadow-kodia">
		<table class="w-full text-left text-sm border-collapse">
			<thead>
				{#each table.getHeaderGroups() as headerGroup}
					<tr class="border-b border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/80">
						{#each headerGroup.headers as header}
							<th class="px-4 py-4 font-black uppercase tracking-widest text-[10px] text-slate-500">
								{#if !header.isPlaceholder}
									<button
										class={cn(
											"flex items-center gap-1 hover:text-primary transition-colors uppercase tracking-widest",
											header.column.getCanSort() && "cursor-pointer select-none"
										)}
										onclick={header.column.getToggleSortingHandler()}
									>
										{flexRender(header.column.columnDef.header, header.getContext())}
										{#if header.column.getCanSort()}
											{#if header.column.getIsSorted() === "asc"}
												<ChevronUp class="w-3 h-3 text-primary" />
											{:else if header.column.getIsSorted() === "desc"}
												<ChevronDown class="w-3 h-3 text-primary" />
											{:else}
												<ChevronsUpDown class="w-3 h-3 opacity-30" />
											{/if}
										{/if}
									</button>
								{/if}
							</th>
						{/each}
					</tr>
				{/each}
			</thead>
			<tbody>
				{#each table.getRowModel().rows as row}
					<tr class="border-b border-slate-100 dark:border-slate-800 hover:bg-primary/5 transition-colors">
						{#each row.getVisibleCells() as cell}
							<td class="px-4 py-4 text-slate-700 dark:text-slate-300 font-medium">
								{flexRender(cell.column.columnDef.cell, cell.getContext())}
							</td>
						{/each}
					</tr>
				{:else}
					<tr>
						<td colspan={columns.length} class="px-4 py-12 text-center text-slate-400 font-bold uppercase tracking-widest italic">
							No results found.
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>

	<!-- Pagination Controls -->
	<div class="flex items-center justify-between px-2">
		<div class="text-[10px] font-black uppercase tracking-widest text-slate-500">
			Page {table.getState().pagination.pageIndex + 1} of {table.getPageCount()}
		</div>
		<div class="flex items-center gap-2">
			<Button
				variant="outline"
				size="icon"
				onclick={() => table.previousPage()}
				disabled={!table.getCanPreviousPage()}
			>
				<ChevronLeft class="w-4 h-4" />
			</Button>
			<Button
				variant="outline"
				size="icon"
				onclick={() => table.nextPage()}
				disabled={!table.getCanNextPage()}
			>
				<ChevronRight class="w-4 h-4" />
			</Button>
		</div>
	</div>
</div>
