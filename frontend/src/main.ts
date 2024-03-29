import "element-plus/dist/index.css";
import "element-plus/theme-chalk/dark/css-vars.css";

import ElementPlus from "element-plus";
import { createApp } from "vue";
import App from "./App.vue";
import Listen from "./events";
import { router } from "./router";

const app = createApp(App);
app.use(ElementPlus);
app.use(router);
app.mount("#app");

// 监听事件
Listen();
