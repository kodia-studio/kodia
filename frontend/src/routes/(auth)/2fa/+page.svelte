<script lang="ts">
  import AuthLayout from "$lib/components/layouts/AuthLayout.svelte";
  import FormField from "$lib/components/forms/FormField.svelte";
  import Input from "$lib/components/forms/Input.svelte";
  import { ShieldCheck, Loader2, ArrowLeft, Smartphone } from "lucide-svelte";
  import { api } from "$lib/api/client.svelte";
  import { authStore } from "$lib/stores/auth.store";
  import { toast } from "svelte-sonner";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  let code = $state("");
  let isLoading = $state(false);
  let mfaToken = $state("");

  onMount(() => {
    mfaToken = sessionStorage.getItem("mfa_token") || "";
    if (!mfaToken) {
      toast.error("MFA session expired. Please login again.");
      goto("/login"); // FIXED: removed /auth/
    }
  });

  async function handleVerify(e: SubmitEvent) {
    e.preventDefault();
    if (code.length !== 6) {
      toast.error("Please enter a valid 6-digit code");
      return;
    }

    isLoading = true;
    try {
      const response = await api.post<any>("/auth/2fa/login-verify", { 
        mfa_token: mfaToken, 
        code 
      });
      
      const { user, access_token } = response.data;
      authStore.login(user, access_token);
      sessionStorage.removeItem("mfa_token");
      toast.success("Identity verified! Welcome, " + user.name);
      goto("/dashboard");
    } catch (err: any) {
      toast.error(err.message || "Invalid 2FA code");
    } finally {
      isLoading = false;
    }
  }
</script>

<AuthLayout>
  <form onsubmit={handleVerify} class="space-y-6">
    <div class="text-center mb-10">
      <div class="w-16 h-16 rounded-2xl bg-amber-500/10 text-amber-500 mx-auto flex items-center justify-center mb-6 shadow-2xl shadow-amber-500/20 rotate-3">
        <Smartphone size={32} />
      </div>
      <h1 class="text-4xl font-black tracking-tight text-slate-900 dark:text-white mb-2">Shield.</h1>
      <p class="text-sm font-medium text-slate-500 dark:text-slate-400">
        Enter the 6-digit code from your authenticator app to authorize this entry.
      </p>
    </div>

    <FormField label="Verification Code">
      <Input 
        type="text" 
        bind:value={code} 
        placeholder="000000"
        maxlength={6}
        required
        class="text-center text-3xl tracking-[0.5em] font-black font-mono py-10"
      >
        {#snippet icon()}
          <ShieldCheck size={20} class="text-slate-400" />
        {/snippet}
      </Input>
    </FormField>

    <button
      type="submit"
      disabled={isLoading}
      class="btn-premium w-full py-4 flex items-center justify-center gap-3 group text-lg"
    >
      {#if isLoading}
        <Loader2 class="w-6 h-6 animate-spin text-white" />
      {:else}
        <ShieldCheck size={20} class="group-hover:scale-110 transition-transform" />
      {/if}
      Verify Identity
    </button>

    <div class="pt-6 border-t border-slate-100 dark:border-slate-800 text-center">
      <a href="/login" class="inline-flex items-center gap-2 text-xs font-bold text-slate-500 dark:text-slate-400 hover:text-primary transition-colors">
        <ArrowLeft size={10} />
        Back to Login
      </a>
    </div>
  </form>
</AuthLayout>
