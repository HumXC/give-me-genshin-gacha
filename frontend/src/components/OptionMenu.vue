<!-- 
    设置侧栏，提供对 Option 的显示与修改
 -->

<script lang="ts" setup>
import { Option } from "../type";
// TODO: 增加 close 事件，实现关闭页面时保存 Option
const e = defineEmits<{
    (e: "save-option", done: () => void): void;
}>();

defineProps<{ data: { isShow: boolean; opt: Option } }>();
const handleClose = (done: () => void) => {
    e("save-option", done);
};
</script>
<template>
    <!-- 设置侧栏 -->
    <el-drawer
        style="text-align: left"
        size="250"
        v-model="data.isShow"
        direction="rtl"
        title="选项"
        :before-close="handleClose"
    >
        <!-- 显示祈愿 -->
        <el-card>
            <template #header>
                <div class="card-header">
                    <span>显示祈愿</span>
                </div>
            </template>
            <div class="switch-item">
                <span>角色活动祈愿</span>
                <el-switch class="switch" v-model="data.opt.showGacha.roleUp" />
            </div>
            <div class="switch-item">
                <span>武器活动祈愿</span>
                <el-switch class="switch" v-model="data.opt.showGacha.armsUp" />
            </div>
            <div class="switch-item">
                <span>常驻祈愿</span>
                <el-switch class="switch" v-model="data.opt.showGacha.permanent" />
            </div>
            <div class="switch-item">
                <span>新手祈愿</span>
                <el-switch class="switch" v-model="data.opt.showGacha.start" />
            </div>
        </el-card>
        <br />
        <!-- 程序设置 -->
        <el-card>
            <template #header>
                <div class="card-header">
                    <span>程序设置</span>
                </div>
            </template>

            <div class="switch-item">
                <el-tooltip content="启动程序时自动同步数据" placement="top-start"
                    ><span>自动同步</span> </el-tooltip
                ><el-switch class="switch" v-model="data.opt.otherOption.autoSync" />
            </div>

            <div class="switch-item">
                <el-tooltip content="同步时优先使用代理服务器获取祈愿链接" placement="top-start"
                    ><span>默认使用代理</span></el-tooltip
                >
                <el-switch class="switch" v-model="data.opt.otherOption.useProxy" />
            </div>
            <div class="switch-item">
                <span>深色主题</span>
                <el-switch class="switch" v-model="data.opt.otherOption.darkTheme" />
            </div>
        </el-card>
    </el-drawer>
</template>
<style scoped>
.switch-item {
    display: flex;
    flex-direction: row;
    line-height: 35px;
    justify-content: space-between;
}
.switch {
    --el-switch-on-color: #13ce66;
    --el-switch-off-color: #ff4949;
}
</style>
