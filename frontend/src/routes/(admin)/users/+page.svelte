<script lang="ts">
  import AdminLayout from "$lib/components/layouts/AdminLayout.svelte";
  import DataTable from "$lib/components/data/DataTable.svelte";
  import { Shield, Mail, MoreHorizontal, UserPlus, Filter, Trash2, Edit2, Loader2, Search, Activity } from "lucide-svelte";
  import type { ColumnDef } from "@tanstack/table-core";
  import type { User } from "$lib/types/models.types";
  import { api } from "$lib/api/client.svelte";
  import { toast } from "svelte-sonner";
  import { cn } from "$lib/utils/styles";

  let users = $state<User[]>([]);
  let isLoading = $state(true);
  let searchQuery = $state("");
  let initialized = $state(false);

  async function fetchUsers() {
    isLoading = true;
    try {
      const res = await api.get<any[]>("/api/users");
      // API client returns data directly, no .data wrapper
      users = res.map((u: any) => ({
        ...u,
        createdAt: u.created_at,
        updatedAt: u.updated_at
      }));
    } catch (err: any) {
      toast.error(err.message || "Failed to load users");
    } finally {
      isLoading = false;
    }
  }

  async function deleteUser(id: string) {
    if (!confirm("Are you sure you want to delete this user? This action is permanent.")) return;
    try {
      await api.delete(`/api/users/${id}`);
      users = users.filter(u => u.id !== id);
      toast.success("User deleted");
    } catch (err: any) {
      toast.error(err.message || "Failed to delete user");
    }
  }

  const columns: ColumnDef<User>[] = [
    {
      accessorKey: "name",
      header: "User",
      cell: (info) => {
        const user = info.row.original;
        return `
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-lg bg-primary/10 flex items-center justify-center text-primary font-bold uppercase">
              ${user.name[0]}
            </div>
            <div>
              <p class="font-bold leading-none capitalize">${user.name}</p>
              <p class="text-xs text-muted-foreground mt-1">${user.email}</p>
            </div>
          </div>
        `;
      }
    },
    {
      accessorKey: "role",
      header: "Role",
      cell: (info) => {
        const role = info.getValue() as string;
        const colors: Record<string, string> = {
          admin: "bg-purple-500/10 text-purple-500",
          moderator: "bg-blue-500/10 text-blue-500",
          user: "bg-muted text-muted-foreground"
        };
        return `<span class="px-2.5 py-1 rounded-lg text-[10px] font-bold uppercase tracking-wider ${colors[role] || colors.user}">${role}</span>`;
      }
    },
    {
      accessorKey: "createdAt",
      header: "Joined",
      cell: (info) => new Date(info.getValue() as string).toLocaleDateString("en-US", { month: 'short', day: 'numeric', year: 'numeric' })
    },
    {
      id: "actions",
      header: "",
      cell: (info) => {
        const user = info.row.original;
        // Don't allow deleting yourself if you browse the table
        return `
          <div class="flex items-center justify-end gap-2">
            <button class="p-2 hover:bg-muted rounded-lg transition-colors text-muted-foreground hover:text-foreground">
              <i data-lucide="edit-2" class="w-4 h-4"></i>
            </button>
            <button 
              class="p-2 hover:bg-destructive/10 rounded-lg transition-colors text-muted-foreground hover:text-destructive"
              onclick="window.dispatchEvent(new CustomEvent('delete-user', { detail: '${user.id}' }))"
            >
              <i data-lucide="trash-2" class="w-4 h-4"></i>
            </button>
          </div>
        `;
      }
    }
  ];

  $effect(() => {
    if (!initialized) {
      fetchUsers();
      initialized = true;
    }
    // Handler for the custom event from the table
    const handleEvent = (e: any) => deleteUser(e.detail);
    window.addEventListener('delete-user', handleEvent);
    return () => window.removeEventListener('delete-user', handleEvent);
  });

  const filteredUsers = $derived(
    users.filter(u => 
      u.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
      u.email.toLowerCase().includes(searchQuery.toLowerCase())
    )
  );
</script>

