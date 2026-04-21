<script lang="ts">
  import AuthLayout from "$lib/components/layouts/AuthLayout.svelte";
  import FormField from "$lib/components/forms/FormField.svelte";
  import Input from "$lib/components/forms/Input.svelte";
  import Checkbox from "$lib/components/forms/Checkbox.svelte";
  import { Mail, Lock, LogIn, Loader2, ArrowRight } from "lucide-svelte";
  import { api } from "$lib/api/client.svelte";
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
      const response = await api.post<any>("/auth/login", { email, password });
      
      if (response.mfa_required) {
        sessionStorage.setItem("mfa_token", response.mfa_token);
        toast.info("Two-Factor Authentication Required");
        goto("/2fa");
        return;
      }

      const { user, access_token } = response;
      authStore.login(user, access_token);
      toast.success("Welcome back, " + user.name);
      goto("/dashboard");
    } catch (err: any) {
      toast.error(err.message || "Invalid credentials");
    } finally {
      isLoading = false;
    }
  }
</script>

<AuthLayout>
  <form onsubmit={handleLogin} class="space-y-6">
    <div class="text-center mb-10">
      <h1 class="text-4xl font-black tracking-tight text-slate-900 dark:text-white mb-2">Sign In.</h1>
      <p class="text-sm font-medium text-slate-500 dark:text-slate-400">Enter your colony credentials to access the hub.</p>
    </div>

    <div class="space-y-4">
      <FormField label="Email Address">
        <Input 
          type="email" 
          bind:value={email} 
          placeholder="name@example.com"
          required
        >
          {#snippet icon()}
            <Mail size={18} />
          {/snippet}
        </Input>
      </FormField>

      <FormField label="Password">
        <Input 
          type="password" 
          bind:value={password} 
          placeholder="••••••••"
          required
        >
          {#snippet icon()}
            <Lock size={18} />
          {/snippet}
        </Input>
      </FormField>
    </div>

    <div class="flex items-center justify-between">
      <Checkbox bind:checked={rememberMe} label="Remember me" />
      <a href="/forgot-password" class="text-xs font-bold text-primary hover:text-primary/80 transition-colors">
        Forgot password?
      </a>
    </div>

    <button
      type="submit"
      disabled={isLoading}
      class="btn-premium w-full py-4 flex items-center justify-center gap-3 group text-lg"
    >
      {#if isLoading}
        <Loader2 class="w-6 h-6 animate-spin text-white" />
      {:else}
        <LogIn size={20} class="group-hover:translate-x-1 transition-transform" />
      {/if}
      Sign In to Hive
    </button>

    <div class="pt-6 border-t border-slate-100 dark:border-slate-800 text-center">
      <p class="text-xs font-bold text-slate-500 dark:text-slate-400">
        Don't have an account yet? 
        <a href="/register" class="inline-flex items-center gap-1 text-primary hover:underline ml-1">
          Create One
          <ArrowRight size={10} />
        </a>
      </p>
    </div>
  </form>
</AuthLayout>
