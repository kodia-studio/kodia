<script lang="ts">
  import { UserPlus, Search, Trash2, Pencil } from "lucide-svelte";
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
    if (!confirm("Are you sure you want to delete this user?")) return;
    try {
      await api.delete(`/api/users/${id}`);
      users = users.filter(u => u.id !== id);
      toast.success("User deleted");
    } catch (err: any) {
      toast.error(err.message || "Failed to delete user");
    }
  }

  $effect(() => {
    if (!initialized) {
      fetchUsers();
      initialized = true;
    }
  });

  const filteredUsers = $derived(
    users.filter(u =>
      u.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      u.email.toLowerCase().includes(searchQuery.toLowerCase())
    )
  );
</script>

<div class="space-y-8">
  <!-- Header -->
  <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
    <div>
      <h1 class="text-3xl font-bold text-slate-900 dark:text-white">Users</h1>
      <p class="text-sm text-slate-600 dark:text-slate-400 mt-1">Manage system users and access control</p>
    </div>
    <button class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2 text-sm font-semibold transition-colors">
      <UserPlus class="w-4 h-4" />
      Add User
    </button>
  </div>

  <!-- Search -->
  <div class="relative">
    <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
    <input
      type="text"
      bind:value={searchQuery}
      placeholder="Search by name or email..."
      class="w-full pl-10 pr-4 py-2 bg-white dark:bg-slate-900/50 border border-slate-200 dark:border-slate-800 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
    />
  </div>

  <!-- Users Table -->
  <div class="bg-white dark:bg-slate-900/50 rounded-xl border border-slate-200 dark:border-slate-800 overflow-hidden">
    {#if isLoading}
      <div class="h-64 flex items-center justify-center">
        <div class="text-center">
          <div class="w-10 h-10 rounded-full border-2 border-slate-200 dark:border-slate-700 border-t-blue-500 animate-spin mx-auto mb-3"></div>
          <p class="text-sm text-slate-600 dark:text-slate-400">Loading users...</p>
        </div>
      </div>
    {:else if filteredUsers.length === 0}
      <div class="h-64 flex items-center justify-center">
        <div class="text-center">
          <p class="text-sm text-slate-600 dark:text-slate-400">No users found</p>
        </div>
      </div>
    {:else}
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-800/50">
              <th class="px-6 py-3 text-left text-xs font-semibold text-slate-600 dark:text-slate-400">Name</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-slate-600 dark:text-slate-400">Email</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-slate-600 dark:text-slate-400">Role</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-slate-600 dark:text-slate-400">Joined</th>
              <th class="px-6 py-3 text-right text-xs font-semibold text-slate-600 dark:text-slate-400">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200 dark:divide-slate-800">
            {#each filteredUsers as user (user.id)}
              <tr class="hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                <td class="px-6 py-4">
                  <div class="flex items-center gap-3">
                    <div class="w-9 h-9 rounded-lg bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center text-blue-600 dark:text-blue-400 font-semibold text-sm">
                      {user.name[0].toUpperCase()}
                    </div>
                    <span class="font-medium text-slate-900 dark:text-white">{user.name}</span>
                  </div>
                </td>
                <td class="px-6 py-4 text-sm text-slate-600 dark:text-slate-400">{user.email}</td>
                <td class="px-6 py-4">
                  <span class={cn("inline-block px-2.5 py-1 rounded-full text-xs font-semibold",
                    user.role === 'admin' ? 'bg-orange-100 dark:bg-orange-900/30 text-orange-700 dark:text-orange-400' :
                    user.role === 'moderator' ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400' :
                    'bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300'
                  )}>
                    {user.role}
                  </span>
                </td>
                <td class="px-6 py-4 text-sm text-slate-600 dark:text-slate-400">
                  {new Date(user.createdAt).toLocaleDateString("en-US", { month: 'short', day: 'numeric', year: 'numeric' })}
                </td>
                <td class="px-6 py-4 text-right">
                  <div class="flex items-center justify-end gap-2">
                    <button class="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white">
                      <Pencil class="w-4 h-4" />
                    </button>
                    <button
                      onclick={() => deleteUser(user.id)}
                      class="p-2 hover:bg-red-100 dark:hover:bg-red-900/30 rounded-lg transition-colors text-red-600 dark:text-red-400"
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
