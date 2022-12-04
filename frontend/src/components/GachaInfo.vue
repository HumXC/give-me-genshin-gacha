<!-- 
    在主界面上的祈愿数据概览页面，绘制和显示饼图
 -->
<script lang="ts" setup>
import { Ref, ref } from "vue";
import { GetPieDatas } from "../../wailsjs/go/main/App";
import { Option, PieData } from "../type";
import Pie from "./GachaPie.vue";

const props = defineProps<{ option: Option }>();
// 打开 GachaData 页面
defineEmits<{
    (e: "pie-click", gachaType: string): void;
}>();

const datas: Ref<PieData> = ref(new PieData());
const refresh = async () => {
    if (props.option.controlBar.selectedUid === "") {
        return;
    }
    let d = await GetPieDatas(props.option.controlBar.selectedUid);
    datas.value.usedCosts = d.usedCosts;
    // 设置之间按钮的值
    for (let i = 0; i < d.usedCosts.length; i++) {
        const element = d.usedCosts[i];
        switch (element.gachaType) {
            case "301":
                pie1.value.usedCost = element.cost;
                break;
            case "302":
                pie2.value.usedCost = element.cost;
                break;
            case "200":
                pie3.value.usedCost = element.cost;
                break;
            case "100":
                pie4.value.usedCost = element.cost;
                break;
        }
    }

    pie1.value.items = d.totals.t301;
    pie2.value.items = d.totals.t302;
    pie3.value.items = d.totals.t200;
    pie4.value.items = d.totals.t100;

    refreshPie1.value();
    refreshPie2.value();
    refreshPie3.value();
    refreshPie4.value();
};
defineExpose(refresh);
const pie1 = ref({
    gachaType: "301",
    usedCost: 0,
    items: new Array<{
        rankType: string;
        itemType: string;
        total: number;
    }>(),
});
const pie2 = ref({
    gachaType: "302",
    usedCost: 0,
    items: new Array<{
        rankType: string;
        itemType: string;
        total: number;
    }>(),
});
const pie3 = ref({
    gachaType: "200",
    usedCost: 0,
    items: new Array<{
        rankType: string;
        itemType: string;
        total: number;
    }>(),
});
const pie4 = ref({
    gachaType: "100",
    usedCost: 0,
    items: new Array<{
        rankType: string;
        itemType: string;
        total: number;
    }>(),
});
const refreshPie1 = ref(() => {});
const refreshPie2 = ref(() => {});
const refreshPie3 = ref(() => {});
const refreshPie4 = ref(() => {});
</script>
<template>
    <!-- 祈愿概览 -->
    <div class="gacha-info">
        <!-- 角色活动祈愿 -->
        <Pie
            :data="pie1"
            @pie-click="(gachaType) => $emit('pie-click', gachaType)"
            v-if="option.showGacha.roleUp"
            ref="refreshPie1"
        />
        <!-- 武器活动祈愿 -->
        <Pie
            :data="pie2"
            @pie-click="(gachaType) => $emit('pie-click', gachaType)"
            v-if="option.showGacha.armsUp"
            ref="refreshPie2"
        />
        <!-- 常驻祈愿 -->
        <Pie
            :data="pie3"
            @pie-click="(gachaType) => $emit('pie-click', gachaType)"
            v-if="option.showGacha.permanent"
            ref="refreshPie3"
        />
        <!-- 新手祈愿 -->
        <Pie
            :data="pie4"
            @pie-click="(gachaType) => $emit('pie-click', gachaType)"
            v-if="option.showGacha.start"
            ref="refreshPie4"
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
