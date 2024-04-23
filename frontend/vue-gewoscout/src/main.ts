import { createApp } from 'vue'
import App from './App.vue'
import PrimeVue from 'primevue/config';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Menubar from 'primevue/menubar';
import Card from 'primevue/card';

import './assets/main.css'
import 'primevue/resources/themes/lara-light-amber/theme.css'
import 'primeicons/primeicons.css'
import 'primeflex/primeflex.css';


const app = createApp(App);
app.use(PrimeVue);

app.component('vueButton', Button);
app.component('vueMenubar', Menubar);
app.component('vueCard', Card);
app.component('vueInputText', InputText);

app.mount('#app');
