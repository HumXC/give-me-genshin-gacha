<script lang="ts" setup></script>
<template>
    <div style="height: 100%; overflow: hidden">
        <div class="header">
            <h2>- 祈愿 -</h2>
        </div>
        <el-scrollbar style="width: 100%; height: calc(100% - 75px)" :view-style="''">
            <div class="info">
                <h3>角色活动祈愿 - 100</h3>
                <p>五星武器 - 20</p>
                <p>五星角色 - 20</p>
                <p>四星武器 - 10</p>
                <p>四星角色 - 49</p>
            </div>
            <div class="info">
                <h3>武器活动祈愿 - 100</h3>
                <p>五星武器 - 20</p>
                <p>五星角色 - 20</p>
                <p>四星武器 - 10</p>
                <p>四星角色 - 49</p>
            </div>
            <div class="info">
                <h3>常驻祈愿 - 100</h3>
                <p>五星武器 - 20</p>
                <p>五星角色 - 20</p>
                <p>四星武器 - 10</p>
                <p>四星角色 - 49</p>
            </div>
            <div class="info">
                <h3>新手祈愿 - 100</h3>
                <p>五星武器 - 20</p>
                <p>五星角色 - 20</p>
                <p>四星武器 - 10</p>
                <p>四星角色 - 49</p>
            </div>
        </el-scrollbar>
    </div>
</template>
<style scoped>
.info {
    margin-left: 12px;
    margin-right: 12px;
    overflow: hidden;
    margin-bottom: 16px;
    text-align: left;
    padding: 12px;
    background-color: var(--item);
    border-radius: 12px;
}
.info:hover {
    background-color: var(--item-hover);
}
.header {
    background-color: var(--bg);
}
</style>
