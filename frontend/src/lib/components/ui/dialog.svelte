<script lang="ts">
	import X from '~icons/lucide/x';
	import type { Snippet } from 'svelte';

	let {
		open = $bindable(false),
		title,
		children,
		footer,
		header
	}: {
		open?: boolean;
		title?: string;
		children: Snippet;
		footer?: Snippet;
		header?: Snippet;
	} = $props();

	let dialog: HTMLDialogElement;
	$effect(() => {
		if (open) {
			document.body.style.overflow = 'hidden';
			dialog.showModal();
		} else {
			document.body.style.overflow = 'auto';
			dialog.close();
		}
	});
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
<div class="absolute inset-0 border-none {!open && 'pointer-events-none'}">
	<dialog
		bind:this={dialog}
		onclose={() => {
			open = false;
		}}
		class="mx-auto mt-auto w-full max-w-xl rounded-t-xl border border-gray-200 bg-white px-4 pb-4 sm:m-auto sm:rounded-xl"
	>
		<div class="flex max-h-[calc(100vh-4rem)] flex-col">
			<div class="flex bg-white py-4">
				{#if header}
					{@render header()}
				{:else}
					<h2 class="text-2xl font-bold">{title}</h2>
				{/if}
				<button
					class="ml-auto"
					onclick={() => {
						open = false;
					}}
				>
					<X class="h-6 w-6" />
				</button>
			</div>
			<div class="flex flex-col gap-4 overflow-y-auto">
				{@render children()}
			</div>
			{#if footer}
				<div class="flex w-full pt-4 gap-4">
					{@render footer()}
				</div>
			{/if}
		</div>
	</dialog>
</div>
