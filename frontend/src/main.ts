import ElementPlus, { ElMessage, MessageParams } from "element-plus";
import "element-plus/dist/index.css";
import { createApp, watch } from "vue";
import { EventsOn } from "../wailsjs/runtime/runtime";
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

export function setWindow(option: Option) {
    watch(option.showGacha, (data: any) => {
        let count = 0;
        for (const key in data) {
            if (data[key] === true) {
                count++;
            }
        }
    });
}
