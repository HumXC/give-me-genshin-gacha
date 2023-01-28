<script lang="ts" setup>
import { onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { GetConfig } from "../wailsjs/go/app/App";
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
            <button class="button-home" @click="routeTo('/')">Home</button>
            <button class="button-settings" @click="routeTo('/settings')">Settings</button>
        </div>
        <!-- 右边展示页 -->
        <div style="flex: 1">
            <el-scrollbar :wrap-class="'scrollbar-wrap'" :view-style="['width:100%']">
                <router-view></router-view>
            </el-scrollbar>
        </div>
    </div>

    <!-- <router-link to="/home">Go to Home</router-link> -->
</template>
<style scoped>
.button-home {
    width: 50px;
    height: 50px;
}
.button-settings {
    width: 50px;
    height: 50px;
}
#leftbar {
    margin: 0;
    padding: 0;
    height: 100%;
    flex: 0 0 50px;
    color: black;
    display: flex;
    flex-wrap: wrap;
    flex-flow: column;
}
</style>
<style>
.scrollbar-wrap {
    /* height: calc(100% - 45px); */
    height: calc(100% - 45px);
    width: calc(100% - 45px);
    margin-top: 10px;
    margin-left: 10px;
    margin-right: 10px;
    display: flex;
    flex-flow: column;
    align-items: center;
    border-radius: 8px;
    padding: 12px 12px 12px 12px;
    color: var(--el-text-color-primary);
    background-color: var(--el-bg-color-page);
}
</style>
