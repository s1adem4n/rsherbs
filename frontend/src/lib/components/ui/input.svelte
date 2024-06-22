<script lang="ts">
	let {
		value = $bindable(''),
		placeholder,
		type = 'text',
		id = '',
		onenter
	}: {
		value?: string;
		type?: string;
		placeholder?: string;
		id?: string;
		onenter?: () => void;
	} = $props();
</script>

<input
	class="w-full rounded-xl border border-gray-200 px-4 py-2 shadow-sm focus-visible:border-blue-500 focus-visible:outline-none"
	{placeholder}
	{id}
	{type}
	onkeydown={(e) => {
		if (e.key === 'Enter') {
			onenter?.();
		}
	}}
	bind:value
/>
