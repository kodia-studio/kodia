<script lang="ts">
  import { authStore } from "$lib/stores/auth.store";
  import { slide } from "svelte/transition";

  interface Props {
    role?: string | string[];
    permission?: string | string[];
    fallback?: import('svelte').Snippet;
    children: import('svelte').Snippet;
  }

  let { role, permission, fallback, children }: Props = $props();

  const isAuthorized = $derived.by(() => {
    if (!$authStore.isAuthenticated) return false;
    
    const userRole = $authStore.user?.role;
    const userPerms = $authStore.user?.permissions || [];

    // Check Role
    if (role) {
      const roles = Array.isArray(role) ? role : [role];
      if (!userRole || !roles.includes(userRole)) return false;
    }

    // Check Permission
    if (permission) {
      const perms = Array.isArray(permission) ? permission : [permission];
      const hasPerm = perms.some(p => userPerms.includes(p));
      if (!hasPerm) return false;
    }

    return true;
  });
</script>

{#if isAuthorized}
  {@render children()}
{:else if fallback}
  {@render fallback()}
{/if}
