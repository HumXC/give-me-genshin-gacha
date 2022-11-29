<!-- 
    在主界面上的祈愿数据概览页面，绘制和显示饼图
 -->
<script lang="ts" setup>
import { Ref, ref } from "vue";
import { Option, PieData } from "../type";
import Pie from "./GachaPie.vue";

defineProps<{ option: Option }>();
// 打开 GachaData 页面
defineEmits<{
    (e: "pie-click", gachaType: string, showType: string): void;
}>();
const gachaInfo301: Ref<PieData> = ref({
    gachaType: "301",
    usedCost: 30,
    arms3Total: 32,
    role4Total: 311,
    arms4Total: 13,
    arms5Total: 23,
    role5Total: 2,
});

const gachaInfo302: Ref<PieData> = ref({
    gachaType: "302",
    usedCost: 80,
    arms3Total: 32,
    role4Total: 211,
    arms4Total: 111,
    arms5Total: 76,
    role5Total: 28,
});
const gachaInfo200: Ref<PieData> = ref({
    gachaType: "200",
    usedCost: 10,
    arms3Total: 3552,
    role4Total: 311,
    arms4Total: 141,
    arms5Total: 53,
    role5Total: 51,
});
const gachaInfo100: Ref<PieData> = ref({
    gachaType: "100",
    usedCost: 54,
    arms3Total: 312,
    role4Total: 311,
    arms4Total: 131,
    arms5Total: 213,
    role5Total: 21,
});
</script>
<template>
    <!-- 祈愿概览 -->
    <div class="gacha-info">
        <!-- 角色活动祈愿 -->
        <Pie
            :data="gachaInfo301"
            @pie-click="(gachaType, showType) => $emit('pie-click', gachaType, showType)"
            v-if="option.showGacha.roleUp"
        />
        <!-- 武器活动祈愿 -->
        <Pie
            :data="gachaInfo302"
            @pie-click="(gachaType, showType) => $emit('pie-click', gachaType, showType)"
            v-if="option.showGacha.armsUp"
        />
        <!-- 常驻祈愿 -->
        <Pie
            :data="gachaInfo200"
            @pie-click="(gachaType, showType) => $emit('pie-click', gachaType, showType)"
            v-if="option.showGacha.permanent"
        />
        <!-- 新手祈愿 -->
        <Pie
            :data="gachaInfo100"
            @pie-click="(gachaType, showType) => $emit('pie-click', gachaType, showType)"
            v-if="option.showGacha.start"
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
