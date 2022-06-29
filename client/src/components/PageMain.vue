<script>
import PageNote from './PageNote.vue'

export default {
	name: 'PageMain',
	components: {
		PageNote
	}
}
</script>

<template>
	<main>
		<PageNote />
	</main>
</template>

<style>
main {
	display: flex;
	flex-direction: row;
	justify-content: center;
	align-items: center;
	margin: 1rem 0;
}
</style>
