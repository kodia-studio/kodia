/**
 * Notification Sound System
 * Plays notification sound alerts with configurable settings
 */

export interface NotificationSoundConfig {
    enabled: boolean;
    volume: number; // 0-1
    soundType: 'chime' | 'bell' | 'ping' | 'pop'; // Different sound options
}

const DEFAULT_CONFIG: NotificationSoundConfig = {
    enabled: true,
    volume: 0.5,
    soundType: 'chime'
};

class NotificationSoundManager {
    private config: NotificationSoundConfig = { ...DEFAULT_CONFIG };
    private audioContext: AudioContext | null = null;
    private readonly STORAGE_KEY = 'notification_sound_config';

    constructor() {
        this.loadConfig();
    }

    /**
     * Initialize AudioContext (required for user interaction)
     */
    initAudioContext(): void {
        if (typeof window === 'undefined' || this.audioContext) return;

        try {
            const AudioContextClass = (window as any).AudioContext || (window as any).webkitAudioContext;
            if (AudioContextClass) {
                this.audioContext = new AudioContextClass();
            }
        } catch (err) {
            console.warn('[Notification Sound] Failed to initialize AudioContext:', err);
        }
    }

    /**
     * Play notification sound
     */
    play(notificationType: 'success' | 'error' | 'warning' | 'info' = 'info'): void {
        if (!this.config.enabled) return;
        if (!this.audioContext && typeof window !== 'undefined') {
            this.initAudioContext();
        }
        if (!this.audioContext) return;

        try {
            this.playSound(notificationType);
        } catch (err) {
            console.error('[Notification Sound] Failed to play sound:', err);
        }
    }

    /**
     * Generate and play sound using Web Audio API
     */
    private playSound(notificationType: 'success' | 'error' | 'warning' | 'info'): void {
        if (!this.audioContext) return;

        const ctx = this.audioContext;
        const now = ctx.currentTime;
        const volume = this.config.volume;

        // Create gain node for volume control
        const gainNode = ctx.createGain();
        gainNode.gain.value = volume;
        gainNode.connect(ctx.destination);

        // Different sound patterns based on type
        switch (this.config.soundType) {
            case 'chime':
                this.playChime(ctx, gainNode, now, notificationType);
                break;
            case 'bell':
                this.playBell(ctx, gainNode, now);
                break;
            case 'ping':
                this.playPing(ctx, gainNode, now, notificationType);
                break;
            case 'pop':
                this.playPop(ctx, gainNode, now);
                break;
            default:
                this.playChime(ctx, gainNode, now, notificationType);
        }
    }

    /**
     * Play chime sound (musical)
     */
    private playChime(ctx: AudioContext, gainNode: GainNode, now: number, type: string): void {
        const frequencies = this.getFrequenciesForType(type);

        frequencies.forEach((freq, index) => {
            const osc = ctx.createOscillator();
            osc.type = 'sine';
            osc.frequency.value = freq;

            const envGain = ctx.createGain();
            const delay = index * 0.1; // Stagger notes
            const duration = 0.2;

            envGain.gain.setValueAtTime(0, now + delay);
            envGain.gain.linearRampToValueAtTime(1, now + delay + 0.05);
            envGain.gain.exponentialRampToValueAtTime(0.01, now + delay + duration);

            osc.connect(envGain);
            envGain.connect(gainNode);

            osc.start(now + delay);
            osc.stop(now + delay + duration);
        });
    }

    /**
     * Play bell sound
     */
    private playBell(ctx: AudioContext, gainNode: GainNode, now: number): void {
        const osc = ctx.createOscillator();
        osc.type = 'sine';
        osc.frequency.setValueAtTime(800, now);
        osc.frequency.exponentialRampToValueAtTime(600, now + 0.5);

        const envGain = ctx.createGain();
        envGain.gain.setValueAtTime(1, now);
        envGain.gain.exponentialRampToValueAtTime(0.01, now + 0.8);

        osc.connect(envGain);
        envGain.connect(gainNode);

        osc.start(now);
        osc.stop(now + 0.8);
    }

    /**
     * Play ping sound
     */
    private playPing(ctx: AudioContext, gainNode: GainNode, now: number, type: string): void {
        const freq = type === 'error' ? 400 : type === 'warning' ? 600 : 800;
        const duration = 0.1;

        const osc = ctx.createOscillator();
        osc.type = 'square';
        osc.frequency.value = freq;

        const envGain = ctx.createGain();
        envGain.gain.setValueAtTime(1, now);
        envGain.gain.exponentialRampToValueAtTime(0.01, now + duration);

        osc.connect(envGain);
        envGain.connect(gainNode);

        osc.start(now);
        osc.stop(now + duration);
    }

    /**
     * Play pop sound
     */
    private playPop(ctx: AudioContext, gainNode: GainNode, now: number): void {
        const osc = ctx.createOscillator();
        osc.type = 'sine';
        osc.frequency.setValueAtTime(150, now);
        osc.frequency.exponentialRampToValueAtTime(50, now + 0.1);

        const envGain = ctx.createGain();
        envGain.gain.setValueAtTime(1, now);
        envGain.gain.exponentialRampToValueAtTime(0.01, now + 0.1);

        osc.connect(envGain);
        envGain.connect(gainNode);

        osc.start(now);
        osc.stop(now + 0.1);
    }

    /**
     * Get frequencies based on notification type
     */
    private getFrequenciesForType(type: string): number[] {
        switch (type) {
            case 'success':
                return [523.25, 659.25, 783.99]; // C5, E5, G5
            case 'error':
                return [261.63, 261.63, 261.63]; // C4 repeated
            case 'warning':
                return [440, 440, 550]; // A4, A4, C#5
            case 'info':
            default:
                return [523.25, 523.25, 659.25]; // C5, C5, E5
        }
    }

    /**
     * Update config
     */
    setConfig(config: Partial<NotificationSoundConfig>): void {
        this.config = { ...this.config, ...config };
        this.saveConfig();
    }

    /**
     * Get current config
     */
    getConfig(): NotificationSoundConfig {
        return { ...this.config };
    }

    /**
     * Toggle sound on/off
     */
    toggle(): void {
        this.config.enabled = !this.config.enabled;
        this.saveConfig();
    }

    /**
     * Save config to localStorage
     */
    private saveConfig(): void {
        if (typeof window !== 'undefined') {
            try {
                localStorage.setItem(this.STORAGE_KEY, JSON.stringify(this.config));
            } catch (err) {
                console.error('[Notification Sound] Failed to save config:', err);
            }
        }
    }

    /**
     * Load config from localStorage
     */
    private loadConfig(): void {
        if (typeof window !== 'undefined') {
            try {
                const stored = localStorage.getItem(this.STORAGE_KEY);
                if (stored) {
                    this.config = { ...this.config, ...JSON.parse(stored) };
                }
            } catch (err) {
                console.error('[Notification Sound] Failed to load config:', err);
            }
        }
    }
}

export const notificationSoundManager = new NotificationSoundManager();
