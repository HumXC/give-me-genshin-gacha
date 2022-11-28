<!--
    主界面顶部的控制区域
    能够触发同步事件和打开设置界面事件
-->

<script lang="ts" setup>
import { More, Sort } from "@element-plus/icons-vue";
import { onMounted, ref } from "vue";
import { GetUids } from "../../wailsjs/go/main/App";
import { Option } from "../type";
import OptionMenu from "./OptionMenu.vue";
// 选择框里已经选择的 uid
var isUidSelectorDisabled = ref(false);

const props = defineProps<{ option: Option }>();
const data = ref({
    isSynchronizing: false,
    uids: new Array<string>(),
});
const e = defineEmits<{
    (e: "open-option-menu"): void;
    (e: "start-sync", done: () => void): void;
    (e: "change-select", uid: string): void;
}>();
defineExpose(refresh);
function handleSync() {
    data.value.isSynchronizing = true;
    // 用于取消同步按钮的 loading 状态
    let done = () => {
        data.value.isSynchronizing = false;
    };
    e("start-sync", done);
}
async function refresh() {
    // 从后端获取uid
    let uids = await GetUids();
    data.value.uids = uids;
    if (uids.length != 0) {
        if (props.option.controlBar.selectedUid == "") {
            props.option.controlBar.selectedUid = uids[0];
        }
    } else {
        props.option.controlBar.selectedUid = "";
    }
    if (uids.length <= 1) {
        isUidSelectorDisabled.value = true;
    }
}
async function init() {
    refresh();
    if (props.option.otherOption.autoSync) {
        handleSync();
    }
}
onMounted(() => {
    init();
});
</script>
<template>
    <div style="height: 10px"></div>
    <div class="control-bar">
        <!-- 同步按钮 -->
        <el-button
            type="primary"
            id="sync_button"
            :icon="Sort"
            :loading="data.isSynchronizing"
            @click="handleSync"
            circle
        >
        </el-button>
        <!-- uid 选择框 -->
        <el-select
            v-bind:disabled="isUidSelectorDisabled"
            v-model="option.controlBar.selectedUid"
            id="uid_selector"
            placeholder="先点击左侧按钮同步"
            @change="(uid:string) => $emit('change-select', uid)"
        >
            <el-option v-for="uid in data.uids" :key="uid" :label="uid" :value="uid" />
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
