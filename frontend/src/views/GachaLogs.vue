<script lang="ts" setup>
import { onMounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { GetConfig, PutConfig } from "../../wailsjs/go/app/App";
import { GetGachaLogs } from "../../wailsjs/go/app/GachaMan";
import { config, models } from "../../wailsjs/go/models";
import GachaItem from "../components/GachaItem.vue";
import SwitchItem from "../components/SwitchItem.vue";
import { gachaTypeToName } from "../util";
// 为了获取 gachaType，故引入路由，在 onMounted 处使用
const route = useRoute();
const gachaType = ref("");
const isShowFilter = ref(false);
// 排序选项一般是临时开启，故单独拎出来不放在 filterOption 里，不做持久化
const sortDESC = ref(false);
let sortDESCHasChanged = false;
watch(sortDESC, () => {
    sortDESCHasChanged = !sortDESCHasChanged;
});
const filterOption = ref(new config.FilterOption());
const isEmpty = ref(true);
const isLoading = ref(true);
let isEnd = false;
// 按照过滤选项请求祈愿记录
async function getGachaLog() {
    if (isEnd) return;
    isLoading.value = true;
    let c = await GetConfig();
    let logs = await GetGachaLogs(
        page.value,
        c.selectedUid,
        c.gachaLang,
        gachaType.value,
        filterOption.value,
        sortDESC.value
    );
    if (logs.length === 0) {
        isEmpty.value = true;
        isEnd = true;
        isLoading.value = true;
    }
    gachaData.value.push(...logs);
    if (gachaData.value.length !== 0) isEmpty.value = false;
    isLoading.value = false;
    page.value++;
}
const page = ref(0);
const gachaData = ref(new Array<models.FullGachaLog>());
// 保存过滤设置,如果值没有改变返回 false
async function saveFilterOption(): Promise<boolean> {
    let c = await GetConfig();
    let f = filterOption.value;
    let cf = c.filterOption;
    let v = false;
    let hasChange = (v1: boolean, v2: boolean) => {
        if (v) return v;
        v = !(v1 == v2);
        return v;
    };

    hasChange(cf.avatar4, f.avatar4);
    hasChange(cf.avatar5, f.avatar5);
    hasChange(cf.weapon3, f.weapon3);
    hasChange(cf.weapon4, f.weapon4);
    hasChange(cf.weapon5, f.weapon5);
    let hasChanged = v || sortDESCHasChanged;
    if (hasChanged) {
        cf.avatar4 = f.avatar4;
        cf.avatar5 = f.avatar5;
        cf.weapon3 = f.weapon3;
        cf.weapon4 = f.weapon4;
        cf.weapon5 = f.weapon5;
        PutConfig(c);
        sortDESCHasChanged = false;
    }
    return hasChanged;
}

async function handDrawerClose(done: () => void) {
    let hasChanged = await saveFilterOption();
    if (hasChanged) {
        isEnd = false;
        gachaData.value = [];
        page.value = 0;
        await getGachaLog();
    }
    done();
}
onMounted(async () => {
    let c = await GetConfig();
    filterOption.value.avatar4 = c.filterOption.avatar4;
    filterOption.value.avatar5 = c.filterOption.avatar5;
    filterOption.value.weapon3 = c.filterOption.weapon3;
    filterOption.value.weapon4 = c.filterOption.weapon4;
    filterOption.value.weapon5 = c.filterOption.weapon5;
    gachaType.value = route.query.gachaType as string;
    await getGachaLog();
});
</script>
<template>
    <div style="height: 100%; overflow: hidden">
        <div style="position: relative">
            <div class="back" @click="$router.back()">
                <svg
                    width="28"
                    height="28"
                    viewBox="0 0 48 48"
                    fill="none"
                    xmlns="http://www.w3.org/2000/svg"
                >
                    <path
                        d="M31 36L19 24L31 12"
                        stroke="var(--el-text-color-placeholder)"
                        stroke-width="4"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    />
                </svg>
            </div>
            <h2>- {{ gachaTypeToName($route.query.gachaType as string) }} -</h2>
            <div class="filter" @click="isShowFilter = true">
                <svg
                    width="30"
                    height="30"
                    viewBox="0 0 48 48"
                    fill="none"
                    xmlns="http://www.w3.org/2000/svg"
                >
                    <circle cx="12" cy="24" r="3" fill="var(--el-text-color-placeholder)" />
                    <circle cx="24" cy="24" r="3" fill="var(--el-text-color-placeholder)" />
                    <circle cx="36" cy="24" r="3" fill="var(--el-text-color-placeholder)" />
                </svg>
            </div>
            <el-drawer
                v-model="isShowFilter"
                direction="rtl"
                size="40%"
                :before-close="
                    (done:() => void) => handDrawerClose(done)
                "
            >
                <template #header><span style="text-align: left">筛选</span></template>
                <el-scrollbar style="width: 100%; height: 100%">
                    <div class="box">
                        <span class="filter-title-top">角色</span>
                        <!-- 武器活动祈愿 (302) 不会出现五星角色，故不显示此选项 -->
                        <SwitchItem
                            v-if="gachaType != '302'"
                            class="filter-item"
                            v-model="filterOption.avatar5"
                            >五星</SwitchItem
                        >
                        <SwitchItem class="filter-item" v-model="filterOption.avatar4"
                            >四星</SwitchItem
                        >
                        <span class="filter-title">武器</span>
                        <!-- 角色活动祈愿 (301) 不会出现五星武器，故不显示此选项 -->
                        <SwitchItem
                            v-if="gachaType != '301'"
                            class="filter-item"
                            v-model="filterOption.weapon5"
                            >五星</SwitchItem
                        >
                        <SwitchItem class="filter-item" v-model="filterOption.weapon4"
                            >四星</SwitchItem
                        >
                        <SwitchItem class="filter-item" v-model="filterOption.weapon3"
                            >三星</SwitchItem
                        >
                        <span class="filter-title">排序</span>
                        <SwitchItem class="filter-item" v-model="sortDESC">倒序</SwitchItem>
                    </div></el-scrollbar
                >
            </el-drawer>
        </div>
        <el-scrollbar style="width: 100%; height: calc(100% - 74px)">
            <div
                infinite-scroll-distance="100"
                infinite-scroll-immediate="false"
                v-infinite-scroll="getGachaLog"
                :infinite-scroll-disabled="isLoading"
            >
                <div v-for="item in gachaData" style="margin-bottom: 16px; height: 100%">
                    <GachaItem :data="item"></GachaItem>
                </div>
            </div>
        </el-scrollbar>
        <h2 class="empty" v-if="isEmpty">空空如也</h2>
    </div>
</template>
<style scoped>
.empty {
    letter-spacing: 4px;
    position: absolute;
    left: calc(50% - 56px);
    top: 30%;
    color: var(--el-text-color-disabled);
}

.box {
    text-align: left;
    display: flex;
    flex-flow: column;
    margin-right: 12px;
}
.back {
    cursor: pointer;
    position: absolute;
    top: 6px;
    left: 20px;
    width: 50px;
}
.filter-title-top {
    font-size: 14px;
    color: var(--el-text-color-regular);
}
.filter-title {
    margin-top: 16px;
    font-size: 14px;
    color: var(--el-text-color-regular);
}
.filter-item {
    margin-top: 6px;
    margin-bottom: 4px;
}
.filter {
    cursor: pointer;
    position: absolute;
    top: 5px;
    right: 24px;
    width: 50px;
}
.back:hover > svg > path {
    stroke: var(--el-text-color-primary);
}
.filter:hover > svg > circle {
    fill: var(--el-text-color-primary);
}
</style>
