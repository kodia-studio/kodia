<script lang="ts">
  import AdminLayout from "$lib/components/layouts/AdminLayout.svelte";
  import Container from "$lib/components/shared/Container.svelte";
  import DataTable from "$lib/components/data/DataTable.svelte";
  import { Shield, Mail, MoreHorizontal, UserPlus, Filter } from "lucide-svelte";
  import type { ColumnDef } from "@tanstack/svelte-table";
  import type { User } from "$lib/types/models.types";

  const columns: ColumnDef<User>[] = [
    {
      accessorKey: "name",
      header: "User",
      cell: (info) => {
        const user = info.row.original;
        return `
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-lg bg-primary/10 flex items-center justify-center text-primary font-bold">
              ${user.name[0]}
            </div>
            <div>
              <p class="font-bold leading-none">${user.name}</p>
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
          admin: "bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400",
          moderator: "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400",
          user: "bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-400"
        };
        return `<span class="px-2.5 py-1 rounded-lg text-xs font-bold uppercase tracking-wider ${colors[role] || colors.user}">${role}</span>`;
      }
    },
    {
      accessorKey: "createdAt",
      header: "Joined",
      cell: (info) => new Date(info.getValue() as string).toLocaleDateString()
    },
    {
      id: "actions",
      header: "",
      cell: () => `
        <button class="p-2 hover:bg-muted rounded-lg transition-colors">
          <MoreHorizontal class="w-4 h-4 text-muted-foreground" />
        </button>
      `
    }
  ];

  // Mock data
  const users: User[] = [
    { id: "1", name: "Andi Aryatno", email: "andi@kodia.dev", role: "admin", permissions: ["*"], createdAt: "2024-01-15", updatedAt: "2024-04-19" },
    { id: "2", name: "Jane Smith", email: "jane@example.com", role: "moderator", permissions: ["users:read"], createdAt: "2024-02-10", updatedAt: "2024-04-10" },
    { id: "3", name: "Bob Wilson", email: "bob@test.com", role: "user", permissions: [], createdAt: "2024-03-05", updatedAt: "2024-03-05" },
    { id: "4", name: "Alice Brown", email: "alice@kodia.dev", role: "user", permissions: [], createdAt: "2024-03-20", updatedAt: "2024-03-20" },
    { id: "5", name: "Charlie Davis", email: "charlie@kodia.dev", role: "user", permissions: [], createdAt: "2024-04-01", updatedAt: "2024-04-01" },
  ];
</script>

<AdminLayout>
  <div class="space-y-8">
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-3xl font-bold font-heading tracking-tight">User Management</h1>
        <p class="text-muted-foreground mt-1.5">Manage and monitor access across your Kodia application.</p>
      </div>

      <div class="flex items-center gap-3">
        <button class="flex items-center gap-2 px-4 py-2.5 bg-muted/50 hover:bg-muted border rounded-xl text-sm font-bold transition-all active:scale-95">
          <Filter class="w-4 h-4" />
          Filter
        </button>
        <button class="btn-primary py-2.5 px-6 flex items-center gap-2">
          <UserPlus class="w-5 h-5" />
          Add User
        </button>
      </div>
    </div>

    <!-- Data Table -->
    <DataTable {columns} data={users} />
  </div>
</AdminLayout>
