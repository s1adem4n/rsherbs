<script lang="ts">
	import ChevronUp from '~icons/lucide/chevron-up';
	import ChevronDown from '~icons/lucide/chevron-down';

	let {
		value = $bindable(0),
		placeholder = '',
		id = '',
		min,
		max,
		onchange
	}: {
		value?: number;
		placeholder?: string;
		id?: string;
		min?: number;
		max?: number;
		onchange?: (event: Event & { currentTarget: HTMLInputElement }) => void;
	} = $props();
</script>

<div class="relative flex w-full items-center">
	<input
		type="number"
		class="w-full rounded-xl border border-gray-200 px-4 py-2 appearance-none shadow-sm focus-visible:border-blue-500 focus-visible:outline-none"
		{placeholder}
		{id}
		{min}
		{max}
		{onchange}
		bind:value
	/>
	<div class="absolute right-4 flex flex-col">
		<button
			onclick={() => {
				if (max && value < max) {
					value += 1;
				} else if (!max) {
					value += 1;
				}
			}}
		>
			<ChevronUp class="h-4 w-4" />
		</button>
		<button
			onclick={() => {
				if (min && value > min) {
					value -= 1;
				} else if (!min) {
					value -= 1;
				}
			}}
		>
			<ChevronDown class="h-4 w-4" />
		</button>
	</div>
</div>

<style>
	input {
		appearance: textfield;
	}
</style>
