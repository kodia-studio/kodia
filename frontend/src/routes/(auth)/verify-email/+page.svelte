<script lang="ts">
  import AuthLayout from "$lib/components/layouts/AuthLayout.svelte";
  import { ShieldCheck, ShieldAlert, Loader2, Mail, ArrowRight } from "lucide-svelte";
  import { api } from "$lib/api/client.svelte";
  import { page } from "$app/state";
  import { onMount } from "svelte";

  let status = $state<"loading" | "success" | "error">("loading");
  let message = $state("Verifying your email address...");

  const token = page.url.searchParams.get("token") || "";

  onMount(async () => {
    if (!token) {
      status = "error";
      message = "No verification token found. The link is invalid or missing.";
      return;
    }

    try {
      await api.get(`/auth/verify-email?token=${token}`);
      status = "success";
      message = "Your email has been verified successfully. You can now sign in.";
    } catch (err: any) {
      status = "error";
      message = err.message || "Email verification failed. The link might have expired.";
    }
  });
</script>

<AuthLayout>
  <div class="text-center py-10">
    {#if status === "loading"}
      <div class="w-16 h-16 rounded-2xl bg-primary/10 text-primary mx-auto flex items-center justify-center mb-6 shadow-2xl shadow-primary/20 rotate-3">
        <Loader2 size={32} class="animate-spin" />
      </div>
      <h1 class="text-4xl font-black tracking-tight text-slate-900 dark:text-white mb-4">Verifying Email</h1>
      <p class="text-sm font-medium text-slate-500 dark:text-slate-400">{message}</p>
    {:else if status === "success"}
      <div class="w-16 h-16 rounded-2xl bg-emerald-500/10 text-emerald-500 mx-auto flex items-center justify-center mb-6 shadow-2xl shadow-emerald-500/20 rotate-12">
        <ShieldCheck size={32} />
      </div>
      <h1 class="text-4xl font-black tracking-tight text-slate-900 dark:text-white mb-4">Email Verified</h1>
      <p class="text-sm font-medium text-slate-500 dark:text-slate-400 mb-10 leading-relaxed">{message}</p>

      <a href="/login" class="btn-premium w-full py-4 inline-flex items-center justify-center gap-3 text-lg font-bold">
        Continue to Sign In
        <ArrowRight size={20} />
      </a>
    {:else}
      <div class="w-16 h-16 rounded-2xl bg-rose-500/10 text-rose-500 mx-auto flex items-center justify-center mb-6 shadow-2xl shadow-rose-500/20 -rotate-6">
        <ShieldAlert size={32} />
      </div>
      <h1 class="text-4xl font-black tracking-tight text-slate-900 dark:text-white mb-4">Verification Failed</h1>
      <p class="text-sm font-medium text-slate-500 dark:text-slate-400 mb-10 leading-relaxed">{message}</p>

      <div class="space-y-4">
        <a href="/forgot-password" class="btn-premium w-full py-4 inline-flex items-center justify-center gap-2 font-bold">
          Request New Link
        </a>
        <a href="/login" class="text-xs font-bold text-slate-500 hover:text-primary transition-colors block">
          Back to Sign In
        </a>
      </div>
    {/if}
  </div>
</AuthLayout>
