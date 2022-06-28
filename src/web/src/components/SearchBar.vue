<script>
export default {
	name: 'SearchBar',
}
</script>
<template>
	<div id="search-bar">
		<input id="search-bar-input" type="text" placeholder="Search something here..." />
		<span>/</span>
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
