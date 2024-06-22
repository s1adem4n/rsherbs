<script lang="ts">
	import '../app.css';
	let { children } = $props();
</script>

<div class="flex flex-col gap-2 w-full max-w-2xl mx-auto p-4">
	{@render children()}
</div>
