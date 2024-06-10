import { createApp } from 'vue';
import App from './App.vue';
import PrimeVue from 'primevue/config';

import Ripple from 'primevue/ripple';

import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import InputNumber from 'primevue/inputnumber';
import Dropdown from 'primevue/dropdown';
import InputSwitch from 'primevue/inputswitch';
import MultiSelect from 'primevue/multiselect';
import SelectButton from 'primevue/selectbutton';
import Menubar from 'primevue/menubar';
import Card from 'primevue/card';
import Avatar from 'primevue/avatar';
import Divider from 'primevue/divider';
import Menu from 'primevue/menu';
import Dialog from 'primevue/dialog';

import './assets/main.css';
import 'primevue/resources/themes/lara-light-amber/theme.css';
import 'primeicons/primeicons.css';
import 'primeflex/primeflex.css';

import { createPinia } from 'pinia';

const pinia = createPinia();
const app = createApp(App);
app.use(PrimeVue, { ripple: true });

app.directive('ripple', Ripple);

app.component('vueButton', Button);
app.component('vueInputText', InputText);
app.component('vueInputNumber', InputNumber);
app.component('vueDropdown', Dropdown);
app.component('vueInputSwitch', MultiSelect);
app.component('vueMultiSelect', MultiSelect);
app.component('vueSelectButton', SelectButton);
app.component('vueMenubar', Menubar);
app.component('vueCard', Card);
app.component('vueAvatar', Avatar);
app.component('vueDivider', Divider);
app.component('vueMenu', Menu);
app.component('vueDialog', Dialog);

app.use(pinia);

app.mount('#app');
