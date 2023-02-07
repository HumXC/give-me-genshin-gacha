import { useDark, useToggle } from "@vueuse/core";
import {
    LogDebug,
    LogError,
    LogFatal,
    LogInfo,
    LogPrint,
    LogTrace,
    LogWarning,
    WindowSetDarkTheme,
    WindowSetLightTheme,
} from "../wailsjs/runtime";
export async function sleep(timeout: number) {
    return new Promise<void>((resolve) => {
        setTimeout(() => {
            resolve();
        }, timeout);
    });
}
export const logger = {
    Debug: (msg: string) => LogDebug("[Webview] " + msg),
    Error: (msg: string) => LogError("[Webview] " + msg),
    Fatal: (msg: string) => LogFatal("[Webview] " + msg),
    Info: (msg: string) => LogInfo("[Webview] " + msg),
    Print: (msg: string) => LogPrint("[Webview] " + msg),
    Trace: (msg: string) => LogTrace("[Webview] " + msg),
    Warning: (msg: string) => LogWarning("[Webview] " + msg),
};

// TODO:
// 关于半透明颜色的 bug
// 在 main.go 中的 wails.run 函数中设置
//     WebviewIsTransparent: true,
//	   WindowIsTranslucent:  true,
//	   BackdropType:         windows.Acrylic,
// 使用 _toggledark 切换颜色
// 会导致前端页面 background-color 的半透明效果丢失
// 需要鼠标重新点击屏幕恢复

const _isDark = useDark();
const _toggleDark = useToggle(_isDark);
export const toggleTheme = async (isDark: boolean) => {
    _toggleDark(isDark);
    if (isDark) {
        WindowSetDarkTheme();
    } else {
        WindowSetLightTheme();
    }
};

// 将 src 的值赋值到 dest 上
export function copyObjTo(src: any, dest: any) {
    if (dest === undefined) {
        dest = src;
        return;
    }
    for (let key in dest) {
        if (!src.hasOwnProperty(key)) {
            continue;
        }
        if (typeof src[key] === "object") {
            copyObjTo(src[key], dest[key]);
            continue;
        } else if (typeof dest[key] === "function") {
            continue;
        }
        dest[key] = src[key];
    }
}

export function gachaTypeToName(gachaType: string): string {
    switch (gachaType) {
        case "301":
            return "角色活动祈愿";

        case "302":
            return "武器活动祈愿";

        case "200":
            return "常驻祈愿";

        case "100":
            return "新手祈愿";

        default:
            return "未知祈愿类型: " + gachaType;
    }
}

export function formatTime(time: any): string {
    return (time as string).substring(0, 19).replace("T", " ");
}
