import './assets/css/materialize.min.css'
import './assets/css/icons.css'
import App from './App.svelte'
import { mount } from 'svelte'

const app = mount(App, {
  target: document.getElementById('app')!,
})

export default app
