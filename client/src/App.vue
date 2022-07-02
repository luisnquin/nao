<script setup>
import PageHeader from './components/PageHeader.vue'
import PageMain from './components/PageMain.vue'

</script>

<template>
	<div id="app">
		<PageHeader />
		<PageMain />
	</div>
</template>

<style>
* {
	font-family: Verdana, Geneva, Tahoma, sans-serif;
	padding: 0;
	margin: 0;
}

body {
	background-color: #e9eaeb;
}

@font-face {
	font-family: 'Cascadia Code';
	src: URL('./assets/CascadiaCode.ttf') format('truetype');
}
</style>
