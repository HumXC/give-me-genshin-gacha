<script lang="ts" setup>
import ControlBar from "./components/ControlBar.vue";
import OptionMenu from "./components/OptionMenu.vue";
import { GetMsg } from "../wailsjs/go/main/Bridge";
import { ElMessage, MessageParams } from "element-plus";
import { ref, Ref } from "vue";
import { Option } from "./type";
import GachaInfo from "./components/GachaInfo.vue";
import { fa } from "element-plus/es/locale";
import GachaData from "./components/GachaData.vue";
import { useDark, useToggle } from "@vueuse/core";

type Message = {
    type: string;
    msg: string;
};
// 接收后端消息，也可能是前端发送到后端再被捕获的
var messager = new Promise<void>(async () => {
    while (true) {
        var msg: Message = await GetMsg();
        var show: MessageParams;
        switch (msg.type) {
            case "success":
                show = {
                    duration: 1000,
                    type: msg.type,
                };
                break;
            case "warning":
                show = {
                    showClose: true,
                    duration: 5000,
                    type: msg.type,
                };
                break;
            case "error":
                show = {
                    showClose: true,
                    duration: 0,
                    type: msg.type,
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
        show.message = msg.msg;
        ElMessage(show);
    }
});

// 配置文件
// 这个变量被其他组件所关联，起重要作用
// TODO: 应该从后端获取配置文件
// TODO: 退出 OptionMenu 后自动保存配置文件
const option: Ref<Option> = ref(<Option>{
    isShow: false,
    showGacha: {
        roleUp: true,
        armsUp: true,
        permanent: true,
        start: false,
    },
    otherOption: {
        autoSync: false,
        useProxy: false,
        darkTheme: false,
    },
});

// 控制选项侧栏的开启与关闭
const optionMenuData: Ref<{ isShow: boolean; opt: Option }> = ref(<
    { isShow: boolean; opt: Option }
>{
    isShow: false,
    opt: option.value,
});
// 打开选项侧栏
function openOptionMenu() {
    optionMenuData.value.isShow = true;
}

// 打开祈愿数据页面
// TODO
const gachaDataData: Ref<{
    isShow: boolean;
}> = ref(<{ isShow: boolean }>{
    isShow: false,
});
function openGachaDataPage(name: string, count: number) {
    gachaDataData.value.isShow = true;
}
</script>

<template>
    <OptionMenu :data="optionMenuData" />
    <GachaData :data="gachaDataData" />
    <ControlBar @open-option-menu="openOptionMenu" />
    <GachaInfo :data="option" @open-gacha-data-page="openGachaDataPage" />
</template>
