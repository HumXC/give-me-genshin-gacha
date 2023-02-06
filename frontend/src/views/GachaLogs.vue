<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { GetConfig, PutConfig } from "../../wailsjs/go/app/App";
import { config } from "../../wailsjs/go/models";
import SwitchItem from "../components/SwitchItem.vue";
import { gachaTypeToName } from "../util";
// 为了获取 gachaType，故引入路由，在 onMounted 处使用
const route = useRoute();
const gachaType = ref("");
const isShowFilter = ref(false);
// 排序选项一般是临时开启，故单独拎出来不放在 filterOption 里，不做持久化
const sortDESC = ref(false);
const filterOption = ref(new config.FilterOption());
const isEmpty = ref(false);
// 按照过滤选项请求祈愿记录
async function getGachaLog() {}

// 保存过滤设置
async function saveFilterOption() {
    let c = await GetConfig();
    c.filterOption.avatar4 = filterOption.value.avatar4;
    c.filterOption.avatar5 = filterOption.value.avatar5;
    c.filterOption.weapon3 = filterOption.value.weapon3;
    c.filterOption.weapon4 = filterOption.value.weapon4;
    c.filterOption.weapon5 = filterOption.value.weapon5;
    PutConfig(c);
}
onMounted(async () => {
    let c = await GetConfig();
    filterOption.value.avatar4 = c.filterOption.avatar4;
    filterOption.value.avatar5 = c.filterOption.avatar5;
    filterOption.value.weapon3 = c.filterOption.weapon3;
    filterOption.value.weapon4 = c.filterOption.weapon4;
    filterOption.value.weapon5 = c.filterOption.weapon5;
    gachaType.value = route.query.gachaType as string;
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
                    (done: () => void) => {
                       saveFilterOption();
                       getGachaLog()
                       done()
                    }
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
        <el-scrollbar style="width: 100%"> </el-scrollbar>
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
