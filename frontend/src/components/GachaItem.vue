<!-- 
    在 GachaData 中显示祈愿物品的组件
 -->
<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { models } from "../../wailsjs/go/models";
import rank3bg from "../assets/images/rank3.png";
import rank4bg from "../assets/images/rank4.png";
import rank5bg from "../assets/images/rank5.png";
let GetNumWithLast = (a: any, b: any, c: any) => {
    return 0;
};
const iconSrc = ref("");
// 进度条的颜色集合
const progressColor = {
    // 进度条底色
    bg: {
        now: "#557de1",
        rank4: "#8f62d4",
        rank5: "#da8d4a",
    },
    // 进度条进度的颜色
    in: {
        now: "#70bff3",
        rank4: "#d9a2e7de",
        rank5: "#ffd52c",
    },
};
const backgroundImages = ref({
    now: "url(" + rank3bg + ")",
    rank5: "url(" + rank5bg + ")",
    rank4: "url(" + rank4bg + ")",
});

// 物品不同品质的颜色
const itemColors = {
    rank5: {
        text: "#fffdf5",
        bg: "#f39d53",
    },
    rank4: {
        text: "#fff2fa",
        bg: "#a37dde",
    },
    rank3: {
        text: "#fdfeff",
        bg: "#5ba8d6",
    },
};
const props = defineProps<{ uid: string; gachaLog: models.GachaLog }>();
const data = ref({
    // 距离上一次同等品质的物品出货所花费的次数
    usedCost: 0,
    percentage: 100,
    bgColor: itemColors.rank3.bg,
    textColor: itemColors.rank3.text,
});

onMounted(async () => {
    var d = data.value;
    var g = props.gachaLog;
    iconSrc.value = "icon/" + props.gachaLog.name + ".png";
    if (g.rank_type === "3") return;
    // 四星，五星物品显示额外内容
    d.usedCost = await GetNumWithLast(props.uid, props.gachaLog.gacha_type, props.gachaLog.id);
    if (g.rank_type === "4") {
        progressColor.bg.now = progressColor.bg.rank4;
        progressColor.in.now = progressColor.in.rank4;
        d.percentage = (d.usedCost / 10) * 100;
        if (d.percentage < 11) d.percentage = 11;
        d.bgColor = itemColors.rank4.bg;
        d.textColor = itemColors.rank4.text;
        backgroundImages.value.now = backgroundImages.value.rank4;
    } else {
        progressColor.bg.now = progressColor.bg.rank5;
        progressColor.in.now = progressColor.in.rank5;
        d.percentage = (d.usedCost / 90) * 100;
        if (d.percentage < 11) d.percentage = 11;
        d.bgColor = itemColors.rank5.bg;
        d.textColor = itemColors.rank5.text;
        backgroundImages.value.now = backgroundImages.value.rank5;
    }
});
</script>
<template>
    <div class="item">
        <img :alt="gachaLog.name" class="icon" :src="iconSrc" />
        <div style="width: 10px" />
        <div class="info" :id="gachaLog.id">
            <div class="item-name">
                <span>{{ gachaLog.name }}</span>
                <!-- 下面这行不能删，删了的话如果进度条满了，四星，五星物品的进度条颜色就会变成三星的蓝色。但是为什么？ -->
                <span v-if="data.usedCost != 0"></span>
                <span class="time">{{ gachaLog.time }}</span>
            </div>
            <el-progress
                :text-inside="true"
                :stroke-width="24"
                :percentage="data.percentage"
                :color="progressColor.in.now"
                status="success"
            >
                <span class="used-cost">
                    {{ data.usedCost !== 0 ? data.usedCost : "" }}
                </span>
            </el-progress>
        </div>
    </div>
</template>
<style scoped>
:deep(.el-progress-bar__outer) {
    background-color: v-bind("progressColor.bg.now");
}
.used-cost {
    margin-right: 2px;
    font-size: 16px;
}
.icon {
    min-height: 64px;
    min-width: 64px;
    max-width: 64px;
    max-height: 64px;
    color: #161616;
    background-image: v-bind("backgroundImages.now");
    background-size: contain;
    border-radius: 10px;
}
.item {
    display: flex;
    flex-direction: row;
    align-items: center;
    height: 70px;
    width: 320px;
    margin-top: 10px;
    margin-right: 16px;
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
    border: 2px solid #d0d0d0;
}
.item-name {
    font-size: 14px;
    color: v-bind("data.textColor");
    margin-bottom: 4px;
    margin-left: 4px;
}
.time {
    position: absolute;
    right: 14px;
}
</style>
