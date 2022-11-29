import ElementPlus, { ElMessage, MessageParams } from "element-plus";
import "element-plus/dist/index.css";
import { createApp } from "vue";
import { EventsOn, WindowCenter, WindowSetSize } from "../wailsjs/runtime/runtime";
import App from "./App.vue";
import { GachaTypeWithName, Message, Option } from "./type";

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
EventsOn("alert", (message: Message) => {
    var show: MessageParams;
    switch (message.type) {
        case "success":
            show = {
                duration: 1000,
                type: message.type,
            };
            break;
        case "warning":
            show = {
                showClose: true,
                duration: 5000,
                type: message.type,
            };
            break;
        case "error":
            show = {
                showClose: true,
                duration: 0,
                type: message.type,
            };
            break;
        default:
            show = {
                showClose: true,
                duration: 2000,
                type: "info",
            };
            break;
    }
    show.message = message.msg;
    ElMessage(show);
});

export async function setWindowSize(option: Option) {
    // 设置窗口大小
    let count = 0;
    for (const key in option.showGacha) {
        if ((<any>option.showGacha)[key] === true) {
            count++;
        }
    }
    let w = 0;
    let h = 518;
    switch (count) {
        case 1:
            w = 388;
            break;

        case 2:
        case 4:
            w = 770;
            break;

        case 3:
            w = 1120;
            break;
    }
    await WindowSetSize(w, h);
    WindowCenter();
}
