<!-- 
    在主界面上的祈愿数据概览页面，绘制和显示饼图
 -->
<script lang="ts" setup>
import Pie from "./GachaPie.vue";
import { PieData, Option } from "../type";
import { Ref, ref } from "vue";

defineProps<{ data: Option }>();
// 打开 GachaData 页面
// name: 祈愿名称, count: 几发没出
defineEmits<{ (e: "open-gacha-data-page", name: string, count: number): void }>();
const gachaInfo1: Ref<PieData> = ref({
    name: "角色活动祈愿",
    几发未出金: 30,
    三星武器: 32,
    四星角色: 311,
    四星武器: 13,
    五星武器: 23,
    五星角色: 2,
});

const gachaInfo2: Ref<PieData> = ref({
    name: "武器活动祈愿",
    几发未出金: 80,
    三星武器: 32,
    四星角色: 211,
    四星武器: 111,
    五星武器: 76,
    五星角色: 28,
});
const gachaInfo3: Ref<PieData> = ref({
    name: "常驻祈愿",
    几发未出金: 54,
    三星武器: 312,
    四星角色: 311,
    四星武器: 131,
    五星武器: 213,
    五星角色: 21,
});
const gachaInfo4: Ref<PieData> = ref({
    name: "新手祈愿",
    几发未出金: 10,
    三星武器: 3552,
    四星角色: 311,
    四星武器: 141,
    五星武器: 53,
    五星角色: 51,
});
</script>
<template>
    <!-- 祈愿概览 -->
    <div class="gacha-info">
        <!-- 角色活动祈愿 -->
        <Pie
            :data="gachaInfo1"
            @open-click="$emit('open-gacha-data-page', gachaInfo1.name, gachaInfo1.几发未出金)"
            v-if="data.showGacha.roleUp"
        />
        <!-- 武器活动祈愿 -->
        <Pie
            :data="gachaInfo2"
            @open-click="$emit('open-gacha-data-page', gachaInfo2.name, gachaInfo2.几发未出金)"
            v-if="data.showGacha.armsUp"
        />
        <!-- 常驻祈愿 -->
        <Pie
            :data="gachaInfo3"
            @open-click="$emit('open-gacha-data-page', gachaInfo3.name, gachaInfo3.几发未出金)"
            v-if="data.showGacha.permanent"
        />
        <!-- 新手祈愿 -->
        <Pie
            :data="gachaInfo4"
            @open-click="$emit('open-gacha-data-page', gachaInfo4.name, gachaInfo4.几发未出金)"
            v-if="data.showGacha.start"
        />
    </div>
</template>
<style scoped>
.gacha-info {
    display: flex;
    flex-shrink: 0;
    flex-wrap: wrap;
    flex-direction: row;
    justify-content: center;
}
</style>
