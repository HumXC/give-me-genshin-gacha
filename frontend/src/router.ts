import { createRouter, createWebHashHistory } from "vue-router";
import Home from "./views/Home.vue";
import Settings from "./views/Settings.vue";
export const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        { path: "/home", component: Home },
        { path: "/settings", component: Settings },
    ],
});
