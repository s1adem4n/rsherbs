<script lang="ts">
	import Printer from '~icons/lucide/printer';

	import { BASE_URL, type Plant } from '$lib/pb';
	import Button from './ui/button.svelte';
	import Dialog from './ui/dialog.svelte';
	import Label from './ui/label.svelte';
	import NumberInput from './ui/number-input.svelte';

	let {
		open = $bindable(false),
		plant
	}: {
		plant: Plant;
		open?: boolean;
	} = $props();

	let quantity = $state(1);
	let width = $state(80);
</script>

<Dialog bind:open title="Etiketten">
	<div>
		<Label for="quantity">Anzahl</Label>
		<NumberInput id="quantity" bind:value={quantity} min={1} />
	</div>
	<div>
		<Label for="width">Breite</Label>
		<NumberInput id="width" bind:value={width} min={1} />
	</div>

	{#snippet footer()}
		<Button
			onclick={() => {
				const params = new URLSearchParams();
				params.append('quantity', quantity.toString());
				params.append('width', width.toString());
				params.append('ids', plant.id);

				window.open(`${BASE_URL}/labels?${params.toString()}`, '_blank');
				open = false;
			}}
		>
			Öffnen
		</Button>
		<button
			onclick={() => {
				const params = new URLSearchParams();
				params.append('quantity', quantity.toString());
				params.append('width', width.toString());
				params.append('print', 'true');
				params.append('ids', plant.id);

				window.open(`${BASE_URL}/labels?${params.toString()}`, '_blank');
				open = false;
			}}
		>
			<Printer class="h-8 w-8" />
		</button>
	{/snippet}
</Dialog>
