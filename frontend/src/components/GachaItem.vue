<!-- 
    在 GachaData 中显示祈愿物品的组件
 -->
<script lang="ts" setup>
import { onMounted, ref } from "vue";
import l3bg from "../assets/images/l3.png";
import l4bg from "../assets/images/l4.png";
import l5bg from "../assets/images/l5.png";
import { GachaLog } from "../type";

// 进度条的颜色集合
const colors = [
    { color: "#f56c6c", percentage: 20 },
    { color: "#e6a23c", percentage: 40 },
    { color: "#5cb87a", percentage: 60 },
    { color: "#1989fa", percentage: 80 },
    { color: "#6f7ad3", percentage: 99 },
    { color: "#7191b1", percentage: 100 },
];
const backgroundImages = ref({
    now: "url(" + l3bg + ")",
    l5: "url(" + l5bg + ")",
    l4: "url(" + l4bg + ")",
    l3: "url(" + l3bg + ")",
});

// 物品不同品质的颜色
const itemColors = {
    l5: {
        text: "#edeff1",
        bg: "#f39d53",
    },
    l4: {
        text: "#e7e7e7",
        bg: "#a37dde",
    },
    l3: {
        text: "#e7e7e7",
        bg: "#5ba8d6",
    },
};
const props = defineProps<{ gachaLog: GachaLog }>();
const data = ref({
    id: "",
    name: "",
    date: "",
    // 距离上一次同等品质的物品出货所花费的次数
    usedCost: 0,
    assess: "",
    percentage: 100,
    bgColor: itemColors.l3.bg,
    textColor: itemColors.l3.text,
});

onMounted(() => {
    var d = data.value;
    var g = props.gachaLog;
    var date = new Date();

    d.name = "name";
    d.date =
        date.getFullYear() +
        "-" +
        date.getMonth() +
        "-" +
        date.getDate() +
        " " +
        date.getHours() +
        ":" +
        date.getMinutes();
    if (g.rank_type != "3") {
        // 四星，五星物品显示额外内容
        // TODO: 从后端查询
        d.usedCost = 40;
        d.percentage = (d.usedCost / 90) * 100;
        if (g.rank_type === "4") {
            d.bgColor = itemColors.l4.bg;
            d.textColor = itemColors.l4.text;
            backgroundImages.value.now = backgroundImages.value.l4;
        } else {
            d.bgColor = itemColors.l5.bg;
            d.textColor = itemColors.l5.text;
            backgroundImages.value.now = backgroundImages.value.l5;
        }
    } else {
    }
});
</script>
<template>
    <div class="item">
        <div class="icon"></div>
        <div style="width: 10px" />
        <div class="info" :id="data.id">
            <div class="item-name">
                <span>{{ data.name }}</span>
                <span> &nbsp&nbsp</span>
                <span v-if="data.usedCost != 0">{{ data.usedCost }}</span>
                <span class="date">{{ data.date }}</span>
            </div>
            <el-progress
                :text-inside="true"
                :stroke-width="24"
                :percentage="data.percentage"
                :color="colors"
                status="success"
                ><span>{{ data.assess }}</span></el-progress
            >
        </div>
    </div>
</template>
<style scoped>
.icon {
    min-height: 64px;
    min-width: 64px;
    max-width: 64px;
    max-height: 64px;
    color: #161616;
    background-image: v-bind("backgroundImages.now");
    background-color: aqua;
    background-size: contain;
    border-radius: 10px;
}
.item {
    display: flex;
    flex-direction: row;
    align-items: center;
    width: 100%;
    border-radius: 10px;
    margin-bottom: 10px;
}
.info {
    position: relative;
    margin: auto;
    height: 52px;
    width: 100%;
    display: flex;
    background-color: v-bind("data.bgColor");
    border: 92px;
    border-color: black;
    flex-direction: column;
    text-align: left;
    border-radius: 10px;
    padding-left: 10px;
    padding-right: 10px;
    padding-top: 8px;
    padding-bottom: 4px;
    border: 2px solid #d4d4d4;
    box-shadow: 0 0 12px rgb(142, 142, 142);
}
.item-name {
    font-size: 14px;
    color: v-bind("data.textColor");
    margin-bottom: 4px;
    margin-left: 4px;
}
.date {
    position: absolute;
    right: 14px;
}
</style>
