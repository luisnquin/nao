import { createApp } from 'vue'
import App from './App.vue'

document.onkeyup = (e) => {
	if (e.key === '/') {
		document.getElementById('search-bar-input').focus()
	}
}

const app = createApp(App)
app.mount('body')
