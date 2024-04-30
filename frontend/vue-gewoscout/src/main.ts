import { createApp } from 'vue'
import App from './App.vue'
import PrimeVue from 'primevue/config';

import Ripple from 'primevue/ripple';

import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Dropdown from 'primevue/dropdown';
import Menubar from 'primevue/menubar';
import Card from 'primevue/card';
import Avatar from 'primevue/avatar';

import './assets/main.css'
import 'primevue/resources/themes/lara-light-amber/theme.css'
import 'primeicons/primeicons.css'
import 'primeflex/primeflex.css';


const app = createApp(App);
app.use(PrimeVue, { ripple: true });

app.directive('ripple', Ripple);

app.component('vueButton', Button);
app.component('vueInputText', InputText);
app.component('vueDropdown', Dropdown);
app.component('vueMenubar', Menubar);
app.component('vueCard', Card);
app.component('vueAvatar', Avatar);

app.mount('#app');
