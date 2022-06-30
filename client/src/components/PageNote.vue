<script>
export default {
	name: 'PageNote',
	data() {
		return {
			previousWasTab: false
		}
	},
	mounted: () => {
		const annotations = document.getElementsByTagName('textarea')
		for (let i = 0; i < annotations.length; i++) {
			annotations[i].setAttribute('style', 'height:' + annotations[i].scrollHeight + 'px;overflow-y:hidden;')
			annotations[i].oninput = function () {
				this.style.height = 'auto'
				this.style.height = this.scrollHeight + 'px'
			}
			annotations[i].onkeydown = function (e) {
				if (e.key == 'Tab') {
					e.preventDefault()

					const start = this.selectionStart

					this.value = this.value.substring(0, start) + '\t' + this.value.substring(this.selectionEnd)
					this.selectionStart = this.selectionEnd = start + 1

					this.previousWasTab = true

					return
				}

				if ((e.key == 'z' || e.key == 'Z') && e.ctrlKey && this.previousWasTab) {
					e.preventDefault()

					const start = this.selectionStart

					this.value =
						this.value.substring(0, start).replace('\t', '') + this.value.substring(this.selectionEnd)

					return
				}

				if (!e.ctrlKey) {
					this.previousWasTab = false
				}
			}
		}
	}
}
</script>

<template>
	<article class="note">
		<textarea spellcheck="false">
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla cursus ullamcorper leo, ac condimentum magna vestibulum vitae. Vestibulum consectetur risus at enim volutpat, id sodales neque iaculis. Aenean nunc enim, cursus at nisi quis, ultricies blandit mauris. Maecenas mollis molestie pulvinar. Aenean bibendum enim at ligula ultricies, non dapibus libero aliquet. Vivamus in massa id libero facilisis rutrum. Sed posuere tortor quis pharetra porttitor. Pellentesque bibendum, diam et pretium molestie, dui sapien pellentesque tellus, sit amet accumsan elit felis in odio. Nunc pharetra quam erat, vel sodales metus facilisis mattis. Proin cursus ipsum sed egestas semper. Nam pharetra, elit ut sodales suscipit, nulla enim laoreet leo, dictum luctus lacus nisl et lectus. Suspendisse vitae enim purus. Sed eu interdum turpis, vel fermentum mauris. Donec faucibus bibendum magna, eu varius mi dapibus ultrices. Pellentesque ac eleifend ipsum, in tincidunt diam.

Praesent id velit sit amet leo aliquet rhoncus. Nulla vestibulum ante nisl, in ultricies ipsum egestas nec. Suspendisse at feugiat ex. Morbi dictum blandit nisl nec iaculis. Quisque nisi nisl, aliquam et tortor blandit, accumsan lobortis ex. Interdum et malesuada fames ac ante ipsum primis in faucibus. Curabitur consequat malesuada nunc quis aliquam. Nunc id mauris sed arcu venenatis semper. Sed ut orci eu eros aliquam porta. In lectus nunc, bibendum scelerisque congue et, bibendum et felis. Vivamus semper suscipit gravida. Nullam at tortor elementum, lacinia quam sit amet, congue erat. Morbi a efficitur nisl. Suspendisse tristique tellus vitae aliquet vulputate.

Quisque egestas elit ligula, in interdum augue congue eget. In sodales dapibus ipsum, in tincidunt eros efficitur vel. Suspendisse ut lacus tempus, egestas tortor in, placerat mauris. Mauris et pharetra justo. In varius magna vitae dapibus euismod. Maecenas dictum, odio nec consectetur auctor, ligula ligula venenatis nulla, tincidunt pellentesque erat mi in urna. In in finibus lectus, et cursus ligula. Aenean vitae lorem vitae massa volutpat accumsan tempus sed turpis. Aenean tellus magna, vulputate at purus a, fermentum tristique orci. Fusce rutrum, massa sed eleifend malesuada, dolor ante elementum lacus, ac aliquet ipsum diam non nunc. Proin turpis sem, semper quis bibendum nec, luctus at turpis. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Integer mollis sapien non molestie dapibus. Nulla tempor aliquam neque, sit amet condimentum eros ornare nec.</textarea>
	</article>
</template>

<style>
.note {
	width: 60%;
}

.note textarea {
	padding: 50px;
	width: 80%;
	font-size: 18px;
	outline: none;
	resize: none;
	tab-size: 4;
	border-radius: 3px;
	box-shadow: 4px 3px rgb(226, 226, 226);
}
</style>
