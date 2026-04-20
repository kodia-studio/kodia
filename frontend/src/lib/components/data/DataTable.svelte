<script lang="ts">
  import { 
    createSvelteTable, 
    flexRender, 
    getCoreRowModel, 
    getPaginationRowModel,
    getSortedRowModel,
    getFilteredRowModel,
    type ColumnDef, 
    type SortingState,
    type TableOptions 
  } from "@tanstack/svelte-table";
  import { cn } from "$lib/utils/styles";
  import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight, ArrowUpDown } from "lucide-svelte";

  interface Props<TData, TValue> {
    columns: ColumnDef<TData, TValue>[];
    data: TData[];
    class?: string;
  }

  let { columns, data, class: className }: Props<any, any> = $props();

  let sorting = $state<SortingState>([]);
  let pagination = $state({ pageIndex: 0, pageSize: 10 });

  const table = createSvelteTable({
    get data() { return data; },
    get columns() { return columns; },
    state: {
      get sorting() { return sorting; },
      get pagination() { return pagination; },
    },
    onSortingChange: (updater) => {
      if (typeof updater === "function") {
        sorting = updater(sorting);
      } else {
        sorting = updater;
      }
    },
    onPaginationChange: (updater) => {
      if (typeof updater === "function") {
        pagination = updater(pagination);
      } else {
        pagination = updater;
      }
    },
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
  });
</script>

<div class={cn("w-full space-y-4", className)}>
  <div class="rounded-2xl border bg-card overflow-hidden shadow-sm">
    <table class="w-full text-sm">
      <thead class="bg-muted/50 border-b">
        {#each $table.getHeaderGroups() as headerGroup}
          <tr>
            {#each headerGroup.headers as header}
              <th class="h-12 px-4 text-left align-middle font-semibold text-muted-foreground">
                {#if !header.isPlaceholder}
                  <button
                    class={cn(
                      "flex items-center gap-2 hover:text-foreground transition-colors",
                      header.column.getCanSort() && "cursor-pointer select-none"
                    )}
                    onclick={header.column.getToggleSortingHandler()}
                  >
                    {flexRender(header.column.columnDef.header, header.getContext())}
                    {#if header.column.getCanSort()}
                      <ArrowUpDown class="w-3.5 h-3.5" />
                    {/if}
                  </button>
                {/if}
              </th>
            {/each}
          </tr>
        {/each}
      </thead>
      <tbody class="divide-y">
        {#each $table.getRowModel().rows as row}
          <tr class="hover:bg-muted/30 transition-colors">
            {#each row.getVisibleCells() as cell}
              <td class="p-4 align-middle">
                {flexRender(cell.column.columnDef.cell, cell.getContext())}
              </td>
            {/each}
          </tr>
        {:else}
          <tr>
            <td colspan={columns.length} class="h-24 text-center text-muted-foreground">
              No results found.
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>

  <!-- Pagination -->
  <div class="flex items-center justify-between px-2">
    <div class="text-sm text-muted-foreground">
      Page {pagination.pageIndex + 1} of {$table.getPageCount()}
    </div>
    <div class="flex items-center gap-2">
      <button
        class="p-2 border rounded-xl hover:bg-muted disabled:opacity-50 transition-colors"
        onclick={() => $table.setPageIndex(0)}
        disabled={!$table.getCanPreviousPage()}
      >
        <ChevronsLeft class="w-4 h-4" />
      </button>
      <button
        class="p-2 border rounded-xl hover:bg-muted disabled:opacity-50 transition-colors"
        onclick={() => $table.previousPage()}
        disabled={!$table.getCanPreviousPage()}
      >
        <ChevronLeft class="w-4 h-4" />
      </button>
      <button
        class="p-2 border rounded-xl hover:bg-muted disabled:opacity-50 transition-colors"
        onclick={() => $table.nextPage()}
        disabled={!$table.getCanNextPage()}
      >
        <ChevronRight class="w-4 h-4" />
      </button>
      <button
        class="p-2 border rounded-xl hover:bg-muted disabled:opacity-50 transition-colors"
        onclick={() => $table.setPageIndex($table.getPageCount() - 1)}
        disabled={!$table.getCanNextPage()}
      >
        <ChevronsRight class="w-4 h-4" />
      </button>
    </div>
  </div>
</div>
