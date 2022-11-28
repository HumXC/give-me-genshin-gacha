import { createApp } from "vue";
import App from "./App.vue";
import ElementPlus from "element-plus";
import "element-plus/dist/index.css";
import { GachaTypeWithName } from "./type";

const app = createApp(App);

app.use(ElementPlus);
app.mount("#app");
export function gachaTypeToName(type: string): string {
    let n = GachaTypeWithName.get(type);
    if (n === undefined) {
        return "";
    }
    return n;
}
