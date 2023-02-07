<script lang="ts" setup>
import { onMounted, ref, watch } from "vue";
import { GetConfig, PutConfig } from "../../wailsjs/go/app/App";
import { config } from "../../wailsjs/go/models";
import SwitchItem from "../components/SwitchItem.vue";
import { copyObjTo, toggleTheme } from "../util";
// 子组件不能修改 props，拷贝一下
const conf = ref(
    new config.Config({
        showGacha: new config.ShowGacha(),
    })
);
onMounted(async () => {
    let c = await GetConfig();
    copyObjTo(c, conf.value);
    watch(conf.value, async () => {
        let config = await GetConfig();
        config.isAutoSync = conf.value.isAutoSync;
        config.isDarkTheme = conf.value.isDarkTheme;
        config.showGacha = conf.value.showGacha;
        config.language = conf.value.language;
        config.isShowRank3Item = conf.value.isShowRank3Item;
        PutConfig(config);
    });
});
</script>
<template>
    <div style="height: 100%; overflow: hidden">
        <h2>- 设置 -</h2>
        <el-scrollbar style="width: 100%; height: calc(100% - 74px)">
            <div class="box">
                <span class="title-top">全局</span>
                <SwitchItem
                    class="item"
                    @change="(val) => toggleTheme(val)"
                    v-model="conf.isDarkTheme"
                >
                    深色模式</SwitchItem
                >
                <SwitchItem class="item" v-model="conf.isAutoSync">自动同步</SwitchItem>

                <span class="title">祈愿</span>
                <SwitchItem class="item" v-model="conf.showGacha.g301">角色活动祈愿</SwitchItem>
                <SwitchItem class="item" v-model="conf.showGacha.g302">武器活动祈愿</SwitchItem>
                <SwitchItem class="item" v-model="conf.showGacha.g200">常驻祈愿</SwitchItem>
                <SwitchItem class="item" v-model="conf.showGacha.g100">新手祈愿</SwitchItem>
                <SwitchItem class="item" v-model="conf.isShowRank3Item">显示三星物品</SwitchItem>
                <div style="min-height: 10px"></div>
            </div>
        </el-scrollbar>
    </div>
</template>
<style scoped>
.item {
    margin-top: 6px;
    margin-bottom: 4px;
}
.title-top {
    font-size: 14px;
    margin-right: auto;
    color: var(--el-text-color-regular);
}
.title {
    margin-top: 16px;
    font-size: 14px;
    margin-right: auto;
    color: var(--el-text-color-regular);
}
.box {
    display: flex;
    flex-flow: column;
    margin-left: 12px;
    margin-right: 12px;
}
</style>
