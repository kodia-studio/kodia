<script lang="ts">
  import AuthLayout from "$lib/components/layouts/AuthLayout.svelte";
  import FormField from "$lib/components/forms/FormField.svelte";
  import Input from "$lib/components/forms/Input.svelte";
  import Checkbox from "$lib/components/forms/Checkbox.svelte";
  import { Mail, Lock, User, UserPlus, Loader2 } from "lucide-svelte";
  import { api } from "$lib/api/client";
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
      await api.post("/auth/register", { name, email, password });
      toast.success("Account created successfully! Please sign in.");
      goto("/auth/login");
    } catch (err: any) {
      toast.error(err.message || "Registration failed");
    } finally {
      isLoading = false;
    }
  }
</script>

<AuthLayout>
  <form onsubmit={handleRegister} class="space-y-6">
    <div class="text-center mb-8">
      <h1 class="text-2xl font-bold tracking-tight">Create Account</h1>
      <p class="text-sm text-muted-foreground mt-2">Join Kodia and start building faster</p>
    </div>

    <FormField label="Full Name" required>
      <Input 
        type="text" 
        bind:value={name} 
        placeholder="John Doe"
        required
      >
        {#snippet icon()}
          <User class="w-5 h-5" />
        {/snippet}
      </Input>
    </FormField>

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

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
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

      <FormField label="Confirm" required>
        <Input 
          type="password" 
          bind:value={confirmPassword} 
          placeholder="••••••••"
          required
        >
          {#snippet icon()}
            <Lock class="w-5 h-5" />
          {/snippet}
        </Input>
      </FormField>
    </div>

    <Checkbox 
      bind:checked={agreeTerms} 
      label="I agree to the Terms of Service and Privacy Policy" 
    />

    <button
      type="submit"
      disabled={isLoading}
      class="btn-primary w-full py-3 flex items-center justify-center gap-2 group"
    >
      {#if isLoading}
        <Loader2 class="w-5 h-5 animate-spin" />
      {:else}
        <UserPlus class="w-5 h-5 group-hover:scale-110 transition-transform" />
      {/if}
      Create Account
    </button>

    <div class="text-center mt-6">
      <p class="text-sm text-muted-foreground">
        Already have an account? 
        <a href="/auth/login" class="font-semibold text-primary hover:underline">Sign In</a>
      </p>
    </div>
  </form>
</AuthLayout>
