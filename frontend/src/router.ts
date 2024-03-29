import { createRouter, createWebHashHistory } from "vue-router";
import Gacha from "./views/Gacha.vue";
import GachaLogs from "./views/GachaLogs.vue";
import Home from "./views/Home.vue";
import Settings from "./views/Settings.vue";
export const router = createRouter({
    history: createWebHashHistory(),
    // `params` 不能与 `path` 一起使用
    routes: [
        { path: "/", name: "home", component: Home, meta: { keepAlive: true } },
        { path: "/settings", name: "settings", component: Settings },
        { path: "/gacha", name: "gacha", component: Gacha },
        { path: "/gacha/logs", name: "gacha_logs", component: GachaLogs },
    ],
});
