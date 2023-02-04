<script lang="ts" setup>
import { onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { GetConfig } from "../wailsjs/go/app/App";
import iconSetting from "./components/icons/setting.vue";
import iconUser from "./components/icons/user.vue";

import { toggleTheme } from "./util";
const route = useRoute();
const router = useRouter();
const routeTo = (path: string) => {
    router.push({
        path: path,
        query: { ...route.query },
    });
};

onMounted(async () => {
    let config = await GetConfig();
    toggleTheme(config.isDarkTheme);
});
</script>

<template>
    <div style="display: flex; height: 100%">
        <!-- 左边功能栏 -->
        <div id="leftbar">
            <div class="mid-line"></div>
            <icon-user class="button" @click="routeTo('/')"></icon-user>
            <icon-setting class="button-setting" @click="routeTo('/settings')">Set</icon-setting>
        </div>
        <!-- 右边展示页 -->
        <div style="flex: 1">
            <el-scrollbar
                :wrap-class="'scrollbar-wrap'"
                :view-style="['width:100%', 'height:100%', 'display:flex', 'flex-flow:column']"
            >
                <router-view v-slot="{ Component }">
                    <keep-alive>
                        <component :is="Component" />
                    </keep-alive>
                </router-view>
            </el-scrollbar>
        </div>
    </div>

    <!-- <router-link to="/home">Go to Home</router-link> -->
</template>
<style scoped>
.mid-line {
    position: absolute;
    min-width: 8px;
    height: 90%;
    background-color: var(--el-fill-color);
    left: calc(50% - 4px);
    border-radius: 10px;
    z-index: -1;
}
.button {
    width: 36px;
    height: 36px;
    margin: 10px 0 10px 0;
    border: 2px solid var(--el-border-color);
    border-radius: 100%;
    background-color: var(--el-fill-color-lighter);
}
.button-setting {
    position: absolute;
    width: 30px;
    height: 30px;
    margin: 10px 0 10px 0;
    border-style: dashed;
    border: 2px solid var(--el-border-color);
    border-radius: 100%;
    bottom: 0px;
    background-color: var(--el-fill-color-lighter);
}

#leftbar {
    position: relative;
    margin: 0;
    padding: 0px 8px;
    height: 100%;
    flex: 0 0 10px;
    display: flex;
    flex-wrap: wrap;
    flex-flow: column;
    justify-content: center;
    align-items: center;
}
</style>
<style>
.scrollbar-wrap {
    /* height: calc(100% - 45px); */
    height: calc(100% - 45px);
    margin-top: 12px;
    margin-right: 12px;
    display: flex;
    flex-flow: column;
    align-items: center;
    border-radius: 8px;
    padding: 12px 12px 12px 12px;
    color: var(--el-text-color-primary);
    background-color: var(--el-bg-color-page);
}
</style>
