<!-- 
    详细祈愿数据的展示页面
 -->

<script lang="ts" setup>
import { Close } from "@element-plus/icons-vue";
import { ref, watch } from "vue";
import { GetLogs } from "../../wailsjs/go/main/App";
import { gachaTypeToName, sleep } from "../main";
import { GachaLog } from "../type";
import GachaItem from "./GachaItem.vue";
let data = ref({ isShow: false, uid: "", gachaType: "", isLoading: false });

const defaultGachaLogs = () => {
    return new Map([
        ["301", { count: 0, logs: new Array<GachaLog>(), page: 0 }],
        ["302", { count: 0, logs: new Array<GachaLog>(), page: 0 }],
        ["200", { count: 0, logs: new Array<GachaLog>(), page: 0 }],
        ["100", { count: 0, logs: new Array<GachaLog>(), page: 0 }],
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
watch(checks.value, () => {
    load(false);
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
const nowLogs = ref({ count: 0, logs: new Array<GachaLog>(), page: 0 });

const load = async (isWait: boolean = true) => {
    console.log("load");

    if (data.value.isLoading) return;
    data.value.isLoading = true;
    let uid = data.value.uid;
    let gachaType = data.value.gachaType;
    if (uid == "") {
        // uid 为空时保留禁用状态，防止无限循环
        // data.value.isLoading = false;
        return;
    }
    console.log(
        "从数据库获取祈愿记录: uid=" +
            uid +
            " gachaType=" +
            gachaType +
            " page=" +
            nowLogs.value.page
    );

    let logs = await GetLogs(uid, gachaType, 10, nowLogs.value.page);
    console.log("有" + logs.length + "条记录");
    if (logs.length == 0) {
        if (isWait) await sleep(3000);
        data.value.isLoading = false;
        return;
    }
    nowLogs.value.page++;
    nowLogs.value.logs.push(...logs);
    nowLogs.value.count += logs.length;
    data.value.isLoading = false;
};

const refresh = () => {
    data.value.isLoading = false;
    gachaLogs = defaultGachaLogs();
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

defineExpose({
    show,
    refresh,
});
</script>
<template>
    <el-drawer size="100%" :show-close="false" v-model="data.isShow" direction="ttb">
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
        <el-scrollbar height="100%" v-if="data.isShow"
            ><div
                class="item"
                v-infinite-scroll="load"
                infinite-scroll-delay="200"
                infinite-scroll-distance="100"
                :infinite-scroll-disabled="data.isLoading"
            >
                <div v-for="i in nowLogs.count">
                    <GachaItem
                        v-if="checkItemIsShow(nowLogs.logs[i - 1])"
                        :gacha-log="nowLogs.logs[i - 1]"
                        :uid="data.uid"
                    />
                </div></div
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
.item {
    height: 70px;
    margin-right: 16px;
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
