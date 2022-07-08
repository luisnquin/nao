<script>
import * as indentation from 'indent-textarea'

export default {
	name: 'PageNote',
	props: {
		content: {
			required: true,
		},
		id: {
			required: true
		}
	},
	data() {
		return {
			mutableContent: JSON.parse(JSON.stringify(this.content))
		}
	},
	async mounted() {
		const annotations = document.querySelectorAll('.annotation')
		annotations.forEach((a) => {
			a.setAttribute('style', 'height:' + a.scrollHeight + 'px;overflow-y:hidden;')
			a.oninput = function () {
				this.style.height = 'auto'
				this.style.height = this.scrollHeight + 'px'
			}
		})

		indentation.watch(annotations)
	},
	methods: {
		save(e) {
			if (e.key == 's' && e.ctrlKey) {
				e.preventDefault()

				console.log(this.content)
				console.log(this.newContent)

				fetch('http://localhost:5000/notes/' + this.id, {
					method: 'PATCH',
					body: JSON.stringify({ content: this.content }),
					headers: {
						'Content-Type': 'application/json'
					}
				})
					.then(res => console.log(res))
					.catch(err => console.error(err))
			}
		}
	},
}
</script>

<template>
	<article class="annotation-container">
		<textarea class="annotation" spellcheck="false" placeholder="Write something here..." @keydown="save"
			v-model="mutableContent"></textarea>
	</article>
</template>

<style>
.annotation-container {
	width: 60%;
}

.annotation-container textarea {
	padding: 50px;
	width: 80%;
	font-size: 18px;
	outline: none;
	resize: none;
	tab-size: 4;
	border-radius: 3px;
	border: 0.2px solid #949292;
	box-shadow: 4px 3px rgb(226, 226, 226);
	font-family: 'Cascadia Code', monospace;
}
</style>

<!--
	#3e464f && #fff ?
-->