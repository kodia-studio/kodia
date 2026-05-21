<script lang="ts">
  import { Shield, Smartphone, User, Save, Key, Check, Loader, ClipboardCheck, X, AlertTriangle } from "lucide-svelte";
  import PasswordStrengthChecker from "$lib/components/forms/PasswordStrengthChecker.svelte";
  import { api } from "$lib/api/client.svelte";
  import { authStore } from "$lib/stores/auth.store";
  import { toast } from "svelte-sonner";
  import { cn } from "$lib/utils/styles";
  import { onMount } from "svelte";
  import { browser } from "$app/environment";
  import { goto } from "$app/navigation";

  let activeTab = $state("profile");
  let isLoading = $state(false);
  let copiedCode = $state(false);

  let profileName = $state($authStore.user?.name || "");
  let currentPassword = $state("");
  let newPassword = $state("");
  let confirmPassword = $state("");

  let setupData = $state<any>(null);
  let verificationCode = $state("");
  let recoveryCodes = $state<string[]>([]);
  let showSetup = $state(false);

  let showPasswordModal = $state(false);
  let pendingAction = $state<"enable" | "disable" | "delete" | null>(null);
  let modalPassword = $state("");
  let modalPasswordError = $state("");

  onMount(async () => {
    try {
      const user = await api.get(`/api/users/me`);
      authStore.update(s => ({ ...s, user }));
      saveUserToLocalStorage(user);
    } catch (err) {
      console.error('Failed to refresh user data:', err);
    }
  });

  async function updateProfile() {
    if (!profileName) return;
    isLoading = true;
    try {
      const res = await api.patch<any>(`/api/users/${$authStore.user?.id}`, { name: profileName });
      authStore.update(s => ({ ...s, user: res }));
      saveUserToLocalStorage(res);
      toast.success("Profile updated successfully");
    } catch (err: any) {
      toast.error(err.message || "Update failed");
    } finally {
      isLoading = false;
    }
  }

  async function updatePassword() {
    if (newPassword !== confirmPassword) {
      toast.error("Passwords do not match");
      return;
    }
    isLoading = true;
    try {
      await api.post("/api/users/me/change-password", {
        current_password: currentPassword,
        new_password: newPassword
      });
      toast.success("Password changed successfully");
      currentPassword = "";
      newPassword = "";
      confirmPassword = "";
    } catch (err: any) {
      toast.error(err.message || "Failed to change password");
    } finally {
      isLoading = false;
    }
  }

  function saveUserToLocalStorage(user: any) {
    if (browser && user) {
      localStorage.setItem('user', JSON.stringify(user));
    }
  }

  function openPasswordModal(action: "enable" | "disable") {
    pendingAction = action;
    modalPassword = "";
    modalPasswordError = "";
    showPasswordModal = true;
  }

  async function confirmPasswordAction() {
    if (!modalPassword) {
      modalPasswordError = "Password is required";
      return;
    }

    isLoading = true;
    modalPasswordError = "";

    try {
      if (pendingAction === "enable") {
        const response = await api.post<any>("/api/auth/2fa/enable", { password: modalPassword });
        setupData = response;
        showSetup = true;
      } else if (pendingAction === "disable") {
        await api.delete("/api/auth/2fa/disable", { password: modalPassword });
        toast.success("Two-Factor Authentication disabled");
        if ($authStore.user) {
          const updatedUser = { ...($authStore.user), two_factor_enabled: false };
          authStore.update(s => ({
            ...s,
            user: updatedUser
          }));
          saveUserToLocalStorage(updatedUser);
        }
      } else if (pendingAction === "delete") {
        await api.delete("/api/users/me");
        toast.success("Account deleted successfully");
        authStore.logout();
        goto("/");
        return;
      }
      showPasswordModal = false;
      pendingAction = null;
    } catch (err: any) {
      const errorMsg = err.message || err?.error?.message || "";
      if (errorMsg.toLowerCase().includes("invalid password") || errorMsg.toLowerCase().includes("invalid credentials")) {
        modalPasswordError = "Invalid password";
      } else {
        toast.error(errorMsg || `Failed to ${pendingAction} 2FA`);
      }
    } finally {
      isLoading = false;
    }
  }

  async function verifySetup() {
    if (verificationCode.length !== 6) return;
    isLoading = true;
    try {
      const response = await api.post<any>("/api/auth/2fa/verify", { code: verificationCode });
      recoveryCodes = response.recovery_codes || [];
      toast.success("Two-Factor Authentication enabled successfully!");
      if ($authStore.user) {
        const updatedUser = { ...($authStore.user), two_factor_enabled: true };
        authStore.update(s => ({
          ...s,
          user: updatedUser
        }));
        saveUserToLocalStorage(updatedUser);
      }
    } catch (err: any) {
      toast.error(err.message || "Invalid code");
    } finally {
      isLoading = false;
    }
  }
