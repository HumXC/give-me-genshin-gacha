<script lang="ts" setup>
import { init, EChartsOption, use } from "echarts";
import { ElCard } from "element-plus";
import { onMounted, ref } from "vue";
import { PieData } from "../type";

// TODO: 添加更多的参数以支持不同的点击位置
defineEmits<{
    (e: "open-click"): void;
}>();
const props = defineProps<{ data: PieData }>();
onMounted(() => {
    var d = props.data;
    var opt = {
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
                name: d.name,
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
                    { value: d.三星武器, name: "三星武器" },
                    { value: d.四星角色, name: "四星角色" },
                    { value: d.四星武器, name: "四星武器" },
                    { value: d.五星角色, name: "五星角色" },
                    { value: d.五星武器, name: "五星武器" },
                ],
            },
        ],
    };
    let myChart = init(<HTMLElement>document.getElementById(d.name));
    myChart.setOption(opt);
    // TODO: 通过点击饼图上不同的色块，打开对应的 GachaData 页面
    myChart.on("click", (p) => {
        console.log(p);
    });
});
</script>

<template>
    <el-card class="gacha-box">
        <template #header>
            <div class="card-header">
                <span>{{ data.name }}</span>
            </div>
        </template>
        <div :id="data.name" class="chart-pie"></div>
        <el-tooltip show-after="1000" content="距上一次出五星的祈愿次数" placement="top">
            <el-button class="used-num" @click.prevent="$emit('open-click')" circle>{{
                data.几发未出金
            }}</el-button>
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
