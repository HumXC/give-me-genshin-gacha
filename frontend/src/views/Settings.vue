<script lang="ts" setup>
import { onMounted, ref, watch } from "vue";
import { GetConfig, PutConfig } from "../../wailsjs/go/main/App";
import { config } from "../../wailsjs/go/models";
import SwitchItem from "../components/SwitchItem.vue";
import { copyObjTo, toggleTheme } from "../util";
// 子组件不能修改 props，拷贝一下
const conf = ref(new config.Config());
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
        <SwitchItem @change="(val) => toggleTheme(val)" v-model="conf.isDarkTheme"
            >深色模式</SwitchItem
        >
    </div>
</template>
<style scoped>
.title {
    margin-bottom: 8px;
    font-size: 14px;
    margin-right: auto;
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
