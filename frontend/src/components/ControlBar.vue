<!--
    主界面顶部的控制区域
    能够触发同步事件和打开设置界面事件
-->

<script lang="ts" setup>
import { More, Sort } from "@element-plus/icons-vue";
import { ref } from "vue";
const selectUid = ref("");
var uids = new Array<string>();
var isUidSelectorDisabled = ref(false);
const syncButton = ref({
    isCancel: false,
    isSyncying: false,
    tip: "",
});
defineEmits<{
    (e: "open-option-menu"): void;
}>();
function init() {
    // 从后端获取uid
    uids = new Array<string>();
    selectUid.value = uids[0];
    if (uids.length <= 1) {
        isUidSelectorDisabled.value = true;
    }
}
function sync() {
    var data = syncButton.value;
    data.isSyncying = !data.isSyncying;
}
init();
</script>
<template>
    <div style="height: 10px"></div>
    <div class="control-bar">
        <!-- 同步按钮 -->
        <el-button
            type="primary"
            id="sync_button"
            :icon="Sort"
            :loading="syncButton.isSyncying"
            @click="sync"
            circle
        >
        </el-button>
        <!-- uid 选择框 -->
        <el-select
            v-bind:disabled="isUidSelectorDisabled"
            v-model="selectUid"
            id="uid_selector"
            placeholder="先点击左侧按钮同步"
        >
            <el-option v-for="uid in uids" :key="uid" :label="uid" :value="uid" />
        </el-select>
        <!-- 更多选项按钮 -->
        <el-button
            type="primary"
            @click.prevent="$emit('open-option-menu')"
            :icon="More"
            id="setting_button"
            circle
        />
    </div>
</template>

<style scoped>
.control-bar {
    display: flex;
    flex-direction: row;
    line-height: 35px;
    justify-content: space-between;
    padding-left: 20px;
    padding-right: 20px;
}
</style>
