<script lang="ts" setup>
import { init } from "echarts";
import { ElCard } from "element-plus";
import { onMounted, ref } from "vue";
import { gachaTypeToName } from "../main";
import { PieData } from "../type";

// gachaType: 祈愿类型，showType: 显示类型
// all 是点击饼图中心按钮时传递的参数，其他是点击饼图某一块传递对应的参数
const emit = defineEmits<{
    (e: "pie-click", gachaType: string, showType: string): void;
}>();
const props = defineProps<{ data: PieData }>();
const gachaName = ref(gachaTypeToName(props.data.gachaType));
onMounted(() => {
    let d = props.data;
    let opt = {
        tooltip: {
            trigger: "item",
        },
        legend: {
            bottom: 0,
            left: "center",
        },
        series: [
            {
                top: 10,
                bottom: 60,
                name: gachaName.value,
                type: "pie",
                radius: ["10%", "100%"],
                avoidLabelOverlap: false,
                itemStyle: {
                    borderRadius: 10,
                    borderColor: "#fff",
                    borderWidth: 3,
                },
                label: {
                    show: false,
                },
                labelLine: {
                    show: false,
                },
                data: [
                    { value: d.arms3Total, name: "三星武器" },
                    { value: d.role4Total, name: "四星角色" },
                    { value: d.arms4Total, name: "四星武器" },
                    { value: d.role5Total, name: "五星角色" },
                    { value: d.arms5Total, name: "五星武器" },
                ],
            },
        ],
    };
    let myChart = init(<HTMLElement>document.getElementById(gachaName.value));
    myChart.setOption(opt);
    myChart.on("click", (p: any) => {
        let type = "all";
        switch (p.data.name) {
            case "三星武器":
                type = "arms3";
                break;
            case "四星武器":
                type = "arms4";
                break;
            case "五星武器":
                type = "arms5";
                break;
            case "四星角色":
                type = "role4";
                break;
            case "五星角色":
                type = "role5";
                break;
        }
        emit("pie-click", props.data.gachaType, type);
    });
});
</script>

<template>
    <el-card class="gacha-box">
        <template #header>
            <div class="card-header">
                <span>{{ gachaName }}</span>
            </div>
        </template>
        <div :id="gachaName" class="chart-pie"></div>
        <el-tooltip show-after="1000" content="距上一次出五星的祈愿次数" placement="top">
            <el-button
                class="used-num"
                @click.prevent="$emit('pie-click', data.gachaType, 'all')"
                circle
                >{{ data.usedCost }}</el-button
            >
        </el-tooltip>
    </el-card>
</template>
<style scoped>
.gacha-box {
    align-items: center;
    min-width: 340px;
    min-height: 340px;
    margin-top: 20px;
    margin-left: 10px;
    margin-right: 10px;
    padding-top: 0;
    position: relative;
}
.chart-pie {
    margin: 0;
    height: 300px;
    width: 300px;
    background-color: transparent;
}
.used-num {
    height: 120px;
    width: 120px;
    position: absolute;
    top: 145px;
    left: 110px;
    font-size: 30px;
}
</style>