</script>

<svelte:head>
  <title>Settings | Kodia Console</title>
</svelte:head>

<div class="max-w-5xl space-y-8 md:space-y-12 pb-10">
    <div class="flex flex-col gap-4 pb-4 md:pb-6 border-b border-slate-200/50 dark:border-white/5">
      <div>
        <h1 class="text-2xl md:text-4xl font-black font-heading tracking-tight text-slate-900 dark:text-white leading-none">Settings</h1>
        <p class="text-xs font-medium text-slate-400 mt-2 md:mt-3 leading-none">Manage your account and security preferences</p>
      </div>

      <!-- Quick Status -->
      <div class="flex items-center gap-2 md:gap-3 px-3 md:px-4 py-2 bg-emerald-500/5 rounded-xl border border-emerald-500/20 text-[9px] md:text-[10px] font-medium uppercase tracking-wide md:tracking-widest text-emerald-500 w-fit">
        <div class="w-1.5 h-1.5 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]"></div>
        Account Active
      </div>
    </div>

    <!-- Tabs Navigation -->
    <div class="relative flex items-center gap-1 md:gap-2 p-1 md:p-1.5 bg-slate-100/50 dark:bg-white/5 border border-slate-200/50 dark:border-white/5 rounded-xl md:rounded-2xl overflow-x-auto">
      {#each [
        { id: "profile", label: "Profile", icon: User },
        { id: "security", label: "Security", icon: Shield },
        { id: "danger", label: "Danger Zone", icon: Shield, danger: true }
      ] as tab}
        <button
          onclick={() => activeTab = tab.id}
          class={cn(
            "flex items-center gap-2 px-3 md:px-8 py-2 md:py-3 text-xs md:text-sm font-semibold transition-all duration-300 rounded-lg md:rounded-xl relative whitespace-nowrap",
            activeTab === tab.id
              ? tab.danger
                ? "bg-red-50 dark:bg-red-950/30 text-red-600 dark:text-red-400 shadow-lg ring-1 ring-red-200 dark:ring-red-900/50"
                : "bg-white dark:bg-slate-900 text-primary shadow-lg ring-1 ring-black/5 dark:ring-white/10"
              : tab.danger
                ? "text-red-500/60 hover:text-red-600 dark:hover:text-red-400"
                : "text-slate-500 hover:text-slate-900 dark:hover:text-white"
          )}
        >
          <tab.icon class="w-4 h-4 shrink-0" />
          <span class="hidden sm:inline">{tab.label}</span>
        </button>
      {/each}
    </div>

    <div class="grid grid-cols-1 gap-12">
      {#if activeTab === "profile"}
        <div class="space-y-10 animate-in fade-in slide-in-from-bottom-4 duration-500">
          <section class="glass p-6 md:p-10 rounded-2xl md:rounded-4xl border border-slate-200/50 dark:border-white/5 relative overflow-hidden">
            <div class="absolute top-0 right-0 p-6 md:p-10 opacity-5">
               <User size={80} class="md:w-30 md:h-30 text-primary" />
            </div>

            <div class="relative z-10 max-w-xl">
              <h3 class="text-xl md:text-2xl font-black text-slate-900 dark:text-white leading-none mb-2">Profile Information</h3>
              <p class="text-xs md:text-sm text-slate-400 font-medium mb-6 md:mb-10">Update your personal information and profile details.</p>
              
              <form onsubmit={(e) => { e.preventDefault(); updateProfile(); }} class="space-y-6 md:space-y-8">
                <div class="space-y-2 md:space-y-3">
                  <label for="name" class="text-xs md:text-sm font-semibold text-slate-600 dark:text-slate-300">Full Name</label>
                  <input
                    id="name"
                    bind:value={profileName}
                    class="w-full px-4 md:px-6 py-3 md:py-4 bg-slate-50 dark:bg-white/5 border border-slate-200 dark:border-white/10 rounded-lg md:rounded-2xl text-xs md:text-sm focus:ring-4 focus:ring-primary/10 focus:border-primary/30 outline-none transition-all font-medium"
                    placeholder="Enter your full name"
                  />
                </div>

                <div class="space-y-2 md:space-y-3">
                  <label for="email" class="text-xs md:text-sm font-semibold text-slate-600 dark:text-slate-300">Email Address</label>
                  <div class="relative">
                    <input
                      id="email"
                      value={$authStore.user?.email}
                      class="w-full px-4 md:px-6 py-3 md:py-4 bg-slate-100/50 dark:bg-white/5 border border-slate-200/50 dark:border-white/5 rounded-lg md:rounded-2xl text-xs md:text-sm text-slate-500 cursor-not-allowed font-medium"
                      readonly
                    />
                    <div class="absolute right-4 md:right-6 top-1/2 -translate-y-1/2">
                       <Shield class="w-4 h-4 text-slate-300" />
                    </div>
                  </div>
                  <p class="text-xs text-slate-400">Your email address cannot be changed</p>
                </div>

                <div class="pt-2 md:pt-4">
                  <button type="submit" class="btn-premium px-6 md:px-10 py-3 md:py-4 text-xs md:text-sm font-semibold flex items-center justify-center md:justify-start gap-2 md:gap-3 w-full md:w-auto" disabled={isLoading}>
                    {#if isLoading}<Loader class="w-4 h-4 animate-spin" />{:else}<Save class="w-4 h-4" />{/if}
                    <span>Save Changes</span>
                  </button>
                </div>
              </form>
            </div>
          </section>
        </div>
      {:else if activeTab === "security"}
        <div class="space-y-12 animate-in fade-in slide-in-from-bottom-4 duration-500">
          <!-- Two-Factor Authentication -->
          <section class="relative group">
            <div class="absolute -inset-1 bg-linear-to-r from-primary/20 to-secondary/20 rounded-[40px] blur opacity-50"></div>
            
            <div class="relative glass p-6 md:p-10 rounded-2xl md:rounded-[40px] border border-slate-200/50 dark:border-white/5 bg-white/40 dark:bg-slate-900/60 transition-all duration-500">
              <div class="flex flex-col gap-6 md:gap-10">
                <div class="flex flex-col sm:flex-row sm:items-start gap-4 md:gap-6">
                  <div class="w-12 h-12 md:w-16 md:h-16 rounded-2xl md:rounded-3xl bg-linear-to-tr from-primary to-secondary flex items-center justify-center text-white shadow-xl shadow-primary/20 ring-4 ring-primary/10 transition-transform group-hover:scale-105 duration-500 shrink-0">
                    <Smartphone class="w-6 h-6 md:w-8 md:h-8 fill-white/20" />
                  </div>
                  <div class="flex-1">
                    <h3 class="text-lg md:text-2xl font-black text-slate-900 dark:text-white leading-none">Two-Factor Authentication</h3>
                    <p class="text-xs md:text-sm text-slate-500 dark:text-slate-400 font-medium mt-2 md:mt-3 max-w-sm leading-relaxed">
                      Add an extra layer of security to your account by enabling two-factor authentication. You'll need to enter a code from your phone in addition to your password when logging in.
                    </p>
                    {#if $authStore.user?.two_factor_enabled}
                      <div class="mt-3 md:mt-6 flex items-center gap-2 text-emerald-500 bg-emerald-500/10 px-3 md:px-4 py-2 rounded-full text-xs font-semibold border border-emerald-500/20 shadow-[0_0_15px_rgba(16,185,129,0.1)] w-fit">
                        <Check class="w-4 h-4 shadow-[0_0_10px_rgba(16,185,129,0.5)]" />
                        Enabled
                      </div>
                    {/if}
                  </div>
                </div>

                <div class="flex flex-col sm:flex-row gap-3 sm:justify-end">
                  {#if $authStore.user?.two_factor_enabled}
                    <button onclick={() => openPasswordModal("disable")} class="px-6 md:px-8 py-3 md:py-4 rounded-lg md:rounded-2xl bg-rose-500/10 text-rose-500 border border-rose-500/20 text-xs md:text-sm font-semibold hover:bg-rose-500 hover:text-white transition-all active:scale-95" disabled={isLoading}>
                      Disable 2FA
                    </button>
                  {:else if !showSetup}
                    <button onclick={() => openPasswordModal("enable")} class="btn-premium px-6 md:px-10 py-3 md:py-4 text-xs md:text-sm font-semibold" disabled={isLoading}>
                      {isLoading ? 'Processing...' : 'Enable 2FA'}
                    </button>
                  {/if}
                </div>
              </div>

              {#if showSetup && !recoveryCodes.length}
                <div class="mt-6 md:mt-12 p-6 md:p-10 bg-slate-100/50 dark:bg-black/40 rounded-2xl md:rounded-4xl border border-slate-200 dark:border-white/5 grid grid-cols-1 md:grid-cols-2 gap-6 md:gap-12 items-center animate-in zoom-in-95 duration-500">
                  <div class="flex flex-col items-center md:items-start text-center md:text-left space-y-4 md:space-y-6">
                    <p class="text-xs font-semibold text-primary uppercase tracking-wide">Step 1: Scan QR Code</p>
                    <div class="relative group/qr p-4 md:p-6 bg-white rounded-2xl md:rounded-3xl shadow-2xl transition-transform hover:scale-105 duration-500">
                      <div class="absolute -inset-2 bg-linear-to-r from-primary to-secondary blur opacity-0 group-hover/qr:opacity-20 transition-opacity"></div>
                      <img src={setupData?.qr_code} alt="2FA QR Code" class="w-32 h-32 md:w-44 md:h-44 relative z-10" />
                    </div>
                    <div class="w-full max-w-xs p-3 md:p-4 bg-white/5 border border-white/5 rounded-xl font-mono text-xs text-slate-500 text-center break-all">
                      {setupData?.secret}
                    </div>
                  </div>

                  <div class="space-y-4 md:space-y-8">
                    <p class="text-xs font-semibold text-primary uppercase tracking-wide">Step 2: Verify Code</p>
                    <div class="space-y-3 md:space-y-4">
                      <input
                        type="text"
                        maxlength="6"
                        bind:value={verificationCode}
                        placeholder="000000"
                        class="w-full bg-white dark:bg-slate-950 border border-slate-200 dark:border-white/10 rounded-lg md:rounded-2xl text-center text-2xl md:text-4xl font-black font-mono tracking-[0.3em] md:tracking-[0.5em] h-16 md:h-24 focus:ring-8 focus:ring-primary/10 transition-all outline-none"
                      />
                      <button onclick={verifySetup} class="w-full btn-premium h-12 md:h-16 text-xs md:text-sm font-semibold" disabled={isLoading || verificationCode.length !== 6}>
                        {isLoading ? 'Verifying...' : 'Verify'}
                      </button>
                    </div>
                  </div>
                </div>
              {/if}

              {#if recoveryCodes.length}
                <div class="mt-6 md:mt-12 p-6 md:p-8 bg-emerald-500/5 border border-emerald-500/20 rounded-2xl md:rounded-4xl space-y-4 md:space-y-6 animate-in slide-in-from-top-4">
                  <div class="flex items-center gap-3 md:gap-4 text-emerald-500">
                    <ClipboardCheck class="w-6 h-6 md:w-8 md:h-8 shrink-0" />
                    <h4 class="text-lg md:text-xl font-black">Recovery Codes</h4>
                  </div>
                  <p class="text-xs md:text-sm font-medium text-slate-500 leading-relaxed">Save these recovery codes in a safe place. If you lose access to your authenticator app, you can use these codes to regain access to your account.</p>
                  <div class="grid grid-cols-2 sm:grid-cols-4 gap-2 md:gap-4">
                    {#each recoveryCodes as code}
                      <code class="bg-white/40 dark:bg-white/5 border border-slate-200 dark:border-white/5 p-2 md:p-3 rounded-lg md:rounded-xl text-center font-mono text-xs md:text-sm font-bold text-primary uppercase shadow-sm">{code}</code>
                    {/each}
                  </div>
                  <button onclick={() => { recoveryCodes = []; showSetup = false; }} class="w-full btn-primary py-3 md:py-4 text-xs md:text-sm font-semibold">
                    I've Saved My Recovery Codes
                  </button>
                </div>
              {/if}
            </div>
          </section>

          <!-- Change Password -->
          <section class="glass p-6 md:p-10 rounded-2xl md:rounded-4xl border border-slate-200/50 dark:border-white/5 relative overflow-hidden">
            <div class="absolute top-0 right-0 p-6 md:p-10 opacity-5">
               <Key size={100} class="text-orange-500" />
            </div>

            <div class="relative z-10">
              <div class="flex items-center gap-3 md:gap-4 mb-6 md:mb-10">
                <div class="w-10 h-10 md:w-12 md:h-12 rounded-lg md:rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center border border-orange-500/20 shrink-0">
                  <Key class="w-5 h-5 md:w-6 md:h-6" />
                </div>
                <h3 class="text-lg md:text-xl font-black text-slate-900 dark:text-white">Change Password</h3>
              </div>

              <form onsubmit={(e) => { e.preventDefault(); updatePassword(); }} class="space-y-6 md:space-y-8 max-w-md">
                <div class="space-y-2 md:space-y-3">
                  <label for="cur_pass" class="text-xs md:text-sm font-semibold text-slate-600 dark:text-slate-300">Current Password</label>
                  <input id="cur_pass" type="password" bind:value={currentPassword} class="w-full px-4 md:px-6 py-3 md:py-4 bg-slate-50 dark:bg-white/5 border border-slate-200 dark:border-white/10 rounded-lg md:rounded-2xl text-xs md:text-sm focus:ring-4 focus:ring-primary/10 focus:border-primary/30 outline-none transition-all font-medium" required />
                </div>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 md:gap-6">
                  <div class="space-y-2 md:space-y-3">
                    <label for="new_pass" class="text-xs md:text-sm font-semibold text-slate-600 dark:text-slate-300">New Password</label>
                    <input id="new_pass" type="password" bind:value={newPassword} class="w-full px-4 md:px-6 py-3 md:py-4 bg-slate-50 dark:bg-white/5 border border-slate-200 dark:border-white/10 rounded-lg md:rounded-2xl text-xs md:text-sm focus:ring-4 focus:ring-primary/10 focus:border-primary/30 outline-none transition-all font-medium" required />
                  </div>
                  <div class="space-y-2 md:space-y-3">
                    <label for="conf_pass" class="text-xs md:text-sm font-semibold text-slate-600 dark:text-slate-300">Confirm Password</label>
                    <input id="conf_pass" type="password" bind:value={confirmPassword} class="w-full px-4 md:px-6 py-3 md:py-4 bg-slate-50 dark:bg-white/5 border border-slate-200 dark:border-white/10 rounded-lg md:rounded-2xl text-xs md:text-sm focus:ring-4 focus:ring-primary/10 focus:border-primary/30 outline-none transition-all font-medium" required />
                  </div>
                </div>
                <PasswordStrengthChecker password={newPassword} />
                <div class="pt-2 md:pt-4">
                  <button type="submit" class="w-full py-3 md:py-4 rounded-lg md:rounded-2xl bg-slate-900 dark:bg-white text-white dark:text-slate-950 font-black text-xs md:text-sm uppercase tracking-wide hover:scale-105 active:scale-95 transition-all shadow-xl" disabled={isLoading}>
                    Update Password
                  </button>
                </div>
              </form>
            </div>
          </section>
        </div>
      {:else if activeTab === "danger"}
        <div class="space-y-12 animate-in fade-in slide-in-from-bottom-4 duration-500">
          <!-- Delete Account -->
          <section class="relative group">
            <div class="absolute -inset-1 bg-linear-to-r from-red-500/20 to-orange-500/20 rounded-[40px] blur opacity-50"></div>

            <div class="relative glass p-6 md:p-10 rounded-2xl md:rounded-[40px] border border-red-200/50 dark:border-red-900/30 bg-red-50/40 dark:bg-red-950/20 transition-all duration-500">
              <div class="flex flex-col md:flex-row md:items-start justify-between gap-6 md:gap-10">
                <div class="flex items-start gap-4 md:gap-6">
                  <div class="w-12 h-12 md:w-16 md:h-16 rounded-2xl md:rounded-3xl bg-red-500/20 flex items-center justify-center text-red-600 dark:text-red-400 shadow-xl ring-4 ring-red-500/10 shrink-0">
                    <AlertTriangle class="w-6 h-6 md:w-8 md:h-8" />
                  </div>
                  <div>
                    <h3 class="text-lg md:text-2xl font-black text-red-600 dark:text-red-400 leading-none">Delete Account</h3>
                    <p class="text-xs md:text-sm text-red-600/80 dark:text-red-400/70 font-medium mt-2 md:mt-3 max-w-sm leading-relaxed">
                      Permanently delete your account and all associated data. This action cannot be undone.
                    </p>
                  </div>
                </div>

                <div class="shrink-0 w-full md:w-auto">
                  <button
                    onclick={() => openPasswordModal("delete")}
                    class="w-full md:w-auto px-6 md:px-8 py-3 md:py-4 rounded-lg md:rounded-2xl bg-red-500/10 text-red-600 dark:text-red-400 border border-red-500/20 text-xs md:text-sm font-semibold hover:bg-red-500 hover:text-white transition-all active:scale-95"
                    disabled={isLoading}
                  >
                    Delete My Account
                  </button>
                </div>
              </div>
            </div>
          </section>
        </div>
      {/if}
    </div>
  </div>

<!-- Password Confirmation Modal -->
{#if showPasswordModal}
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
    <div class="bg-white dark:bg-slate-900 rounded-2xl md:rounded-3xl shadow-2xl w-full max-w-md max-h-[90vh] overflow-y-auto flex flex-col animate-in zoom-in-95 duration-300">
      <!-- Sticky Header -->
      <div class="sticky top-0 flex items-center justify-between p-4 md:p-8 bg-white dark:bg-slate-900 border-b border-slate-200/50 dark:border-white/5">
        <h3 class="text-lg md:text-2xl font-black text-slate-900 dark:text-white">Confirm Password</h3>
        <button
          onclick={() => {
            showPasswordModal = false;
            pendingAction = null;
          }}
          class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 transition-colors shrink-0"
        >
          <X class="w-5 h-5 md:w-6 md:h-6" />
        </button>
      </div>

      <!-- Scrollable Content -->
      <div class="flex-1 overflow-y-auto">
        <div class="p-4 md:p-8 space-y-6 md:space-y-8">
          <div class="text-xs md:text-sm text-slate-500 dark:text-slate-400 space-y-2 md:space-y-3">
            {#if pendingAction === "enable"}
              <p>Enter your password to enable two-factor authentication.</p>
            {:else if pendingAction === "disable"}
              <p>Enter your password to disable two-factor authentication.</p>
            {:else if pendingAction === "delete"}
              <p class="text-red-600 dark:text-red-400 font-semibold">Enter your password to permanently delete your account.</p>
              <p class="text-xs text-red-600/80 dark:text-red-400/70">This action cannot be undone.</p>
            {/if}
          </div>

          <div class="space-y-3 md:space-y-4">
            <div class="space-y-2 md:space-y-3">
              <label for="modal-password" class="text-xs md:text-sm font-semibold text-slate-600 dark:text-slate-300">Password</label>
              <input
                id="modal-password"
                type="password"
                bind:value={modalPassword}
                placeholder="Enter your password"
                class="w-full px-4 md:px-6 py-3 md:py-4 bg-slate-50 dark:bg-white/5 border border-slate-200 dark:border-white/10 rounded-lg md:rounded-xl text-xs md:text-sm focus:ring-4 focus:ring-primary/10 focus:border-primary/30 outline-none transition-all font-medium"
                onkeydown={(e) => {
                  if (e.key === "Enter") confirmPasswordAction();
                }}
              />
            </div>
            {#if modalPasswordError}
              <p class="text-xs text-red-600 dark:text-red-400 font-medium">{modalPasswordError}</p>
            {/if}
          </div>
        </div>
      </div>

      <!-- Sticky Footer -->
      <div class="sticky bottom-0 flex gap-2 md:gap-3 p-4 md:p-8 bg-white dark:bg-slate-900 border-t border-slate-200/50 dark:border-white/5">
        <button
          onclick={() => {
            showPasswordModal = false;
            pendingAction = null;
          }}
          class="flex-1 px-4 md:px-6 py-3 md:py-4 rounded-lg md:rounded-xl bg-slate-100 dark:bg-white/10 text-slate-700 dark:text-slate-300 text-xs md:text-sm font-semibold hover:bg-slate-200 dark:hover:bg-white/20 transition-all"
          disabled={isLoading}
        >
          Cancel
        </button>
        <button
          onclick={confirmPasswordAction}
          class="flex-1 px-4 md:px-6 py-3 md:py-4 rounded-lg md:rounded-xl bg-primary text-white text-xs md:text-sm font-semibold hover:scale-105 active:scale-95 transition-all disabled:opacity-50 flex items-center justify-center gap-2"
          disabled={isLoading || !modalPassword}
        >
          {#if isLoading}
            <Loader class="w-4 h-4 animate-spin" />
            <span>Verifying...</span>
          {:else}
            Confirm
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}
