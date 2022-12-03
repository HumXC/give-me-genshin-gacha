<script lang="ts" setup>
import { init } from "echarts";
import { ElCard } from "element-plus";
import { onMounted, ref } from "vue";
import { gachaTypeToName } from "../main";

// gachaType: 祈愿类型，showType: 显示类型
// all 是点击饼图中心按钮时传递的参数，其他是点击饼图某一块传递对应的参数
const emit = defineEmits<{
    (e: "pie-click", gachaType: string): void;
}>();
const props = defineProps<{
    data: {
        gachaType: string;
        usedCost: number;
        items: Array<{
            rankType: string;
            itemType: string;
            total: number;
        }>;
    };
}>();
const gachaName = ref(gachaTypeToName(props.data.gachaType));
const rank3ItemTotal = ref("0/0");
const refresh = () => {
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
                data: [],
            },
        ],
    };
    let total = 0;
    d.items.forEach((e) => {
        if (e.rankType == "3") {
            total += e.total;
            rank3ItemTotal.value = e.total.toString();
        } else {
            let data = <Array<{ value: number; name: string }>>opt.series?.[0].data;
            data.push({
                value: e.total,
                name: e.rankType + "星" + e.itemType,
            });
            total += e.total;
        }
    });
    rank3ItemTotal.value += "/" + total;
    let myChart = init(<HTMLElement>document.getElementById(gachaName.value));
    myChart.setOption(opt);
};
defineExpose(refresh);
onMounted(() => {
    refresh();
});
</script>

<template>
    <el-card class="gacha-box">
        <template #header>
            <div class="card-header">
                <span>{{ gachaName }}</span>
            </div>
        </template>
        <span class="rank3-arms">3星: {{ rank3ItemTotal }}</span>
        <div :id="gachaName" class="chart-pie"></div>
        <el-tooltip show-after="1000" content="距上一次出五星的祈愿次数" placement="top">
            <el-button
                class="used-num"
                @click.prevent="$emit('pie-click', data.gachaType)"
                circle
                >{{ data.usedCost }}</el-button
            >
        </el-tooltip>
    </el-card>
</template>
<style scoped>
.rank3-arms {
    position: absolute;
    left: 10px;
    top: 70px;
}
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
