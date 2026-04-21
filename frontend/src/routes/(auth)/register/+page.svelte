<script lang="ts">
  import AuthLayout from "$lib/components/layouts/AuthLayout.svelte";
  import FormField from "$lib/components/forms/FormField.svelte";
  import Input from "$lib/components/forms/Input.svelte";
  import Checkbox from "$lib/components/forms/Checkbox.svelte";
  import { Mail, Lock, User, UserPlus, Loader2, ArrowLeft } from "lucide-svelte";
  import { api } from "$lib/api/client.svelte";
  import { authStore } from "$lib/stores/auth.store";
  import { toast } from "svelte-sonner";
  import { goto } from "$app/navigation";

  let name = $state("");
  let email = $state("");
  let password = $state("");
  let confirmPassword = $state("");
  let agreeTerms = $state(false);
  let isLoading = $state(false);

  async function handleRegister(e: SubmitEvent) {
    e.preventDefault();
    if (password !== confirmPassword) {
      toast.error("Passwords do not match");
      return;
    }
    if (!agreeTerms) {
      toast.error("You must agree to the terms");
      return;
    }

    isLoading = true;
    try {
      const response = await api.post<any>("/auth/register", { name, email, password });
      
      // Auto-login with the returned token and user data
      const { user, access_token } = response;
      authStore.login(user, access_token);
      
      toast.success("Account created! Welcome to the colony, " + user.name);
      goto("/dashboard"); 
    } catch (err: any) {
      toast.error(err.message || "Registration failed");
    } finally {
      isLoading = false;
    }
  }
</script>

<AuthLayout>
  <form onsubmit={handleRegister} class="space-y-6">
    <div class="text-center mb-10">
      <h1 class="text-4xl font-black tracking-tight text-slate-900 dark:text-white mb-2">Join Colony.</h1>
      <p class="text-sm font-medium text-slate-500 dark:text-slate-400">Start building faster with the world's most artisanal framework.</p>
    </div>

    <div class="space-y-4">
      <FormField label="Full Name">
        <Input 
          type="text" 
          bind:value={name} 
          placeholder="John Doe"
          required
        >
          {#snippet icon()}
            <User size={18} />
          {/snippet}
        </Input>
      </FormField>

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

      <div class="space-y-4">
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

        <FormField label="Confirm">
          <Input 
            type="password" 
            bind:value={confirmPassword} 
            placeholder="••••••••"
            required
          >
            {#snippet icon()}
              <Lock size={18} />
            {/snippet}
          </Input>
        </FormField>
      </div>
    </div>

    <div class="py-2">
      <Checkbox 
        bind:checked={agreeTerms} 
        label="I agree to the artisanal Terms & Conditions" 
      />
    </div>

    <button
      type="submit"
      disabled={isLoading}
      class="btn-premium w-full py-4 flex items-center justify-center gap-3 group text-lg"
    >
      {#if isLoading}
        <Loader2 class="w-6 h-6 animate-spin text-white" />
      {:else}
        <UserPlus size={20} class="group-hover:scale-110 transition-transform" />
      {/if}
      Create My Account
    </button>

    <div class="pt-6 border-t border-slate-100 dark:border-slate-800 text-center">
      <p class="text-xs font-bold text-slate-500 dark:text-slate-400">
        Already have a colony account? 
        <a href="/login" class="inline-flex items-center gap-1 text-primary hover:underline ml-1">
          <ArrowLeft size={10} />
          Sign In
        </a>
      </p>
    </div>
  </form>
</AuthLayout>
