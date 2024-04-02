import { createApp } from 'vue';
import App from './App.vue';
import AppLoading from './components/common/AppLoading.vue';
import 'virtual:uno.css';
import './styles/global.css';
import store from './store';
import router from './router';
import { i18n } from './locales';

// app loading
const appLoading = createApp(AppLoading);
appLoading.mount('#appLoading');

const app = createApp(App);

app.use(store);
app.use(i18n);
app.use(router);

app.mount('#app');
