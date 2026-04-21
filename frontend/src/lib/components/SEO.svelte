<script lang="ts">
    import { page } from '$app/stores';

    interface Props {
        title?: string;
        description?: string;
        keywords?: string;
        image?: string;
        url?: string;
    }

    let { 
        title = 'Kodia Framework - Build World-Class Apps with Go & Svelte 5',
        description = 'The high-performance fullstack framework built for developers who ship. Combine the speed of Go with the reactivity of Svelte 5.',
        keywords = 'Go, Svelte 5, Fullstack, Framework, Golang, Web Development',
        image = '/og-image.png',
        url = $page.url.href
    }: Props = $props();
</script>

<svelte:head>
    <!-- Primary Meta Tags -->
    <title>{title}</title>
    <meta name="title" content={title} />
    <meta name="description" content={description} />
    <meta name="keywords" content={keywords} />

    <!-- Open Graph / Facebook -->
    <meta property="og:type" content="website" />
    <meta property="og:url" content={url} />
    <meta property="og:title" content={title} />
    <meta property="og:description" content={description} />
    <meta property="og:image" content={image} />

    <!-- Twitter -->
    <meta property="twitter:card" content="summary_large_image" />
    <meta property="twitter:url" content={url} />
    <meta property="twitter:title" content={title} />
    <meta property="twitter:description" content={description} />
    <meta property="twitter:image" content={image} />
</svelte:head>
