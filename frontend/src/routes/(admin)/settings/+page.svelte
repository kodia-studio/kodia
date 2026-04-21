<script lang="ts">
  import AdminLayout from "$lib/components/layouts/AdminLayout.svelte";
  import { Shield, Smartphone, QrCode, ClipboardCheck, Loader2, AlertTriangle, CheckCircle2, User, Key, Save } from "lucide-svelte";
  import { api } from "$lib/api/client.svelte";
  import { authStore } from "$lib/stores/auth.store";
  import { toast } from "svelte-sonner";
  import { cn } from "$lib/utils/styles";

  let activeTab = $state("general");
  let isLoading = $state(false);
  
  // Profile State
  let profileName = $state($authStore.user?.name || "");
  
  // Password State
  let currentPassword = $state("");
  let newPassword = $state("");
  let confirmPassword = $state("");

  // 2FA State
  let setupData = $state<any>(null);
  let verificationCode = $state("");
  let recoveryCodes = $state<string[]>([]);
  let showSetup = $state(false);

  async function updateProfile() {
    if (!profileName) return;
    isLoading = true;
    try {
      const res = await api.patch<any>(`/users/${$authStore.user?.id}`, { name: profileName });
      authStore.update(s => ({ ...s, user: res.data }));
      toast.success("Profile updated");
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
      await api.post("/users/me/change-password", {
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

  /* 2FA Logic */
  async function startSetup() {
    isLoading = true;
    try {
      const response = await api.post<any>("/auth/2fa/enable");
      setupData = response;
      showSetup = true;
    } catch (err: any) {
      toast.error(err.message || "Failed to initiate 2FA setup");
    } finally {
      isLoading = false;
    }
  }

  async function verifySetup() {
    if (verificationCode.length !== 6) return;
    isLoading = true;
    try {
      const response = await api.post<any>("/auth/2fa/verify", { code: verificationCode });
      recoveryCodes = response.data.recovery_codes || [];
      toast.success("2FA Enabled Successfully!");
      if ($authStore.user) {
        authStore.update(s => ({ 
          ...s, 
          user: s.user ? { ...s.user, two_factor_enabled: true } : null 
        }));
      }
    } catch (err: any) {
      toast.error(err.message || "Invalid code");
    } finally {
      isLoading = false;
    }
  }

  async function disable2FA() {
    if (!confirm("Are you sure you want to disable 2FA?")) return;
    isLoading = true;
    try {
      await api.delete("/auth/2fa/disable");
      toast.success("2FA Disabled");
      if ($authStore.user) {
        authStore.update(s => ({ 
          ...s, 
          user: s.user ? { ...s.user, two_factor_enabled: false } : null 
        }));
      }
    } catch (err: any) {
      toast.error(err.message || "Failed to disable 2FA");
    } finally {
      isLoading = false;
    }
  }
</script>

<AdminLayout>
  <div class="max-w-5xl space-y-12 pb-10">
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-6 pb-6 border-b border-slate-200/50 dark:border-white/5">
      <div>
        <h1 class="text-4xl font-black font-heading tracking-tight text-slate-900 dark:text-white leading-none">Console Config.</h1>
        <p class="text-xs font-black uppercase tracking-[0.3em] text-slate-400 mt-3 leading-none">Identity & Encryption Protocols</p>
      </div>
      
      <!-- Quick Status -->
      <div class="flex items-center gap-3 px-4 py-2 bg-emerald-500/5 rounded-xl border border-emerald-500/20 text-[10px] font-black uppercase tracking-widest text-emerald-500">
        <div class="w-1.5 h-1.5 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]"></div>
        System Verified.
      </div>
    </div>

    <!-- Elite Tabs Navigation -->
    <div class="relative flex items-center gap-2 p-1.5 bg-slate-100/50 dark:bg-white/5 border border-slate-200/50 dark:border-white/5 rounded-2xl w-fit">
      {#each [
        { id: "general", label: "Identity", icon: User },
        { id: "security", label: "Security", icon: Shield }
      ] as tab}
        <button
          onclick={() => activeTab = tab.id}
          class={cn(
            "flex items-center gap-3 px-8 py-3 text-[11px] font-black uppercase tracking-widest transition-all duration-300 rounded-xl relative",
            activeTab === tab.id 
              ? "bg-white dark:bg-slate-900 text-primary shadow-lg ring-1 ring-black/5 dark:ring-white/10" 
              : "text-slate-500 hover:text-slate-900 dark:hover:text-white"
          )}
        >
          <tab.icon class="w-4 h-4" />
          {tab.label}
        </button>
      {/each}
    </div>

    <div class="grid grid-cols-1 gap-12">
      {#if activeTab === "general"}
        <div class="space-y-10 animate-in fade-in slide-in-from-bottom-4 duration-500">
          <section class="glass p-10 rounded-[32px] border border-slate-200/50 dark:border-white/5 relative overflow-hidden">
            <div class="absolute top-0 right-0 p-10 opacity-5">
               <User size={120} class="text-primary" />
            </div>
            
            <div class="relative z-10 max-w-xl">
              <h3 class="text-2xl font-black text-slate-900 dark:text-white leading-none mb-2">Profile Intel.</h3>
              <p class="text-sm text-slate-400 font-medium mb-10">Update your operational identity and communication matrix.</p>
              
              <form onsubmit={(e) => { e.preventDefault(); updateProfile(); }} class="space-y-8">
                <div class="space-y-3">
                  <label for="name" class="text-[10px] font-black uppercase tracking-[0.2em] text-slate-500">Clearance Name</label>
                  <input 
                    id="name" 
                    bind:value={profileName} 
                    class="w-full px-6 py-4 bg-slate-50 dark:bg-white/5 border border-slate-200 dark:border-white/10 rounded-2xl text-sm focus:ring-4 focus:ring-primary/10 focus:border-primary/30 outline-none transition-all font-bold" 
                    placeholder="Enter full name" 
                  />
                </div>
                
                <div class="space-y-3 group/readonly">
                  <label for="email" class="text-[10px] font-black uppercase tracking-[0.2em] text-slate-500 group-hover/readonly:text-rose-500 transition-colors">Identity Hash (Email)</label>
                  <div class="relative">
                    <input 
                      id="email" 
                      value={$authStore.user?.email} 
                      class="w-full px-6 py-4 bg-slate-100/50 dark:bg-white/5 border border-slate-200/50 dark:border-white/5 rounded-2xl text-sm text-slate-400 cursor-not-allowed font-mono" 
                      readonly 
                    />
                    <div class="absolute right-6 top-1/2 -translate-y-1/2">
                       <Shield class="w-4 h-4 text-slate-300" />
                    </div>
                  </div>
                </div>

                <div class="pt-4">
                  <button type="submit" class="btn-premium px-10 py-4 text-xs uppercase tracking-widest flex items-center gap-3" disabled={isLoading}>
                    {#if isLoading}<Loader2 class="w-4 h-4 animate-spin" />{:else}<Save class="w-4 h-4" />{/if}
                    Save Matrix.
                  </button>
                </div>
              </form>
            </div>
          </section>
        </div>
      {:else if activeTab === "security"}
        <div class="space-y-12 animate-in fade-in slide-in-from-bottom-4 duration-500">
          <!-- Artisanal 2FA Stage -->
          <section class="relative group">
            <div class="absolute -inset-1 bg-linear-to-r from-primary/20 to-secondary/20 rounded-[40px] blur opacity-50"></div>
            
            <div class="relative glass p-10 rounded-[40px] border border-slate-200/50 dark:border-white/5 bg-white/40 dark:bg-slate-900/60 transition-all duration-500">
              <div class="flex flex-col lg:flex-row lg:items-start justify-between gap-10">
                <div class="flex items-start gap-6">
                  <div class="w-16 h-16 rounded-3xl bg-linear-to-tr from-primary to-secondary flex items-center justify-center text-white shadow-xl shadow-primary/20 ring-4 ring-primary/10 transition-transform group-hover:scale-105 duration-500">
                    <Smartphone class="w-8 h-8 fill-white/20" />
                  </div>
                  <div>
                    <h3 class="text-2xl font-black text-slate-900 dark:text-white leading-none">Security Shield.</h3>
                    <p class="text-sm text-slate-500 dark:text-slate-400 font-medium mt-3 max-w-sm leading-relaxed">
                      Implement Multi-Factor Authentication to solidify your terminal's perimeter and prevent unauthorized access.
                    </p>
                    {#if $authStore.user?.two_factor_enabled}
                      <div class="mt-6 flex items-center gap-2 text-emerald-500 bg-emerald-500/10 px-4 py-2 rounded-full text-[10px] font-black uppercase tracking-widest border border-emerald-500/20 shadow-[0_0_15px_rgba(16,185,129,0.1)] w-fit">
                        <CheckCircle2 class="w-4 h-4 shadow-[0_0_10px_rgba(16,185,129,0.5)]" />
                        Protocol Active.
                      </div>
                    {/if}
                  </div>
                </div>
                
                <div class="shrink-0">
                  {#if $authStore.user?.two_factor_enabled}
                    <button onclick={disable2FA} class="px-8 py-4 rounded-2xl bg-rose-500/10 text-rose-500 border border-rose-500/20 text-xs font-black uppercase tracking-widest hover:bg-rose-500 hover:text-white transition-all active:scale-95" disabled={isLoading}>
                      Suspend Shield.
                    </button>
                  {:else if !showSetup}
                    <button onclick={startSetup} class="btn-premium px-10 py-4 text-xs uppercase tracking-widest" disabled={isLoading}>
                      {isLoading ? 'Processing...' : 'Deploy 2FA.'}
                    </button>
                  {/if}
                </div>
              </div>

              {#if showSetup && !recoveryCodes.length}
                <div class="mt-12 p-10 bg-slate-100/50 dark:bg-black/40 rounded-[32px] border border-slate-200 dark:border-white/5 grid grid-cols-1 lg:grid-cols-2 gap-12 items-center animate-in zoom-in-95 duration-500">
                  <div class="flex flex-col items-center lg:items-start text-center lg:text-left space-y-6">
                    <p class="text-[10px] font-black uppercase tracking-[0.3em] text-primary">Protocol: Scan Identity</p>
                    <div class="relative group/qr p-6 bg-white rounded-3xl shadow-2xl transition-transform hover:scale-105 duration-500">
                      <div class="absolute -inset-2 bg-linear-to-r from-primary to-secondary blur opacity-0 group-hover/qr:opacity-20 transition-opacity"></div>
                      <img src={setupData?.qr_code} alt="2FA QR" class="w-44 h-44 relative z-10" />
                    </div>
                    <div class="w-full max-w-xs p-4 bg-white/5 border border-white/5 rounded-xl font-mono text-[10px] text-slate-500 text-center break-all">
                      {setupData?.secret}
                    </div>
                  </div>
                  
                  <div class="space-y-8">
                    <p class="text-[10px] font-black uppercase tracking-[0.3em] text-primary">Protocol: Finalize Verify</p>
                    <div class="space-y-4">
                      <input 
                        type="text" 
                        maxlength="6"
                        bind:value={verificationCode}
                        placeholder="000000"
                        class="w-full bg-white dark:bg-slate-950 border border-slate-200 dark:border-white/10 rounded-2xl text-center text-4xl font-black font-mono tracking-[0.5em] h-24 focus:ring-8 focus:ring-primary/10 transition-all outline-none"
                      />
                      <button onclick={verifySetup} class="w-full btn-premium h-16 text-xs uppercase tracking-widest" disabled={isLoading || verificationCode.length !== 6}>
                        {isLoading ? 'Verifying...' : 'Establish Connection.'}
                      </button>
                    </div>
                  </div>
                </div>
              {/if}

              {#if recoveryCodes.length}
                <div class="mt-12 p-8 bg-emerald-500/5 border border-emerald-500/20 rounded-[32px] space-y-6 animate-in slide-in-from-top-4">
                  <div class="flex items-center gap-4 text-emerald-500">
                    <ClipboardCheck class="w-8 h-8" />
                    <h4 class="text-xl font-black">Emergency Access Protocols.</h4>
                  </div>
                  <p class="text-sm font-medium text-slate-500 leading-relaxed">Save these recovery identity hashes. They are the final contingency for terminal access should primary verification fail.</p>
                  <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
                    {#each recoveryCodes as code}
                      <code class="bg-white/40 dark:bg-white/5 border border-slate-200 dark:border-white/5 p-3 rounded-xl text-center font-mono text-sm font-black text-primary uppercase shadow-sm">{code}</code>
                    {/each}
                  </div>
                  <button onclick={() => { recoveryCodes = []; showSetup = false; }} class="w-full btn-premium py-4">
                    Secure Matrix Saved.
                  </button>
                </div>
              {/if}
            </div>
          </section>

          <!-- Password Redesign -->
          <section class="glass p-10 rounded-[32px] border border-slate-200/50 dark:border-white/5 relative overflow-hidden">
            <div class="absolute top-0 right-0 p-10 opacity-5">
               <Key size={100} class="text-orange-500" />
            </div>
            
            <div class="relative z-10">
              <div class="flex items-center gap-4 mb-10">
                <div class="w-12 h-12 rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center border border-orange-500/20">
                  <Key class="w-6 h-6" />
                </div>
                <h3 class="text-xl font-black text-slate-900 dark:text-white">Cipher Rotation.</h3>
              </div>
              
              <form onsubmit={(e) => { e.preventDefault(); updatePassword(); }} class="space-y-8 max-w-md">
                <div class="space-y-3">
                  <label for="cur_pass" class="text-[10px] font-black uppercase tracking-[0.2em] text-slate-500">Legacy Cipher</label>
                  <input id="cur_pass" type="password" bind:value={currentPassword} class="w-full px-6 py-4 bg-slate-50 dark:bg-white/5 border border-slate-200 dark:border-white/10 rounded-2xl text-sm focus:ring-4 focus:ring-primary/10 focus:border-primary/30 outline-none transition-all font-bold" required />
                </div>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
                  <div class="space-y-3">
                    <label for="new_pass" class="text-[10px] font-black uppercase tracking-[0.2em] text-slate-500">New Cipher</label>
                    <input id="new_pass" type="password" bind:value={newPassword} class="w-full px-6 py-4 bg-slate-50 dark:bg-white/5 border border-slate-200 dark:border-white/10 rounded-2xl text-sm focus:ring-4 focus:ring-primary/10 focus:border-primary/30 outline-none transition-all font-bold" required />
                  </div>
                  <div class="space-y-3">
                    <label for="conf_pass" class="text-[10px] font-black uppercase tracking-[0.2em] text-slate-500">Confirm Cipher</label>
                    <input id="conf_pass" type="password" bind:value={confirmPassword} class="w-full px-6 py-4 bg-slate-50 dark:bg-white/5 border border-slate-200 dark:border-white/10 rounded-2xl text-sm focus:ring-4 focus:ring-primary/10 focus:border-primary/30 outline-none transition-all font-bold" required />
                  </div>
                </div>
                <div class="pt-4">
                  <button type="submit" class="w-full py-4 rounded-2xl bg-slate-900 dark:bg-white text-white dark:text-slate-950 font-black text-[11px] uppercase tracking-[0.2em] hover:scale-105 active:scale-95 transition-all shadow-xl" disabled={isLoading}>
                    Initiate Rotation.
                  </button>
                </div>
              </form>
            </div>
          </section>
        </div>
      {/if}
    </div>
  </div>
</AdminLayout>
