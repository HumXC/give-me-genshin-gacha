<script lang="ts" setup>
import { ref } from "vue";
const value = ref(false);
</script>
<template>
    <h1 style="height: 100px">Hi !</h1>
    <div class="view">
        <div class="left">
            <p>一些信息</p>
            <p>一些信息</p>
            <p>一些信息</p>
            <p>一些信息</p>
            <p>一些信息</p>
        </div>
        <div class="right">
            <div style="height: 50px"></div>
            <el-select> </el-select>
            <div style="height: 50px"></div>
            <div class="sync">
                <el-button circle class="sync-button" type="primary">尝试同步</el-button>
                <el-switch
                    v-model="value"
                    style="
                        --el-switch-on-color: var(--el-color-success);
                        --el-switch-off-color: var(--el-color-success);
                    "
                    class="sync-type"
                    active-text="从网络代理"
                    inactive-text="从游戏缓存"
                />
            </div>
        </div>
    </div>
</template>
<style scoped>
.sync-button {
    height: calc(24vh);
    width: calc(24vh);
    position: absolute;
    top: 18%;
}

.sync-type {
    position: absolute;
    bottom: 10px;
    color: var(--el-text-color-regular);
}
.view {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
}
.left {
    padding-left: 16px;
    text-align: left;
    width: 46%;
    border-radius: 8px;
    height: 100%;
    background-color: var(--el-fill-color-lighter);
}
.right {
    flex: 1;
    margin-left: 5%;
    display: flex;
    flex-flow: column;
    height: 100%;
}
.sync {
    position: relative;
    height: 100%;
    padding-left: 20px;
    padding-right: 20px;
    text-align: left;
    border-radius: 8px;
    display: flex;
    flex-flow: column;
    align-items: center;
    background-color: var(--el-fill-color-lighter);
}
</style>
