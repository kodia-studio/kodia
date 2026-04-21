<script lang="ts">
  import AuthLayout from "$lib/components/layouts/AuthLayout.svelte";
  import FormField from "$lib/components/forms/FormField.svelte";
  import Input from "$lib/components/forms/Input.svelte";
  import { Mail, ArrowLeft, Loader2, Send, CheckCircle2 } from "lucide-svelte";
  import { api } from "$lib/api/client.svelte";
  import { toast } from "svelte-sonner";
  
  let email = $state('');
  let isLoading = $state(false);
  let isSubmitted = $state(false);

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    isLoading = true;

    try {
      await api.post("/auth/forgot-password", { email });
      isSubmitted = true;
      toast.success('Reset link sent successfully!');
    } catch (err: any) {
      toast.error(err.message || 'Failed to send reset link');
    } finally {
      isLoading = false;
    }
  }
</script>

<AuthLayout>
  {#if isSubmitted}
    <div class="text-center py-10">
      <div class="w-16 h-16 rounded-2xl bg-emerald-500/10 text-emerald-500 mx-auto flex items-center justify-center mb-6 shadow-2xl shadow-emerald-500/20 rotate-12">
        <Send size={32} />
      </div>
      <h1 class="text-4xl font-black tracking-tight text-slate-900 dark:text-white mb-4">Sent.</h1>
      <p class="text-sm font-medium text-slate-500 dark:text-slate-400 mb-10 leading-relaxed">
        We've dispatched a sting-ready reset link to <br/>
        <span class="font-black text-primary">{email}</span>
      </p>
      
      <a href="/login" class="btn-premium w-full py-4 inline-flex items-center justify-center gap-3 text-lg font-bold">
        <ArrowLeft size={20} />
        Return to Hive
      </a>
    </div>
  {:else}
    <form onsubmit={handleSubmit} class="space-y-6">
      <div class="text-center mb-10">
        <h1 class="text-4xl font-black tracking-tight text-slate-900 dark:text-white mb-2">Reset.</h1>
        <p class="text-sm font-medium text-slate-500 dark:text-slate-400">Enter your email and we'll send a recovery link to your colony.</p>
      </div>

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

      <button
        type="submit"
        disabled={isLoading}
        class="btn-premium w-full py-4 flex items-center justify-center gap-3 group text-lg"
      >
        {#if isLoading}
          <Loader2 class="w-6 h-6 animate-spin text-white" />
        {:else}
          <Send size={20} class="group-hover:translate-x-1 group-hover:-translate-y-1 transition-transform" />
        {/if}
        Request Access
      </button>

      <div class="pt-6 border-t border-slate-100 dark:border-slate-800 text-center">
        <a href="/login" class="inline-flex items-center gap-2 text-xs font-bold text-slate-500 dark:text-slate-400 hover:text-primary transition-colors">
          <ArrowLeft size={10} />
          Back to Sign In
        </a>
      </div>
    </form>
  {/if}
</AuthLayout>
