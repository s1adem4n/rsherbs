<script lang="ts">
	import CreatePlant from '$lib/components/create-plant.svelte';
	import GetLabels from '$lib/components/get-labels.svelte';
	import PlantItem from '$lib/components/plant-item.svelte';
	import Button from '$lib/components/ui/button.svelte';
	import { type Plant } from '$lib/pb';
	import pb from '$lib/pb';
	import { onMount } from 'svelte';

	let plants: Plant[] = $state([]);
	let createDialogOpen = $state(false);
	let labelsDialogOpen = $state(false);

	$effect(() => {
		plants.sort((a, b) => (a.latin > b.latin ? 1 : -1));
	});

	const mount = async () => {
		const res = await pb.collection('plants').getList(1, 100);
		plants = res.items;
		plants = plants.sort((a, b) => (a.latin > b.latin ? 1 : -1));

		const unsubscribe = await pb.collection('plants').subscribe('*', (e) => {
			switch (e.action) {
				case 'create':
					plants = [...plants, e.record].sort((a, b) => (a.latin > b.latin ? 1 : -1));
					break;
				case 'update':
					plants = plants.map((p) => (p.id === e.record.id ? e.record : p));
					break;
				case 'delete':
					plants = plants.filter((p) => p.id !== e.record.id);
					break;
			}
		});

		window.addEventListener('beforeunload', () => {
			unsubscribe();
		});

		return unsubscribe;
	};

	onMount(() => {
		const callback = mount();

		return () => callback.then((cb) => cb());
	});
</script>

<CreatePlant bind:open={createDialogOpen} />
<GetLabels bind:open={labelsDialogOpen} {plants} />

<h1 class="text-2xl font-bold">Pflanzen</h1>
<Button onclick={() => (createDialogOpen = true)}>Neue Pflanze</Button>
<Button onclick={() => (labelsDialogOpen = true)}>Etiketten für mehrere Pflanzen</Button>
<div class="flex flex-col divide-y divide-gray-200">
	{#each plants as plant}
		<PlantItem {plant} />
	{/each}
</div>
