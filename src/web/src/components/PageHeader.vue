<script>
import SearchBar from './SearchBar.vue'

export default {
	name: 'PageHeader',
	components: {
		SearchBar
	}
}
</script>

<template>
	<header>
		<div id="logo-with-search-bar">
			<img src="./../assets/transparent-logo-light.png" />
			<SearchBar />
		</div>

		<nav>
			<ul>
				<li>
					<a target="_blank" href="https://github.com/luisnquin/nao"
						><img src="./../assets/github-logo.png" />Github</a
					>
				</li>
			</ul>
		</nav>
	</header>
</template>

<style>
header {
	background-color: #4f46e5;
	position: sticky;
	top: 0;
	padding: 1rem;
	display: flex;
	flex-direction: row;
	justify-content: space-between;
	align-items: center;
}

header #logo-with-search-bar {
	display: flex;
	flex-direction: row;
	align-items: center;
	gap: 3rem;
	margin: 0 0 0 10px;
}

header #logo-with-search-bar img {
	max-width: 50px;
}

header nav ul {
	display: flex;
	flex-direction: row;
}

header nav ul li {
	list-style: none;
}

header nav ul li a {
	display: flex;
	flex-direction: row;
	align-items: center;
	gap: 5px;
}

header nav ul li a img {
	width: 1.5rem;
}

header nav ul li a {
	text-decoration: none;
	color: #f3f4f6;
}
</style>
