import { createApp } from 'vue'
import App from './App.vue'

document.onkeyup = (e) => {
	if (e.key === '/' && e.target.tagName != 'TEXTAREA') {
		document.getElementById('search-bar-input').focus()
	}
}

const app = createApp(App)
app.mount('body')