<AdminLayout>
  <div class="space-y-12 pb-10">
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-6 pb-6 border-b border-slate-200/50 dark:border-white/5">
      <div>
        <h1 class="text-4xl font-black font-heading tracking-tight text-slate-900 dark:text-white leading-none">Colony Registry.</h1>
        <p class="text-xs font-black uppercase tracking-[0.3em] text-slate-400 mt-3 leading-none">Access Control & Intelligence</p>
      </div>

      <div class="flex items-center gap-4">
        <div class="relative hidden lg:block group">
          <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 group-focus-within:text-primary transition-colors" />
          <input 
            type="text" 
            bind:value={searchQuery}
            placeholder="Filter identity..." 
            class="pl-11 pr-4 py-2.5 bg-slate-100/50 dark:bg-white/5 border border-slate-200/50 dark:border-white/5 rounded-xl text-sm focus:ring-4 focus:ring-primary/10 focus:border-primary/30 outline-none transition-all w-72 font-medium"
          />
        </div>
        <button class="btn-premium py-3 px-8 flex items-center gap-3 text-xs uppercase tracking-widest group shadow-lg shadow-primary/20">
          <UserPlus class="w-4.5 h-4.5" />
          Join Colony.
        </button>
      </div>
    </div>

    <!-- Stats Preview (Minor) -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-6">
       <div class="glass p-5 rounded-2xl border border-slate-200/50 dark:border-white/5 flex items-center gap-4">
          <div class="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center text-primary border border-primary/20">
             <Shield class="w-5 h-5" />
          </div>
          <div>
             <p class="text-[10px] font-black uppercase tracking-widest text-slate-400">Security Clearance</p>
             <p class="text-lg font-black text-slate-900 dark:text-white leading-none mt-1">Institutional</p>
          </div>
       </div>
       <div class="glass p-5 rounded-2xl border border-slate-200/50 dark:border-white/5 flex items-center gap-4">
          <div class="w-10 h-10 rounded-xl bg-secondary/10 flex items-center justify-center text-secondary border border-secondary/20">
             <Activity class="w-5 h-5" />
          </div>
          <div>
             <p class="text-[10px] font-black uppercase tracking-widest text-slate-400">Active Sessions</p>
             <p class="text-lg font-black text-slate-900 dark:text-white leading-none mt-1">1,204</p>
          </div>
       </div>
    </div>

    <!-- Data Table Showcase -->
    <div class="relative group">
      <div class="absolute -inset-1 bg-linear-to-r from-primary/20 to-secondary/20 rounded-4xl blur opacity-50 transition duration-500"></div>
      
      <div class="relative glass rounded-4xl border border-slate-200/50 dark:border-white/5 overflow-hidden bg-white/40 dark:bg-slate-900/40 backdrop-blur-3xl shadow-2xl">
        {#if isLoading}
          <div class="h-96 flex flex-col items-center justify-center gap-6 text-slate-400">
            <div class="relative">
              <div class="w-16 h-16 rounded-full border-4 border-slate-100 dark:border-white/5 border-t-primary animate-spin"></div>
              <div class="absolute inset-0 flex items-center justify-center">
                 <div class="w-2 h-2 bg-primary rounded-full animate-pulse"></div>
              </div>
            </div>
            <p class="text-[10px] font-black uppercase tracking-[0.2em] animate-pulse">Scanning Identity Matrix...</p>
          </div>
        {:else}
          <div class="overflow-x-auto selection:bg-primary/20">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="border-b border-slate-200/50 dark:border-white/5 bg-slate-50/50 dark:bg-white/5">
                  <th class="px-8 py-5 text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">User Profile</th>
                  <th class="px-8 py-5 text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">Security Rank</th>
                  <th class="px-8 py-5 text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">Joined Sequence</th>
                  <th class="px-8 py-5 text-[10px] font-black uppercase tracking-[0.2em] text-slate-400 text-right">Actions</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 dark:divide-white/5">
                {#each filteredUsers as user}
                  {@const roleColors = user.role === 'admin' ? 'bg-orange-500/20 text-orange-500 border-orange-500/30' : user.role === 'moderator' ? 'bg-primary/20 text-primary border-primary/30' : 'bg-slate-500/20 text-slate-500 border-slate-500/30'}
                  <tr class="hover:bg-primary/5 transition-colors group/row">
                    <td class="px-8 py-6">
                      <div class="flex items-center gap-4">
                        <div class="w-11 h-11 rounded-2xl bg-linear-to-br from-primary/10 to-secondary/10 flex items-center justify-center text-primary font-black uppercase border border-primary/20 shadow-inner group-hover/row:scale-110 transition-transform">
                          {user.name[0]}
                        </div>
                        <div>
                          <p class="text-sm font-black text-slate-900 dark:text-white leading-none uppercase tracking-tight">{user.name}</p>
                          <p class="text-[10px] text-slate-400 font-bold mt-2 uppercase tracking-widest">{user.email}</p>
                        </div>
                      </div>
                    </td>
                    <td class="px-8 py-6">
                      <span class={cn("px-3 py-1 rounded-full text-[9px] font-black uppercase tracking-widest border", roleColors)}>
                        {user.role}
                      </span>
                    </td>
                    <td class="px-8 py-6">
                       <p class="text-xs font-black text-slate-500 dark:text-slate-400 tabular-nums uppercase">
                         {new Date(user.createdAt).toLocaleDateString("en-US", { month: 'short', day: 'numeric', year: 'numeric' })}
                       </p>
                    </td>
                    <td class="px-8 py-6">
                      <div class="flex items-center justify-end gap-3 opacity-0 group-hover/row:opacity-100 transition-opacity">
                        <button class="w-9 h-9 rounded-xl bg-slate-100 dark:bg-white/5 border border-slate-200/50 dark:border-white/5 flex items-center justify-center text-slate-400 hover:text-primary hover:border-primary/30 transition-all">
                          <Edit2 class="w-4 h-4" />
                        </button>
                        <button 
                          onclick={() => deleteUser(user.id)}
                          class="w-9 h-9 rounded-xl bg-rose-500/10 border border-rose-500/20 flex items-center justify-center text-rose-500 hover:bg-rose-500 hover:text-white transition-all"
                        >
                          <Trash2 class="w-4 h-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}
      </div>
    </div>
  </div>
</AdminLayout>
