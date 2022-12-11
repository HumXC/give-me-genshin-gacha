<!-- 
    详细祈愿数据的展示页面
 -->

<script lang="ts" setup>
import { Close } from "@element-plus/icons-vue";
import { onMounted, ref } from "vue";
import { GetLogs } from "../../wailsjs/go/main/App";
import { gachaTypeToName, sleep } from "../main";
import { GachaLog } from "../type";
import GachaItem from "./GachaItem.vue";
let data = ref({ isShow: false, uid: "", gachaType: "" });

const defaultGachaLogs = () => {
    return new Map([
        ["301", { count: 0, logs: new Array<GachaLog>(), page: 0, isLoaded: false }],
        ["302", { count: 0, logs: new Array<GachaLog>(), page: 0, isLoaded: false }],
        ["200", { count: 0, logs: new Array<GachaLog>(), page: 0, isLoaded: false }],
        ["100", { count: 0, logs: new Array<GachaLog>(), page: 0, isLoaded: false }],
    ]);
};
let gachaLogs = defaultGachaLogs();
const title = ref("祈愿");
const checks = ref({
    showArms: true,
    showRole: true,
    showRank3: false,
    showRank4: true,
    showRank5: true,
});
const checkItemIsShow = (log: GachaLog) => {
    if (log.itemType == "武器" && !checks.value.showArms) {
        return false;
    }
    if (log.itemType == "角色" && !checks.value.showRole) {
        return false;
    }
    switch (log.rankType) {
        case "3":
            return checks.value.showRank3;
        case "4":
            return checks.value.showRank4;
        case "5":
            return checks.value.showRank5;
    }
    return true;
};
const nowLogs = ref({ count: 0, logs: new Array<GachaLog>(), page: 0, isLoaded: false });

const load = async () => {
    // 我实在想不出有什么好方法能让他在只选择 5星 时能够按需加载完所有的5星
    while (true) {
        let uid = data.value.uid;
        let gachaType = data.value.gachaType;
        let nl = nowLogs.value;
        if (uid == "") {
            await sleep(200);
            continue;
        }
        if (nl.isLoaded) {
            await sleep(200);
            continue;
        }
        let logs = await GetLogs(uid, gachaType, 100, nowLogs.value.page);

        if (logs.length == 0) {
            nl.isLoaded = true;
            await sleep(200);
            continue;
        }
        console.log("有" + logs.length + "条记录");
        nl.page++;
        nl.logs.push(...logs);
        nl.count += logs.length;
        await sleep(200);
    }
};

const refresh = () => {
    gachaLogs = defaultGachaLogs();
    load();
};
const show = (uid: string, gachaType: string) => {
    title.value = gachaTypeToName(gachaType);
    if (data.value.uid != uid) {
        refresh();
    }
    data.value.uid = uid;
    data.value.gachaType = gachaType;
    console.log("打开祈愿记录" + gachaType);
    let g = gachaLogs.get(gachaType);
    if (g === undefined) {
        console.log("无法匹配 gachaType: " + gachaType);
        return;
    }
    nowLogs.value = g;
    data.value.isShow = true;
};
onMounted(() => {
    load();
});
defineExpose({
    show,
    refresh,
});
</script>
<template>
    <el-drawer
        class="drawer"
        style="padding: 0"
        size="100%"
        :show-close="false"
        v-model="data.isShow"
        direction="ttb"
    >
        <template #title>
            <div class="title">
                <span>{{ title }}</span>
                <el-icon class="close-btn" @click="data.isShow = false">
                    <Close />
                </el-icon>
                <div class="checkbox-box">
                    <el-checkbox
                        class="checkbox"
                        v-model="checks.showRole"
                        label="角色"
                        size="large"
                    />
                    <el-checkbox
                        class="checkbox"
                        v-model="checks.showArms"
                        label="武器"
                        size="large"
                    />
                    <el-checkbox
                        class="checkbox"
                        v-model="checks.showRank3"
                        label="3星"
                        size="large"
                    />
                    <el-checkbox
                        class="checkbox"
                        v-model="checks.showRank4"
                        label="4星"
                        size="large"
                    />
                    <el-checkbox
                        class="checkbox"
                        v-model="checks.showRank5"
                        label="5星"
                        size="large"
                    />
                </div>
            </div>
        </template>
        <!-- TODO: 拖到栏不靠右边，要是能去除 drawer-body的padding就好了 -->
        <el-scrollbar height="100%" v-if="data.isShow">
            <div class="items">
                <div v-for="i in nowLogs.count">
                    <GachaItem
                        v-if="checkItemIsShow(nowLogs.logs[i - 1])"
                        :gacha-log="nowLogs.logs[i - 1]"
                        :uid="data.uid"
                        :id="nowLogs.logs[i - 1].id"
                    />
                </div>
                <i></i> <i></i> <i></i> <i></i> <i></i> <i></i> <i></i> <i></i> <i></i></div
        ></el-scrollbar>
    </el-drawer>
</template>
<style scoped>
.checkbox-box {
    width: 400px;
    position: absolute;
    left: 50%;
    top: 90%;
    margin-left: -190px;
}
.checkbox {
    margin-right: 20px;
}
.items {
    display: flex;
    flex-direction: row;
    justify-content: center;
    flex-wrap: wrap;
}
.items > i {
    width: 320px;
    margin-right: 16px;
    margin-top: 10px;
}
.title {
    position: relative;
    width: 100%;
}
.close-btn {
    position: absolute;
    width: 20px;
    height: 40px;
    top: -10px;
    right: 10px;
}
.close-btn :hover {
    color: var(--el-color-primary);
    cursor: pointer;
}
</style>
