import { createApp } from "vue";
import { WindowCenter, WindowSetSize } from "../wailsjs/runtime/runtime";
import App from "./App.vue";
import { GachaTypeWithName, Option } from "./type";

const app = createApp(App);

app.mount("#app");

export function gachaTypeToName(type: string): string {
    let n = GachaTypeWithName.get(type);
    if (n === undefined) {
        return "";
    }
    return n;
}
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
export async function sleep(timeout: number) {
    return new Promise<void>((resolve) => {
        setTimeout(() => {
            resolve();
        }, timeout);
    });
}
