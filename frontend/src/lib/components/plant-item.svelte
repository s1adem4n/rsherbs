<script lang="ts">
	import Trash from '~icons/lucide/trash';
	import Printer from '~icons/lucide/printer';

	import type { Plant } from '$lib/pb';
	import GetLabel from './get-label.svelte';
	import pb from '$lib/pb';

	let dialogOpen = $state(false);

	let { plant }: { plant: Plant } = $props();
</script>

<GetLabel bind:open={dialogOpen} {plant} />

<div class="py-2 flex gap-2">
	<span>{plant.latin}</span>
	<span class="text-gray-500 italic">{plant.name}</span>
	<button class="hover:underline ml-auto mr-2" onclick={() => (dialogOpen = true)}>
		<Printer />
	</button>
	<button
		class="hover:underline text-red-500"
		onclick={() => pb.collection('plants').delete(plant.id)}><Trash /></button
	>
</div>
