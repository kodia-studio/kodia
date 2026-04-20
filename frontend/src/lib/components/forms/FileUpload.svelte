<script lang="ts">
  import { cn } from "$lib/utils/styles";
  import { Upload, X, File as FileIcon, FileText, Image as ImageIcon, Music, Video } from "lucide-svelte";
  import { fade, slide } from "svelte/transition";

  interface Props {
    accept?: string;
    multiple?: boolean;
    maxSize?: number; // In MB
    files?: File[];
    onFilesChange?: (files: File[]) => void;
    error?: string;
    class?: string;
  }

  let { 
    accept = "*",
    multiple = false,
    maxSize = 10,
    files = $bindable([]),
    onFilesChange,
    error,
    class: className 
  }: Props = $props();

  let isDragging = $state(false);
  let fileInput: HTMLInputElement;

  function handleFiles(newFiles: FileList | null) {
    if (!newFiles) return;
    
    const validFiles = Array.from(newFiles).filter(f => f.size <= maxSize * 1024 * 1024);
    
    if (multiple) {
      files = [...files, ...validFiles];
    } else {
      files = validFiles.slice(0, 1);
    }
    
    onFilesChange?.(files);
  }

  function removeFile(index: number) {
    files = files.filter((_, i) => i !== index);
    onFilesChange?.(files);
  }

  function getFileIcon(type: string) {
    if (type.startsWith("image/")) return ImageIcon;
    if (type.startsWith("video/")) return Video;
    if (type.startsWith("audio/")) return Music;
    if (type.includes("pdf") || type.includes("text")) return FileText;
    return FileIcon;
  }
</script>

<div class={cn("w-full space-y-4", className)}>
  <div
    class={cn(
      "relative border-2 border-dashed rounded-2xl flex flex-col items-center justify-center p-10 transition-all duration-300 group",
      isDragging 
        ? "border-primary bg-primary/5 scale-[1.02] shadow-xl shadow-primary/10" 
        : "border-muted-foreground/20 bg-slate-50/50 dark:bg-slate-900/50 hover:bg-muted/10 hover:border-primary/50",
      error && "border-destructive/50 bg-destructive/5"
    )}
    ondragover={(e) => { e.preventDefault(); isDragging = true; }}
    ondragleave={() => isDragging = false}
    ondrop={(e) => { e.preventDefault(); isDragging = false; handleFiles(e.dataTransfer?.files || null); }}
    onclick={() => fileInput.click()}
    role="button"
    tabindex="0"
    onkeydown={(e) => e.key === "Enter" && fileInput.click()}
  >
    <input
      bind:this={fileInput}
      type="file"
      class="hidden"
      {multiple}
      {accept}
      onchange={(e) => handleFiles(e.currentTarget.files)}
    />

    <div class="w-16 h-16 rounded-2xl bg-background shadow-xl flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
      <Upload class={cn("w-8 h-8", isDragging ? "text-primary" : "text-muted-foreground")} />
    </div>

    <div class="text-center">
      <p class="text-lg font-semibold tracking-tight">
        {isDragging ? "Drop your files here" : "Click or drag to upload"}
      </p>
      <p class="text-sm text-muted-foreground mt-1">
        Maximum file size: {maxSize}MB
      </p>
    </div>
  </div>

  {#if files.length > 0}
    <div class="space-y-2" transition:slide>
      {#each files as file, i}
        {@const Icon = getFileIcon(file.type)}
        <div 
          class="flex items-center gap-4 p-3 bg-card border rounded-xl animate-in fade-in slide-in-from-left-2 duration-300"
          transition:fade
        >
          <div class="w-10 h-10 rounded-lg bg-muted flex items-center justify-center shrink-0">
            <Icon class="w-5 h-5 text-muted-foreground" />
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium truncate">{file.name}</p>
            <p class="text-xs text-muted-foreground">{(file.size / 1024 / 1024).toFixed(2)} MB</p>
          </div>
          <button 
            onclick={(e) => { e.stopPropagation(); removeFile(i); }}
            class="p-2 hover:bg-destructive/10 hover:text-destructive rounded-lg transition-colors"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
      {/each}
    </div>
  {/if}
</div>
