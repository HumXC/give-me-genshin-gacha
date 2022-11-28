<script lang="ts" setup>
import ControlBar from "./components/ControlBar.vue";
import OptionMenu from "./components/OptionMenu.vue";
import { GetMsg, PutMsg } from "../wailsjs/go/main/Bridge";

import { ElMessage, MessageParams } from "element-plus";
import { onMounted, ref, Ref } from "vue";
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
new Promise<void>(async () => {
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
// 需要提供一个默认的对象
const option: Ref<Option> = ref(<Option>{
    showGacha: {
        roleUp: false,
        armsUp: false,
        permanent: false,
        start: false,
    },
    otherOption: {
        autoSync: false,
        useProxy: false,
        darkTheme: false,
    },
    controlBar: {
        selectedUid: "",
    },
});
// 各组件的刷新函数
const controlBarRefresh = ref(() => {});

// 控制选项侧栏的开启与关闭
const optionMenuData: Ref<{ isShow: boolean; opt: Option }> = ref(<
    { isShow: boolean; opt: Option }
>{
    isShow: false,
    opt: option.value,
});
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
function startSync(done: () => void) {
    console.log("开始同步");
    PutMsg({
        type: "message",
        msg: "正在同步哦",
    });
    setTimeout(() => {
        done();
        PutMsg({
            type: "success",
            msg: "同步成功",
        });
        controlBarRefresh.value();
    }, 1000);
}
// 保存配置
function saveOption(done: (() => void) | void) {
    console.log("保存配置");
    SaveOption(main.Option.createFrom(option.value)).then(() => {
        if (done !== undefined) {
            done();
        }
    });
}
// 切换 uid
function changeSelect(uid: string) {
    console.log("切换了id" + uid);
    option.value.controlBar.selectedUid = uid;
    saveOption();
}

async function reFresh() {
    console.log("刷新主界面");
}
onMounted(async () => {
    console.log("初始化");
    let o = await GetOption();

    console.log(o.controlBar.selectedUid);

    option.value.otherOption = o.otherOption;
    option.value.showGacha = o.showGacha;
    option.value.controlBar = o.controlBar;
    console.log("成功获取配置");
});
</script>

<template>
    <OptionMenu :data="optionMenuData" @save-option="(done) => saveOption(done)" />
    <GachaData :data="gachaDataData" />
    <ControlBar
        :option="option"
        @open-option-menu="optionMenuData.isShow = true"
        @start-sync="(done) => startSync(done)"
        @change-select="(uid) => changeSelect(uid)"
        ref="controlBarRefresh"
    />
    <GachaInfo :data="option" @pie-click="openGachaDataPage" />
</template>
