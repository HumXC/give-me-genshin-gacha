import ElementPlus, {
    ElMessage,
    ElNotification,
    MessageParams,
    NotificationHandle,
} from "element-plus";
import "element-plus/dist/index.css";
import {} from "element-plus/es/components/notification";
import { createApp } from "vue";
import { EventsEmit, EventsOn, WindowCenter, WindowSetSize } from "../wailsjs/runtime/runtime";
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
    console.log(message.msg);

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
// 代理服务器
let proxyNotification: NotificationHandle | undefined = undefined;
EventsOn("proxy-started", () => {
    proxyNotification = ElNotification({
        title: "已经开启代理服务器",
        message: "重新在游戏里打开祈愿记录，关闭此通知会关闭代理服务器并取消同步",
        onClose() {
            EventsEmit("stop-proxy");
        },
        duration: 0,
        offset: 50,
    });
});
EventsOn("proxy-stoped", () => {
    proxyNotification?.close();
    proxyNotification = undefined;
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
