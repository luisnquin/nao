<script>
export default {
    name: 'Header',

}
</script>

<template>
    <header>
        <nav>
            <ul>
                <li><a target="_blank" href="https://github.com/luisnquin/nao"><img src="./../assets/github-logo.png">Github</a></li>
            </ul>
        </nav>
    </header>
</template>

<style>
header {
    background-color: #4F46E5;
    position: sticky;
    top: 0;
    padding: 1rem;
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
    color: #F3F4F6;
}
</style>