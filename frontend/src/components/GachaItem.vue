<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { models } from "../../wailsjs/go/models";
import r3bg from "../assets/images/rank3.png";
import r4bg from "../assets/images/rank4.png";
import r5bg from "../assets/images/rank5.png";
import { formatTime } from "../util";
const prop = defineProps<{ data: models.FullGachaLog }>();
const iconBg = ref(r3bg);
const iconSrc = ref("/icon/gacha_item/" + prop.data.itemId + ".png");
const progIn = ref("var(--gacha-item-prog-in-");
const progBg = ref("var(--gacha-item-prog-bg-");
const bg = ref("var(--gacha-item-bg-");
const text = ref("var(--gacha-item-text-");
const progPercentage = ref(100);
onMounted(() => {
    switch (prop.data.rankType) {
        case 4:
            iconBg.value = r4bg;
            progBg.value = progBg.value + 4 + ")";
            progIn.value = progIn.value + 4 + ")";
            bg.value = bg.value + 4 + ")";
            text.value = text.value + 4 + ")";
            progPercentage.value = prop.data.cost * 10;
            break;
        case 5:
            iconBg.value = r5bg;
            progBg.value = progBg.value + 5 + ")";
            progIn.value = progIn.value + 5 + ")";
            bg.value = bg.value + 5 + ")";
            text.value = text.value + 5 + ")";
            // 武器活动祈愿 80 保底，其他 90 保底
            if (prop.data.gachaType == "302") {
                progPercentage.value = (prop.data.cost / 80) * 100;
            } else {
                progPercentage.value = (prop.data.cost / 90) * 100;
            }
            break;
        default:
            progBg.value = progBg.value + 3 + ")";
            progIn.value = progIn.value + 3 + ")";
            bg.value = bg.value + 3 + ")";
            text.value = text.value + 3 + ")";
    }
});
</script>

<template>
    <div class="item">
        <img
            :alt="data.name"
            class="icon"
            :style="'background-image:url(' + iconBg + ')'"
            :src="iconSrc"
        />
        <div style="width: 10px" />
        <div class="info">
            <div class="item-name">
                <span>{{ data.name }}&nbsp;{{ data.cost === 1 ? "" : data.cost }}</span>
                <span class="time">{{ formatTime(data.time) }}</span>
            </div>
            <el-progress
                class="prog"
                :text-inside="true"
                :stroke-width="24"
                :percentage="progPercentage"
                :color="progIn"
            >
                <!-- 空 span 占位用 -->
                <span> </span>
            </el-progress>
        </div>
    </div>
</template>
<style scoped>
.time {
    position: absolute;
    top: 6px;
    right: 12px;
}
.prog {
    margin-top: 4px;
    margin-left: 6px;
    margin-right: 8px;
}
.item-name {
    margin-top: 5px;
    margin-left: 10px;
    color: v-bind("text");
}
.time {
    color: v-bind("text");
    margin-bottom: 12px;
}
.info {
    position: relative;
    height: 62px;
    width: 100%;
    background-color: v-bind("bg");
    text-align: left;
    border-radius: 10px;
    border: 2px solid var(--el-border-color);
}
:deep(.el-progress-bar__outer) {
    background-color: v-bind("progBg");
}
.icon {
    min-height: 64px;
    min-width: 64px;
    max-width: 64px;
    max-height: 64px;
    background-size: contain;
    border-radius: 10px;
}
.item {
    margin-left: 12px;
    margin-right: 12px;
    display: flex;
    align-items: center;
}
</style>
