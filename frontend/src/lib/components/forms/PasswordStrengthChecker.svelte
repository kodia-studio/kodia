<script lang="ts">
  import { CheckCircle2, Circle } from "lucide-svelte";

  interface Props {
    password: string;
  }

  let { password }: Props = $props();

  // Real-time requirement checking
  let minLength = $derived(password.length >= 8);
  let maxLength = $derived(password.length <= 72);
  let hasUppercase = $derived(/[A-Z]/.test(password));
  let hasLowercase = $derived(/[a-z]/.test(password));
  let hasNumber = $derived(/[0-9]/.test(password));
  let hasSymbol = $derived(/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password));

  let isValid = $derived(
    minLength && maxLength && hasUppercase && hasLowercase && hasNumber && hasSymbol
  );

  // Make requirementsList derived so it updates reactively
  const requirementsList = $derived.by(() => [
    { id: "minLength", label: "At least 8 characters", met: minLength },
    { id: "hasUppercase", label: "At least one uppercase letter (A-Z)", met: hasUppercase },
    { id: "hasLowercase", label: "At least one lowercase letter (a-z)", met: hasLowercase },
    { id: "hasNumber", label: "At least one number (0-9)", met: hasNumber },
    { id: "hasSymbol", label: "At least one symbol (!@#$%^&*...)", met: hasSymbol },
  ]);
</script>

{#if password}
  <div class="mt-3 space-y-2 rounded-lg bg-slate-50 dark:bg-slate-800/30 p-4 border border-slate-200 dark:border-slate-700">
    <p class="text-xs font-semibold text-slate-600 dark:text-slate-300 mb-3">Password Requirements:</p>
    <div class="space-y-2.5">
      {#each requirementsList as req (req.id)}
        <div class="flex items-center gap-3 transition-all duration-200">
          {#if req.met}
            <div class="flex-shrink-0 w-5 h-5 rounded-full bg-green-500 dark:bg-green-600 flex items-center justify-center">
              <CheckCircle2 size={16} class="text-white" />
            </div>
            <span class="text-sm font-medium text-green-700 dark:text-green-300">{req.label}</span>
          {:else}
            <div class="flex-shrink-0 w-5 h-5 rounded-full bg-slate-300 dark:bg-slate-600 flex items-center justify-center">
              <Circle size={16} class="text-slate-400 dark:text-slate-500" />
            </div>
            <span class="text-sm text-slate-500 dark:text-slate-400">{req.label}</span>
          {/if}
        </div>
      {/each}
    </div>

    {#if isValid}
      <div class="mt-3 pt-3 border-t border-green-200 dark:border-green-900 flex items-center gap-2 bg-green-50 dark:bg-green-900/20 rounded px-2 py-1.5">
        <CheckCircle2 size={16} class="text-green-600 dark:text-green-400" />
        <p class="text-xs font-semibold text-green-700 dark:text-green-300">Password is strong! ✓</p>
      </div>
    {/if}
  </div>
{/if}
