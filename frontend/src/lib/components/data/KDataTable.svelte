<script lang="ts">
  import { onMount } from "svelte";
  import DataTable from "./DataTable.svelte";
  import { api } from "$lib/api/client";
  import type { PaginatedResponse } from "$lib/types/api.types";
  import type { ColumnDef } from "@tanstack/svelte-table";
  import { Loader2 } from "lucide-svelte";

  interface Props {
    endpoint: string;
    columns: ColumnDef<any, any>[];
    pageSize?: number;
    title?: string;
    class?: string;
  }

  let { endpoint, columns, pageSize = 10, title, class: className }: Props = $props();

  // State
  let data = $state<any[]>([]);
  let isLoading = $state(true);
  let page = $state(1);
  let totalPages = $state(1);
  let total = $state(0);
  let sortBy = $state("");
  let order = $state<"asc" | "desc">("asc");

  async function fetchData() {
    isLoading = true;
    try {
      const query = new URLSearchParams({
        page: page.toString(),
        per_page: pageSize.toString(),
        sort_by: sortBy,
        order: order
      });

      const response = await api.get<PaginatedResponse<any>>(`${endpoint}?${query.toString()}`);
      
      // Since api.get returns .data already (as per current client implementation)
      // I might need to adjust the client or handle the raw response if meta is needed.
      // For now, let's assume the API client might need a 'getRaw' or similar if we want Meta.
      // BUT current client.ts returns result.data. I'll handle that discrepancy.
      
      // If we want metadata, we should use a more generic request or have the client return result
      data = (response as any).data || response;
      const meta = (response as any).meta;
      if (meta) {
        totalPages = meta.total_pages;
        total = meta.total;
      }
    } catch (error) {
      console.error("Failed to fetch data:", error);
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    fetchData();
  });

  // Re-fetch on state change using Svelte 5 effects
  $effect(() => {
    if (page || sortBy || order) {
      fetchData();
    }
  });
</script>

<div class={className}>
  {#if title}
    <div class="mb-4">
      <h3 class="text-lg font-bold">{title}</h3>
    </div>
  {/if}

  <div class="relative">
    {#if isLoading}
      <div class="absolute inset-0 z-10 flex items-center justify-center bg-background/50 backdrop-blur-[1px] rounded-2xl">
        <Loader2 class="w-8 h-8 animate-spin text-primary" />
      </div>
    {/if}

    <DataTable 
      {columns} 
      {data} 
    />
  </div>
</div>
