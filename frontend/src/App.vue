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
import { main } from "../wailsjs/go/models";
import { GetOption, SaveOption } from "../wailsjs/go/main/App";

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
const gachaDataData: Ref<{
    isShow: boolean;
}> = ref(<{ isShow: boolean }>{
    isShow: false,
});
function openGachaDataPage(gachaType: string, showType: string) {
    gachaDataData.value.isShow = true;
}
// 同步数据
function startSync() {
    console.log("开始同步");
}
// 保存配置
function saveOption(done: () => void) {
    console.log("保存配置");
    SaveOption(main.Option.createFrom(option.value)).then(() => {
        done();
    });
}
// 切换 uid
function changeSelect(uid: string) {
    console.log("切换了id");
}

async function init() {
    console.log("初始化");
    let o = await GetOption();
    option.value.otherOption = o.otherOption;
    option.value.showGacha = o.showGacha;
    console.log("成功获取配置");
}
init();
</script>

<template>
    <OptionMenu :data="optionMenuData" @save-option="(done) => saveOption(done)" />
    <GachaData :data="gachaDataData" />
    <ControlBar
        @open-option-menu="openOptionMenu"
        @start-sync="startSync"
        @change-select="(uid) => changeSelect(uid)"
    />
    <GachaInfo :data="option" @pie-click="openGachaDataPage" />
</template>
