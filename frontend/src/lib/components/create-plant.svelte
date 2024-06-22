<script lang="ts">
	import pb from '$lib/pb';
	import Button from './ui/button.svelte';
	import Dialog from './ui/dialog.svelte';
	import Label from './ui/label.svelte';
	import Input from './ui/input.svelte';

	let { open = $bindable(false) }: { open: boolean } = $props();

	let latin = $state('');
	let name = $state('');
</script>

<Dialog bind:open title="Neue Pflanze">
	<div>
		<Label for="latin">Lateinischer Name</Label>
		<Input id="latin" bind:value={latin} placeholder="Alchemilla" />
	</div>
	<div>
		<Label for="name">Deutscher Name</Label>
		<Input id="name" bind:value={name} placeholder="Frauenmantel" />
	</div>
	{#snippet footer()}
		<Button
			disabled={!latin || !name}
			onclick={() => {
				pb.collection('plants').create({ latin, name });
				open = false;
			}}>Erstellen</Button
		>
	{/snippet}
</Dialog>
