<script>
import PageNote from './PageNote.vue'
import axios from 'axios'

export default {
	name: 'PageMain',
	data() {
		return {
			data: [{}]
		}
	},
	components: {
		PageNote
	},
	async mounted() {
		axios.get('http://localhost:5000/sets')
			.then(res => this.data = res.data.data)
			.catch(err => console.error(err))
	}
	
}
</script>

<template>
	<main>
		<PageNote :content="data[0].content" />
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
