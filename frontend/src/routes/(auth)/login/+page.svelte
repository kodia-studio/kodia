<script lang="ts">
  import AuthLayout from "$lib/components/layouts/AuthLayout.svelte";
  import FormField from "$lib/components/forms/FormField.svelte";
  import Input from "$lib/components/forms/Input.svelte";
  import Checkbox from "$lib/components/forms/Checkbox.svelte";
  import { Mail, Lock, LogIn, Loader2 } from "lucide-svelte";
  import { api } from "$lib/api/client";
  import { authStore } from "$lib/stores/auth.store";
  import { toast } from "svelte-sonner";
  import { goto } from "$app/navigation";

  let email = $state("");
  let password = $state("");
  let rememberMe = $state(false);
  let isLoading = $state(false);

  async function handleLogin(e: SubmitEvent) {
    e.preventDefault();
    isLoading = true;

    try {
      // In a real app, we'd use superforms. For this template:
      const response = await api.post<import("$lib/types/auth.types").AuthResponse>("/auth/login", { email, password });
      authStore.login(response.user, response.access_token);
      toast.success("Welcome back, " + response.user.name);
      goto("/admin");
    } catch (err: any) {
      toast.error(err.message || "Invalid credentials");
    } finally {
      isLoading = false;
    }
  }
</script>

<AuthLayout>
  <form onsubmit={handleLogin} class="space-y-6">
    <div class="text-center mb-8">
      <h1 class="text-2xl font-bold tracking-tight">Sign in to Account</h1>
      <p class="text-sm text-muted-foreground mt-2">Enter your credentials to access your dashboard</p>
    </div>

    <FormField label="Email Address" required>
      <Input 
        type="email" 
        bind:value={email} 
        placeholder="name@example.com"
        required
      >
        {#snippet icon()}
          <Mail class="w-5 h-5" />
        {/snippet}
      </Input>
    </FormField>

    <FormField label="Password" required>
      <Input 
        type="password" 
        bind:value={password} 
        placeholder="••••••••"
        required
      >
        {#snippet icon()}
          <Lock class="w-5 h-5" />
        {/snippet}
      </Input>
    </FormField>

    <div class="flex items-center justify-between">
      <Checkbox bind:checked={rememberMe} label="Remember me" />
      <a href="/auth/reset" class="text-sm font-semibold text-primary hover:underline">
        Forgot password?
      </a>
    </div>

    <button
      type="submit"
      disabled={isLoading}
      class="btn-primary w-full py-3 flex items-center justify-center gap-2 group"
    >
      {#if isLoading}
        <Loader2 class="w-5 h-5 animate-spin" />
      {:else}
        <LogIn class="w-5 h-5 group-hover:translate-x-1 transition-transform" />
      {/if}
      Sign In
    </button>

    <div class="text-center mt-6">
      <p class="text-sm text-muted-foreground">
        Don't have an account? 
        <a href="/auth/register" class="font-semibold text-primary hover:underline">Create one</a>
      </p>
    </div>
  </form>
</AuthLayout>
