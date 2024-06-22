<script lang="ts">
	import Printer from '~icons/lucide/printer';

	import { BASE_URL, type Plant } from '$lib/pb';
	import Button from './ui/button.svelte';
	import Checkbox from './ui/checkbox.svelte';
	import Dialog from './ui/dialog.svelte';
	import Label from './ui/label.svelte';
	import NumberInput from './ui/number-input.svelte';

	let {
		open = $bindable(false),
		plants
	}: {
		plants: Plant[];
		open?: boolean;
	} = $props();

	let quantity = $state(1);
	let width = $state(80);
	let selectedPlants: Record<string, boolean> = $state({});

	$effect(() => {
		plants.forEach((plant) => {
			if (selectedPlants[plant.id] === undefined) {
				selectedPlants[plant.id] = true;
			}
		});
	});
</script>

<Dialog bind:open title="Etiketten">
	<div>
		<Label for="width">Breite</Label>
		<NumberInput id="width" bind:value={width} min={1} />
	</div>
	{#each plants as plant}
		<div class="flex gap-2 items-center">
			<Checkbox
				id="qr"
				value={true}
				onchange={(e) => {
					selectedPlants[plant.id] = e.currentTarget.checked;
				}}
			/>
			<Label for="qr">{plant.latin}</Label>
		</div>
	{/each}

	{#snippet footer()}
		<Button
			onclick={() => {
				const params = new URLSearchParams();
				params.append('quantity', quantity.toString());
				params.append('width', width.toString());
				params.append(
					'ids',
					Object.keys(selectedPlants)
						.filter((id) => selectedPlants[id])
						.join(',')
				);

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
				params.append(
					'ids',
					Object.keys(selectedPlants)
						.filter((id) => selectedPlants[id])
						.join(',')
				);

				window.open(`${BASE_URL}/labels?${params.toString()}`, '_blank');
				open = false;
			}}
		>
			<Printer class="h-8 w-8" />
		</button>
	{/snippet}
</Dialog>
