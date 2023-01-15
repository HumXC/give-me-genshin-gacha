<script lang="ts" setup>
import { useDark, useToggle } from "@vueuse/core";
import { Sync } from "../../wailsjs/go/main/App";
import { WindowSetDarkTheme, WindowSetLightTheme } from "../../wailsjs/runtime/runtime";

let sendErr = async () => {
    let data = await Sync(true);
    console.log(JSON.stringify(data));
};

const isDark = useDark();
const toggleDark = useToggle(isDark);

let dark = () => {
    toggleDark();
    if (isDark.value) {
        WindowSetDarkTheme();
    } else {
        WindowSetLightTheme();
    }
};
</script>
<template>
    <div id="right">
        <button @click="sendErr">发送错误</button>
        <button @click="dark">切换深色模式</button>
    </div>
</template>
