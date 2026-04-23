<script lang="ts" module>
	import { type FormPathLeaves } from "sveltekit-superforms";
	type T = Record<string, unknown>;
</script>

<script lang="ts">
	import { formFieldProxy, type SuperForm } from "sveltekit-superforms";
	import Input from "../ui/Input.svelte";

	interface Props {
		form: SuperForm<T>;
		name: FormPathLeaves<T>;
		label?: string;
		placeholder?: string;
		type?: string;
		class?: string;
	}

	let { form, name, label, placeholder, type = "text", class: className }: Props = $props();

	const proxy = $derived(formFieldProxy(form, name));
	const value = $derived(proxy.value);
	const errors = $derived(proxy.errors);
	const constraints = $derived(proxy.constraints);
</script>

<Input
	{label}
	{placeholder}
	{type}
	name={name as string}
	bind:value={$value}
	error={$errors?.[0]}
	required={$constraints?.required}
	containerClass={className}
	{...$constraints}
/>
