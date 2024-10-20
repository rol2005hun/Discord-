<script>
	import Halloween from './Halloween.svelte';
	import Toasts from '$lib/components/Toasts.svelte';
	import '../app.scss';
</script>

<div class="app">
	<main>
		<slot />
		<Halloween />
		<Toasts />
	</main>
</div>

<style>
</style>
