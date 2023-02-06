<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { GetConfig } from "../wailsjs/go/app/App";
import iconBook from "./components/icons/book.vue";
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

const isButtonLighting = ref([false, true]);
function lightButton(buttonNum: number) {
    for (let i = 0; i < isButtonLighting.value.length; i++) {
        if (i === buttonNum) continue;
        isButtonLighting.value[i] = false;
    }
}
onMounted(async () => {
    let config = await GetConfig();
    toggleTheme(config.isDarkTheme);
});
</script>

<template>
    <div style="display: flex; height: 100%">
        <!-- 左边功能栏 -->
        <div class="left-bar">
            <div class="mid-line"></div>
            <icon-user
                v-model="isButtonLighting[1]"
                class="button"
                @click="
                    routeTo('/');
                    lightButton(1);
                "
            />
            <icon-book
                v-model="isButtonLighting[2]"
                class="button"
                @click="
                    routeTo('/gacha');
                    lightButton(2);
                "
            />

            <icon-setting
                v-model="isButtonLighting[0]"
                class="button-setting"
                @click="
                    routeTo('/settings');
                    lightButton(0);
                "
            />
        </div>
        <!-- 右边展示页 -->
        <div style="flex: 1" class="right-bar">
            <router-view v-slot="{ Component, route }" v-if="$route.meta.keepAlive">
                <keep-alive>
                    <component :is="Component" :key="route.path" />
                </keep-alive>
            </router-view>
            <router-view v-slot="{ Component, route }" v-if="!$route.meta.keepAlive">
                <component :is="Component" :key="route.path" />
            </router-view>
        </div>
    </div>
</template>
<style scoped>
.mid-line {
    position: absolute;
    min-width: 8px;
    height: 90%;
    left: calc(50% - 4px);
    border-radius: 10px;
    z-index: -1;
    background-color: var(--line);
}
.button:hover {
    cursor: pointer;
}
.button-setting:hover {
    cursor: pointer;
}
.button {
    padding: 1px;
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

.left-bar {
    position: relative;
    margin: 0;
    padding: 0px 5px;
    height: 100%;
    flex: 0 0 10px;
    display: flex;
    flex-wrap: wrap;
    flex-flow: column;
    justify-content: center;
    align-items: center;
}
.right-bar {
    position: relative;
    height: calc(100% - 28px);
    margin-top: 12px;
    margin-right: 12px;
    border-radius: 8px;
    color: var(--el-text-color-primary);
    background-color: var(--bg);
}
</style>
<style></style>
