import { useDark, useToggle } from "@vueuse/core";
import {
    LogDebug,
    LogError,
    LogFatal,
    LogInfo,
    LogPrint,
    LogTrace,
    LogWarning,
} from "../wailsjs/runtime";
import { WindowSetDarkTheme, WindowSetLightTheme } from "../wailsjs/runtime/runtime";
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
const _toggleDark = useToggle(useDark());
export const toggleTheme = (isDark: boolean) => {
    _toggleDark(isDark);
    if (isDark) {
        WindowSetDarkTheme();
    } else {
        WindowSetLightTheme();
    }
};

// 将 src 的值赋值到 dest 上
export function copyObjTo(src: any, dest: any) {
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
