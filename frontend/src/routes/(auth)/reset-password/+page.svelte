<script lang="ts">
  import AuthLayout from "$lib/components/layouts/AuthLayout.svelte";
  import FormField from "$lib/components/forms/FormField.svelte";
  import Input from "$lib/components/forms/Input.svelte";
  import { Lock, Loader2, KeyRound, CheckCircle2, ArrowRight } from "lucide-svelte";
  import { api } from "$lib/api/client.svelte";
  import { toast } from "svelte-sonner";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";

  let password = $state("");
  let confirmPassword = $state("");
  let isLoading = $state(false);
  let isSuccess = $state(false);

  const token = page.url.searchParams.get("token") || "";

  async function handleReset(e: SubmitEvent) {
    e.preventDefault();
    if (!token) {
      toast.error("Invalid or missing reset token");
      return;
    }
    if (password !== confirmPassword) {
      toast.error("Passwords do not match");
      return;
    }

    isLoading = true;
    try {
      await api.post("/auth/reset-password", { token, new_password: password });
      isSuccess = true;
      toast.success("Password reset successfully!");
      setTimeout(() => goto("/login"), 3000); // FIXED: removed /auth/
    } catch (err: any) {
      toast.error(err.message || "Failed to reset password");
    } finally {
      isLoading = false;
    }
  }
</script>

<AuthLayout>
  {#if isSuccess}
    <div class="text-center py-10">
      <div class="w-16 h-16 rounded-2xl bg-emerald-500/10 text-emerald-500 mx-auto flex items-center justify-center mb-6 shadow-2xl shadow-emerald-500/20 rotate-12">
        <CheckCircle2 size={32} />
      </div>
      <h1 class="text-4xl font-black tracking-tight text-slate-900 dark:text-white mb-4">Secure.</h1>
      <p class="text-sm font-medium text-slate-500 dark:text-slate-400 mb-10 leading-relaxed">
        Your password has been reset successfully. <br/>
        Taking you back to the colony...
      </p>
      
      <button onclick={() => goto("/login")} class="btn-premium w-full py-4 inline-flex items-center justify-center gap-3 text-lg font-bold">
        Sign In Now
        <ArrowRight size={20} />
      </button>
    </div>
  {:else}
    <form onsubmit={handleReset} class="space-y-6">
      <div class="text-center mb-10">
        <h1 class="text-4xl font-black tracking-tight text-slate-900 dark:text-white mb-2">Secure.</h1>
        <p class="text-sm font-medium text-slate-500 dark:text-slate-400">Choose a strong new password to protect your hive access.</p>
      </div>

      {#if !token}
        <div class="p-4 rounded-2xl bg-rose-500/10 border border-rose-500/20 text-rose-500 text-xs font-bold flex items-start gap-4 mb-6">
          <KeyRound size={18} class="shrink-0" />
          <p>This reset link is invalid or has expired. Please request a new one.</p>
        </div>
      {/if}

      <div class="space-y-4">
        <FormField label="New Password">
          <Input 
            type="password" 
            bind:value={password} 
            placeholder="••••••••"
            required
            disabled={!token}
          >
            {#snippet icon()}
              <Lock size={18} />
            {/snippet}
          </Input>
        </FormField>

        <FormField label="Confirm New Password">
          <Input 
            type="password" 
            bind:value={confirmPassword} 
            placeholder="••••••••"
            required
            disabled={!token}
          >
            {#snippet icon()}
              <Lock size={18} />
            {/snippet}
          </Input>
        </FormField>
      </div>

      <button
        type="submit"
        disabled={isLoading || !token}
        class="btn-premium w-full py-4 flex items-center justify-center gap-3 group text-lg"
      >
        {#if isLoading}
          <Loader2 class="w-6 h-6 animate-spin text-white" />
        {:else}
          <KeyRound size={20} class="group-hover:rotate-12 transition-transform" />
        {/if}
        Update Access
      </button>

      <div class="pt-6 border-t border-slate-100 dark:border-slate-800 text-center">
        <a href="/login" class="inline-flex items-center gap-2 text-xs font-bold text-slate-500 dark:text-slate-400 hover:text-primary transition-colors">
          Log In instead
        </a>
      </div>
    </form>
  {/if}
</AuthLayout>
