<script lang="ts" setup>
import { onMounted, ref, watch } from "vue";
import { GetConfig, PutConfig } from "../../wailsjs/go/main/App";
import { config } from "../../wailsjs/go/models";
import SwitchItem from "../components/SwitchItem.vue";
import { copyObjTo, toggleTheme } from "../util";
// 子组件不能修改 props，拷贝一下
const conf = ref(
    new config.Config({
        showGacha: new config.ShowGacha(),
        savedURLs: new Array<config.savedURL>(),
    })
);
onMounted(async () => {
    let c = await GetConfig();
    copyObjTo(c, conf.value);
    watch(conf.value, () => {
        PutConfig(conf.value);
    });
});
</script>
<template>
    <div class="box">
        <span class="title">全局</span>
        <SwitchItem class="item" @change="(val) => toggleTheme(val)" v-model="conf.isDarkTheme">
            深色模式</SwitchItem
        >
        <span class="title">同步</span>
        <SwitchItem class="item" v-model="conf.isAutoSync">自动同步</SwitchItem>
        <SwitchItem class="item" v-model="conf.isUseProxy">使用代理服务器</SwitchItem>

        <span class="title">显示祈愿</span>
        <SwitchItem class="item" v-model="conf.showGacha.g301">角色活动祈愿</SwitchItem>
        <SwitchItem class="item" v-model="conf.showGacha.g302">武器活动祈愿</SwitchItem>
        <SwitchItem class="item" v-model="conf.showGacha.g200">常驻祈愿</SwitchItem>
        <SwitchItem class="item" v-model="conf.showGacha.g100">新手祈愿</SwitchItem>
    </div>
</template>
<style scoped>
.item {
    margin-top: 6px;
}
.title {
    margin-top: 16px;
    font-size: 14px;
    margin-right: auto;
    color: var(--el-text-color-regular);
}
.box {
    height: calc(100% - 60px);
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
