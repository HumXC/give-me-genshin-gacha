export async function sleep(timeout: number) {
    return new Promise<void>((resolve) => {
        setTimeout(() => {
            resolve();
        }, timeout);
    });
}
import {
    LogDebug,
    LogError,
    LogFatal,
    LogInfo,
    LogPrint,
    LogTrace,
    LogWarning,
} from "../wailsjs/runtime";
export const logger = {
    Debug: (msg: string) => LogDebug("[Webview] " + msg),
    Error: (msg: string) => LogError("[Webview] " + msg),
    Fatal: (msg: string) => LogFatal("[Webview] " + msg),
    Info: (msg: string) => LogInfo("[Webview] " + msg),
    Print: (msg: string) => LogPrint("[Webview] " + msg),
    Trace: (msg: string) => LogTrace("[Webview] " + msg),
    Warning: (msg: string) => LogWarning("[Webview] " + msg),
};
