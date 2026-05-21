<script lang="ts">
  import { UserPlus, Search, Trash2, Pencil, Archive, RotateCcw, X } from "lucide-svelte";
  import type { User } from "$lib/types/models.types";
  import { api } from "$lib/api/client.svelte";
  import { toast } from "svelte-sonner";
  import { cn } from "$lib/utils/styles";

  function formatDate(dateString: string | undefined | null): string {
    if (!dateString) return "—";
    try {
      const date = new Date(dateString);
      if (isNaN(date.getTime())) return "—";
      return date.toLocaleDateString("en-US", { month: 'short', day: 'numeric', year: 'numeric' });
    } catch {
      return "—";
    }
  }

  let users = $state<User[]>([]);
  let archivedUsers = $state<User[]>([]);
  let isLoading = $state(true);
  let searchQuery = $state("");
  let initialized = $state(false);
  let activeTab = $state<"active" | "archived">("active");

  // Modal states
  let showModal = $state(false);
  let modalMode = $state<"add" | "edit">("add");
  let editingUser = $state<User | null>(null);
  let formData = $state({
    name: "",
    email: "",
    role: "user"
  });

  // Confirmation modal states
  let showConfirmModal = $state(false);
  let confirmAction = $state<"delete" | "restore" | "force-delete" | null>(null);
  let pendingUser = $state<User | null>(null);
  let isConfirming = $state(false);

  async function fetchUsers() {
    isLoading = true;
    try {
      const res = await api.get<any[]>("/api/users");
      users = res
        .filter(u => !u.deleted_at)
        .map((u: any) => ({
          ...u,
          createdAt: u.created_at,
          updatedAt: u.updated_at,
          deletedAt: u.deleted_at
        }));
    } catch (err: any) {
      toast.error(err.message || "Failed to load users");
    } finally {
      isLoading = false;
    }
  }

  async function fetchArchivedUsers() {
    try {
      const res = await api.get<any[]>("/api/users/trashed");
      archivedUsers = res.map((u: any) => ({
        ...u,
        createdAt: u.created_at,
        updatedAt: u.updated_at,
        deletedAt: u.deleted_at
      }));
    } catch (err: any) {
      toast.error(err.message || "Failed to load archived users");
    }
  }

  function openAddModal() {
    modalMode = "add";
    formData = { name: "", email: "", role: "user" };
    editingUser = null;
    showModal = true;
  }

  function openEditModal(user: User) {
    modalMode = "edit";
    editingUser = user;
    formData = {
      name: user.name,
      email: user.email,
      role: user.role
    };
    showModal = true;
  }

  function closeModal() {
    showModal = false;
    formData = { name: "", email: "", role: "user" };
    editingUser = null;
  }

  async function saveUser() {
    if (!formData.name.trim() || !formData.email.trim()) {
      toast.error("Name and email are required");
      return;
    }

    try {
      if (modalMode === "add") {
        const res = await api.post<User>("/api/users", {
          name: formData.name,
          email: formData.email,
          role: formData.role,
          password: "TempPassword123!" // Default password
        });

        // Ensure created_at is set to today if not provided by API
        const today = new Date().toISOString();
        const newUser = {
          ...res,
          createdAt: res.created_at || res.createdAt || today,
          updatedAt: res.updated_at || res.updatedAt || today
        };

        users = [...users, newUser];
        toast.success("User created successfully");
      } else if (editingUser) {
        await api.patch(`/api/users/${editingUser.id}`, {
          name: formData.name,
          email: formData.email,
          role: formData.role
        });
        users = users.map(u => u.id === editingUser.id ? { ...u, ...formData } : u);
        toast.success("User updated successfully");
      }
      closeModal();
    } catch (err: any) {
      toast.error(err.message || "Failed to save user");
    }
  }

  function openDeleteConfirm(user: User) {
    pendingUser = user;
    confirmAction = "delete";
    showConfirmModal = true;
  }

  function openRestoreConfirm(user: User) {
    pendingUser = user;
    confirmAction = "restore";
    showConfirmModal = true;
  }

  function openForceDeleteConfirm(user: User) {
    pendingUser = user;
    confirmAction = "force-delete";
    showConfirmModal = true;
  }

  function closeConfirmModal() {
    showConfirmModal = false;
    confirmAction = null;
    pendingUser = null;
    isConfirming = false;
  }

  async function executeConfirmedAction() {
    if (!pendingUser || !confirmAction) return;

    isConfirming = true;
    try {
      if (confirmAction === "delete") {
        await api.delete(`/api/users/${pendingUser.id}`);
        users = users.filter(u => u.id !== pendingUser!.id);
        await fetchArchivedUsers();
        toast.success("User moved to archive");
      } else if (confirmAction === "restore") {
        await api.post(`/api/users/${pendingUser.id}/restore`, {});
        archivedUsers = archivedUsers.filter(u => u.id !== pendingUser!.id);
        await fetchUsers();
        toast.success("User restored successfully");
      } else if (confirmAction === "force-delete") {
        await api.post(`/api/users/${pendingUser.id}/force-delete`, {});
        archivedUsers = archivedUsers.filter(u => u.id !== pendingUser!.id);
        toast.success("User permanently deleted");
      }
      closeConfirmModal();
    } catch (err: any) {
      toast.error(err.message || `Failed to ${confirmAction} user`);
    } finally {
      isConfirming = false;
    }
  }

  $effect(() => {
    if (!initialized) {
      fetchUsers();
      fetchArchivedUsers();
      initialized = true;
    }
  });

  const displayedUsers = $derived(
    (activeTab === "active" ? users : archivedUsers).filter(u =>
      u.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      u.email.toLowerCase().includes(searchQuery.toLowerCase())
    )
  );
