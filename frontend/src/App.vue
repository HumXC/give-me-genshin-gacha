<script lang="ts" setup>
import { ElMessage, ElMessageBox } from "element-plus";
import { onMounted, ref, Ref } from "vue";
import { GetOption, SaveOption, Sync } from "../wailsjs/go/main/App";
import { main } from "../wailsjs/go/models";
import ControlBar from "./components/ControlBar.vue";
import GachaData from "./components/GachaData.vue";
import GachaInfo from "./components/GachaInfo.vue";
import OptionMenu from "./components/OptionMenu.vue";
import { setWindowSize } from "./main";
import { Option } from "./type";

// 配置文件
// 这个变量被其他组件所关联，起重要作用
// 需要提供一个默认的对象, 不过并不是配置的默认值。默认值在 go 代码中配置
const option: Ref<Option> = ref(new Option());

// 控制选项侧栏的开启与关闭
const optionMenuStatus = ref({
    isShow: false,
});
// 控制祈愿数据页面
const gachaDataData: Ref<{
    isShow: boolean;
}> = ref(<{ isShow: boolean }>{
    isShow: false,
});

// 刷新各组件
const refreshGachaInfo = ref(async () => {});
const refreshControlBar = ref(async () => {});

const refresh = async () => {
    await refreshControlBar.value();
    await refreshGachaInfo.value();
};
// 打卡详细数据页面
function showGachaDataPage(gachaType: string) {
    gachaDataData.value.isShow = true;
}

// 同步数据
async function startSync(done: () => void, useProxy: boolean = false) {
    console.log("开始同步 useProxy: " + useProxy);
    if (!useProxy) {
        useProxy = option.value.otherOption.useProxy;
        console.log("默认使用代理覆盖参数 useProxy: " + useProxy);
    }
    let result = await Sync(useProxy);
    switch (result) {
        case "authkey timeout":
            ElMessageBox.confirm("要使用代理方式重新同步吗?", "链接已经过期", {
                confirmButtonText: "好",
                cancelButtonText: "别",
            })
                .catch(() => {
                    done();
                    ElMessage.error({
                        message: "同步失败",
                        duration: 0,
                        showClose: true,
                    });
                })
                .then(() => {
                    startSync(done, true);
                });
            break;

        case "cancel":
            done();
            break;
        case "fail":
            done();
            break;
        default:
            done();
            ElMessage.success("同步完成 - " + result);
            break;
    }
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
function changeSelectedUid(uid: string) {
    console.log("切换了id" + uid);
    option.value.controlBar.selectedUid = uid;
    saveOption();
}

onMounted(async () => {
    console.log("初始化");
    let o = await GetOption();

    console.log(o.controlBar.selectedUid);

    option.value.otherOption = o.otherOption;
    option.value.showGacha = o.showGacha;
    option.value.controlBar = o.controlBar;
    console.log("成功获取配置");
    refresh();
    setWindowSize(o);
});
</script>

<template>
    <el-scrollbar>
        <OptionMenu
            :option="option"
            :status="optionMenuStatus"
            @closing="(done) => saveOption(done)" />
        <GachaData :data="gachaDataData" />
        <ControlBar
            :option="option"
            @option-button-click="optionMenuStatus.isShow = true"
            @sync-button-click="(done) => startSync(done)"
            @select-uid="(uid) => changeSelectedUid(uid)"
            ref="refreshControlBar" />
        <GachaInfo :option="option" @pie-click="showGachaDataPage" ref="refreshGachaInfo" />
        <div style="height: 10px"
    /></el-scrollbar>
</template>
