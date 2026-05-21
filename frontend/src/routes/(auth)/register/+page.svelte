<script lang="ts">
  import AuthLayout from "$lib/components/layouts/AuthLayout.svelte";
  import FormField from "$lib/components/forms/FormField.svelte";
  import Input from "$lib/components/forms/Input.svelte";
  import Checkbox from "$lib/components/forms/Checkbox.svelte";
  import PasswordStrengthChecker from "$lib/components/forms/PasswordStrengthChecker.svelte";
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
  let validationErrors = $state<Record<string, string[]>>({});

  // Validate password format
  const isPasswordValid = $derived.by(() => {
    const requirements = {
      minLength: password.length >= 8,
      maxLength: password.length <= 72,
      hasUppercase: /[A-Z]/.test(password),
      hasLowercase: /[a-z]/.test(password),
      hasNumber: /[0-9]/.test(password),
      hasSymbol: /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password),
    };
    return (
      requirements.minLength &&
      requirements.maxLength &&
      requirements.hasUppercase &&
      requirements.hasLowercase &&
      requirements.hasNumber &&
      requirements.hasSymbol
    );
  });

  async function handleRegister(e: SubmitEvent) {
    e.preventDefault();
    validationErrors = {};

    if (password !== confirmPassword) {
      toast.error("Passwords do not match");
      validationErrors.password = ["Passwords do not match"];
      return;
    }
    if (!isPasswordValid) {
      toast.error("Password does not meet requirements");
      return;
    }
    if (!agreeTerms) {
      toast.error("You must agree to the terms");
      return;
    }

    isLoading = true;
    try {
      const response = await api.post<any>("/api/auth/register", { name, email, password });
      
      // Auto-login with the returned token and user data
      const { user, access_token } = response;
      authStore.login(user, access_token);
      
      toast.success("Account created successfully! Welcome, " + user.name);
      goto("/dashboard"); 
    } catch (err: any) {
      toast.error(err.message || "Registration failed");
      if (err.errors && typeof err.errors === 'object') {
        validationErrors = err.errors;
      }
    } finally {
      isLoading = false;
    }
  }
</script>

<AuthLayout>
  <form onsubmit={handleRegister} class="space-y-6">
    <div class="text-center mb-10">
      <h1 class="text-4xl font-black tracking-tight text-slate-900 dark:text-white mb-2">Create Account</h1>
      <p class="text-sm font-medium text-slate-500 dark:text-slate-400">Sign up to get started with your account</p>
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
        {#if validationErrors.name}
          <p class="mt-1 text-xs text-red-600 dark:text-red-400">{validationErrors.name[0]}</p>
        {/if}
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
        {#if validationErrors.email}
          <p class="mt-1 text-xs text-red-600 dark:text-red-400">{validationErrors.email[0]}</p>
        {/if}
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
          <PasswordStrengthChecker {password} />
          {#if validationErrors.password}
            <p class="mt-1 text-xs text-red-600 dark:text-red-400">{validationErrors.password[0]}</p>
          {/if}
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
          {#if password && confirmPassword && password !== confirmPassword}
            <p class="mt-1 text-xs text-red-600 dark:text-red-400">Passwords do not match</p>
          {:else if password && confirmPassword && password === confirmPassword}
            <p class="mt-1 text-xs text-green-600 dark:text-green-400">Passwords match ✓</p>
          {/if}
        </FormField>
      </div>
    </div>

    <div class="py-2">
      <Checkbox
        bind:checked={agreeTerms}
        label="I agree to the Terms & Conditions"
      />
    </div>

    <button
      type="submit"
      disabled={isLoading || !isPasswordValid}
      class="btn-premium w-full py-4 flex items-center justify-center gap-3 group text-lg disabled:opacity-50 disabled:cursor-not-allowed"
    >
      {#if isLoading}
        <Loader2 class="w-6 h-6 animate-spin text-white" />
      {:else}
        <UserPlus size={20} class="group-hover:scale-110 transition-transform" />
      {/if}
      Sign Up
    </button>

    <div class="pt-6 border-t border-slate-100 dark:border-slate-800 text-center">
      <p class="text-xs font-bold text-slate-500 dark:text-slate-400">
        Already have an account?
        <a href="/login" class="inline-flex items-center gap-1 text-primary hover:underline ml-1">
          <ArrowLeft size={10} />
          Sign In
        </a>
      </p>
    </div>
  </form>
</AuthLayout>