</script>

<svelte:head>
  <title>Users | Kodia Console</title>
</svelte:head>

<div class="space-y-8">
  <!-- Header -->
  <div class="flex flex-col gap-4">
    <div>
      <h1 class="text-2xl md:text-3xl font-bold text-slate-900 dark:text-white">Users</h1>
      <p class="text-xs md:text-sm text-slate-600 dark:text-slate-400 mt-1">Manage system users and access control</p>
    </div>
    <div class="flex flex-col sm:flex-row sm:items-center gap-3">
      {#if activeTab === "active"}
        <button
          onclick={openAddModal}
          class="btn-premium text-sm flex items-center justify-center sm:justify-start gap-2 w-full sm:w-auto"
        >
          <UserPlus class="w-4 h-4" />
          Add User
        </button>
      {/if}
      <button
        onclick={() => activeTab = activeTab === "active" ? "archived" : "active"}
        class={cn(
          "px-4 py-2 rounded-lg flex items-center justify-center sm:justify-start gap-2 text-sm font-semibold transition-colors w-full sm:w-auto",
          activeTab === "archived"
            ? "btn-premium"
            : "bg-slate-200 dark:bg-slate-800 text-slate-900 dark:text-white hover:bg-slate-300 dark:hover:bg-slate-700"
        )}
      >
        <Archive class="w-4 h-4" />
        <span class="hidden sm:inline">{archivedUsers.length > 0 ? `Archive (${archivedUsers.length})` : "Archive"}</span>
        <span class="sm:hidden">Archive {archivedUsers.length > 0 ? `(${archivedUsers.length})` : ""}</span>
      </button>
    </div>
  </div>

  <!-- Search -->
  <div class="relative">
    <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
    <input
      type="text"
      bind:value={searchQuery}
      placeholder="Search by name or email..."
      class="w-full pl-10 pr-4 py-2 bg-white dark:bg-slate-900/50 border border-slate-200 dark:border-slate-800 rounded-lg text-xs md:text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
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
    {:else if displayedUsers.length === 0}
      <div class="h-64 flex items-center justify-center">
        <div class="text-center">
          <p class="text-sm text-slate-600 dark:text-slate-400">
            {searchQuery ? "No users found" : activeTab === "active" ? "No active users" : "No archived users"}
          </p>
        </div>
      </div>
    {:else}
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-800/50">
              <th class="px-3 md:px-6 py-3 text-left text-[10px] md:text-xs font-semibold text-slate-600 dark:text-slate-400">Name</th>
              <th class="hidden sm:table-cell px-3 md:px-6 py-3 text-left text-[10px] md:text-xs font-semibold text-slate-600 dark:text-slate-400">Email</th>
              <th class="hidden md:table-cell px-3 md:px-6 py-3 text-left text-[10px] md:text-xs font-semibold text-slate-600 dark:text-slate-400">Role</th>
              <th class="hidden lg:table-cell px-3 md:px-6 py-3 text-left text-[10px] md:text-xs font-semibold text-slate-600 dark:text-slate-400">
                {activeTab === "active" ? "Joined" : "Deleted"}
              </th>
              <th class="px-3 md:px-6 py-3 text-right text-[10px] md:text-xs font-semibold text-slate-600 dark:text-slate-400">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200 dark:divide-slate-800">
            {#each displayedUsers as user (user.id)}
              <tr class="hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                <td class="px-3 md:px-6 py-3 md:py-4">
                  <div class="flex items-center gap-2 md:gap-3">
                    <div class="w-8 h-8 md:w-9 md:h-9 rounded-lg bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center text-blue-600 dark:text-blue-400 font-semibold text-xs md:text-sm shrink-0">
                      {user.name[0].toUpperCase()}
                    </div>
                    <div class="min-w-0 flex-1">
                      <p class="font-medium text-slate-900 dark:text-white text-xs md:text-sm truncate">{user.name}</p>
                      <p class="text-[10px] md:hidden text-slate-600 dark:text-slate-400 truncate">{user.email}</p>
                    </div>
                  </div>
                </td>
                <td class="hidden sm:table-cell px-3 md:px-6 py-3 md:py-4 text-xs md:text-sm text-slate-600 dark:text-slate-400 truncate">{user.email}</td>
                <td class="hidden md:table-cell px-3 md:px-6 py-3 md:py-4">
                  <span class={cn("inline-block px-2.5 py-1 rounded-full text-[10px] md:text-xs font-semibold",
                    user.role === 'admin' ? 'bg-orange-100 dark:bg-orange-900/30 text-orange-700 dark:text-orange-400' :
                    user.role === 'moderator' ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400' :
                    'bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300'
                  )}>
                    {user.role}
                  </span>
                </td>
                <td class="hidden lg:table-cell px-3 md:px-6 py-3 md:py-4 text-xs md:text-sm text-slate-600 dark:text-slate-400">
                  {formatDate(activeTab === "active" ? user.createdAt : user.deletedAt)}
                </td>
                <td class="px-3 md:px-6 py-3 md:py-4 text-right">
                  <div class="flex items-center justify-end gap-1 md:gap-2">
                    {#if activeTab === "active"}
                      <button
                        onclick={() => openEditModal(user)}
                        class="p-1.5 md:p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white"
                        title="Edit user"
                      >
                        <Pencil class="w-4 h-4" />
                      </button>
                      <button
                        onclick={() => openDeleteConfirm(user)}
                        class="p-1.5 md:p-2 hover:bg-red-100 dark:hover:bg-red-900/30 rounded-lg transition-colors text-red-600 dark:text-red-400"
                        title="Move to archive"
                      >
                        <Trash2 class="w-4 h-4" />
                      </button>
                    {:else}
                      <button
                        onclick={() => openRestoreConfirm(user)}
                        class="p-1.5 md:p-2 hover:bg-green-100 dark:hover:bg-green-900/30 rounded-lg transition-colors text-green-600 dark:text-green-400"
                        title="Restore user"
                      >
                        <RotateCcw class="w-4 h-4" />
                      </button>
                      <button
                        onclick={() => openForceDeleteConfirm(user)}
                        class="p-1.5 md:p-2 hover:bg-red-100 dark:hover:bg-red-900/30 rounded-lg transition-colors text-red-600 dark:text-red-400"
                        title="Permanently delete"
                      >
                        <Trash2 class="w-4 h-4" />
                      </button>
                    {/if}
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

<!-- Add/Edit User Modal -->
{#if showModal}
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
    <div class="bg-white dark:bg-slate-900 rounded-xl shadow-2xl max-w-md w-full mx-4 max-h-[90vh] overflow-y-auto">
      <div class="flex items-center justify-between p-4 md:p-6 border-b border-slate-200 dark:border-slate-800 sticky top-0 bg-white dark:bg-slate-900">
        <h2 class="text-base md:text-lg font-bold text-slate-900 dark:text-white">
          {modalMode === "add" ? "Add New User" : "Edit User"}
        </h2>
        <button
          onclick={closeModal}
          class="p-1 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
        >
          <X class="w-5 h-5 text-slate-600 dark:text-slate-400" />
        </button>
      </div>

      <div class="p-4 md:p-6 space-y-4">
        <div>
          <label for="modal-name" class="block text-sm font-medium text-slate-900 dark:text-white mb-2">
            Full Name
          </label>
          <input
            id="modal-name"
            type="text"
            bind:value={formData.name}
            placeholder="John Doe"
            class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label for="modal-email" class="block text-sm font-medium text-slate-900 dark:text-white mb-2">
            Email Address
          </label>
          <input
            id="modal-email"
            type="email"
            bind:value={formData.email}
            placeholder="john@example.com"
            class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label for="modal-role" class="block text-sm font-medium text-slate-900 dark:text-white mb-2">
            Role
          </label>
          <select
            id="modal-role"
            bind:value={formData.role}
            class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="user">User</option>
            <option value="moderator">Moderator</option>
            <option value="admin">Admin</option>
          </select>
        </div>

        {#if modalMode === "add"}
          <div class="bg-blue-50 dark:bg-blue-950/30 border border-blue-200 dark:border-blue-900/50 rounded-lg p-3 space-y-2">
            <p class="text-xs font-semibold text-blue-700 dark:text-blue-400">
              🔑 Default Password Set
            </p>
            <p class="text-xs text-blue-600 dark:text-blue-400 font-mono bg-white dark:bg-slate-800 px-2 py-1 rounded break-all">
              TempPassword123!
            </p>
            <p class="text-xs text-blue-600 dark:text-blue-400">
              User must change this password on first login.
            </p>
          </div>
        {/if}
      </div>

      <div class="flex gap-3 p-4 md:p-6 border-t border-slate-200 dark:border-slate-800 sticky bottom-0 bg-white dark:bg-slate-900">
        <button
          onclick={closeModal}
          class="flex-1 px-4 py-2 text-xs md:text-sm bg-slate-200 dark:bg-slate-800 text-slate-900 dark:text-white rounded-lg font-medium hover:bg-slate-300 dark:hover:bg-slate-700 transition-colors"
        >
          Cancel
        </button>
        <button
          onclick={saveUser}
          class="flex-1 btn-premium text-xs md:text-sm"
        >
          {modalMode === "add" ? "Create User" : "Save Changes"}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Confirmation Modal -->
{#if showConfirmModal && pendingUser}
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
    <div class="bg-white dark:bg-slate-900 rounded-2xl shadow-2xl max-w-md w-full overflow-hidden animate-in zoom-in-95 duration-300 max-h-[90vh] overflow-y-auto">
      <!-- Header -->
      <div class={cn("px-4 md:px-6 py-4 border-b sticky top-0",
        confirmAction === "force-delete"
          ? "bg-red-50 dark:bg-red-950/30 border-red-200 dark:border-red-900/50"
          : "bg-slate-50 dark:bg-slate-800/50 border-slate-200 dark:border-slate-800"
      )}>
        <div class="flex items-start gap-2 md:gap-3">
          {#if confirmAction === "delete"}
            <div class="w-8 h-8 md:w-10 md:h-10 rounded-lg bg-orange-100 dark:bg-orange-900/30 flex items-center justify-center shrink-0">
              <Archive class="w-4 h-4 md:w-5 md:h-5 text-orange-600 dark:text-orange-400" />
            </div>
            <div class="min-w-0">
              <h3 class="text-base md:text-lg font-bold text-slate-900 dark:text-white">Move to Archive?</h3>
              <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">User can be restored later</p>
            </div>
          {:else if confirmAction === "restore"}
            <div class="w-8 h-8 md:w-10 md:h-10 rounded-lg bg-green-100 dark:bg-green-900/30 flex items-center justify-center shrink-0">
              <RotateCcw class="w-4 h-4 md:w-5 md:h-5 text-green-600 dark:text-green-400" />
            </div>
            <div class="min-w-0">
              <h3 class="text-base md:text-lg font-bold text-slate-900 dark:text-white">Restore User?</h3>
              <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">User will be moved back to active</p>
            </div>
          {:else if confirmAction === "force-delete"}
            <div class="w-8 h-8 md:w-10 md:h-10 rounded-lg bg-red-100 dark:bg-red-900/30 flex items-center justify-center shrink-0">
              <Trash2 class="w-4 h-4 md:w-5 md:h-5 text-red-600 dark:text-red-400" />
            </div>
            <div class="min-w-0">
              <h3 class="text-base md:text-lg font-bold text-red-600 dark:text-red-400">Permanently Delete?</h3>
              <p class="text-xs text-red-500/75 dark:text-red-400/75 mt-0.5">This cannot be undone</p>
            </div>
          {/if}
        </div>
      </div>

      <!-- Content -->
      <div class="px-4 md:px-6 py-4 space-y-4">
        <div class="space-y-2">
          <p class="text-sm font-medium text-slate-600 dark:text-slate-400">User Details:</p>
          <div class="bg-slate-50 dark:bg-slate-800/50 rounded-lg p-3 space-y-1.5 text-sm">
            <div>
              <span class="text-slate-500 dark:text-slate-400">Name:</span>
              <span class="ml-2 font-medium text-slate-900 dark:text-white">{pendingUser.name}</span>
            </div>
            <div>
              <span class="text-slate-500 dark:text-slate-400">Email:</span>
              <span class="ml-2 font-medium text-slate-900 dark:text-white break-all">{pendingUser.email}</span>
            </div>
            <div>
              <span class="text-slate-500 dark:text-slate-400">Role:</span>
              <span class={cn("ml-2 font-medium px-2 py-0.5 rounded text-xs",
                pendingUser.role === 'admin'
                  ? 'bg-orange-100 dark:bg-orange-900/30 text-orange-700 dark:text-orange-400' :
                pendingUser.role === 'moderator'
                  ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400'
                  : 'bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300'
              )}>
                {pendingUser.role}
              </span>
            </div>
          </div>
        </div>

        {#if confirmAction === "force-delete"}
          <div class="p-3 bg-red-50 dark:bg-red-950/20 border border-red-200 dark:border-red-900/30 rounded-lg">
            <p class="text-xs text-red-700 dark:text-red-400 font-medium">⚠️ Warning: This action is permanent and cannot be reversed. All data associated with this user will be permanently deleted.</p>
          </div>
        {:else if confirmAction === "delete"}
          <p class="text-sm text-slate-600 dark:text-slate-400">
            The user will be moved to the archive. You can restore them later if needed.
          </p>
        {:else if confirmAction === "restore"}
          <p class="text-sm text-slate-600 dark:text-slate-400">
            The user will be moved back to the active users list and can log in again.
          </p>
        {/if}
      </div>

      <!-- Actions -->
      <div class="flex gap-3 p-4 md:p-6 border-t border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/20 sticky bottom-0">
        <button
          onclick={closeConfirmModal}
          disabled={isConfirming}
          class="flex-1 px-4 py-2 text-xs md:text-sm bg-slate-200 dark:bg-slate-800 text-slate-900 dark:text-white rounded-lg font-medium hover:bg-slate-300 dark:hover:bg-slate-700 transition-colors disabled:opacity-50"
        >
          Cancel
        </button>
        <button
          onclick={executeConfirmedAction}
          disabled={isConfirming}
          class={cn(
            "flex-1 rounded-lg font-medium text-xs md:text-sm transition-colors flex items-center justify-center gap-2 disabled:opacity-50 px-4 py-2",
            confirmAction === "force-delete"
              ? "bg-red-600 text-white hover:bg-red-700 active:scale-95"
              : confirmAction === "restore"
              ? "bg-green-600 text-white hover:bg-green-700 active:scale-95"
              : "btn-premium"
          )}
        >
          {#if isConfirming}
            <div class="w-4 h-4 rounded-full border-2 border-white/30 border-t-white animate-spin"></div>
          {/if}
          {confirmAction === "force-delete" ? "Delete Permanently" : confirmAction === "restore" ? "Restore User" : "Move to Archive"}
        </button>
      </div>
    </div>
  </div>
{/if}
