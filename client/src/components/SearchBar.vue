<script>
export default {
	name: 'SearchBar',
	data: () => {
		return {
			sets: null
		}
	},
	async created() {
		const res = await fetch('http://localhost:5000/sets', { method: 'GET' })
		const sets = await res.json()
		this.sets = sets.data
	}
}
</script>
<template>
	<div id="search-bar">
		<input id="search-bar-input" type="search" list="search-bar-options" placeholder="Search something here..."
			autocomplete="off" />
		<span>/</span>

		<datalist id="search-bar-options">
			<option v-for="set in sets" v-bind:key="set">{{ set.tag }}</option>
			<option v-for="set in sets" v-bind:key="set">{{ set.key }}</option>
		</datalist>
	</div>
</template>

<style>
#search-bar {
	display: flex;
	align-items: center;
	justify-content: center;
}

#search-bar #search-bar-input {
	padding: 1px 9px;
	border-radius: 5px;
	border: none;
	outline: none;
	font-size: 13px;
	line-height: 2rem;
	width: 18rem;
}

#search-bar span {
	position: relative;
	right: 20px;
	padding: 2px 6px;
	background-color: #414141;
	border-radius: 5px;
	color: aliceblue;
	pointer-events: none;
}
</style>
