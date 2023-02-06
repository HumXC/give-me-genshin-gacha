<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { GetConfig } from "../../wailsjs/go/app/App";
import { GetGachaInfo } from "../../wailsjs/go/app/GachaMan";
import { config, models } from "../../wailsjs/go/models";
import { gachaTypeToName } from "../util";
const gachasInfo = ref(new Array<models.GachaInfo>());
const showGacha = ref(new config.ShowGacha());
const isShowRank3Item = ref(false);
const isEmpty = ref(true);
function getPercentage(n1: number, n2: number): number {
    return parseFloat(((n2 / n1) * 100).toFixed(2));
}

// 此计数如果对于 0, 则代表当前页面是空的, isEmpty 赋值为 true
let count = 0;
function isShowGacha(gachaType: string): boolean {
    let show = showGacha.value;
    let result = false;
    switch (gachaType) {
        case "301":
            result = show.g301;
            break;
        case "302":
            result = show.g302;
            break;
        case "200":
            result = show.g200;
            break;
        case "100":
            result = show.g100;
            break;
        default:
            result = true;
    }
    if (result === true) count++;
    if (count !== 0) isEmpty.value = false;
    return result;
}
onMounted(async () => {
    let c = await GetConfig();
    showGacha.value = c.showGacha;
    isShowRank3Item.value = c.isShowRank3Item;
    let info = await GetGachaInfo(c.selectedUid);
    info.sort((a, b) => {
        // 如果是角色和武器活动祈愿一起比较,把角色活动祈愿放前面
        if (
            (a.gachaType == "301" && b.gachaType == "302") ||
            (b.gachaType == "301" && a.gachaType == "302")
        ) {
            if (a.gachaType > b.gachaType) return 1;
            return -1;
        }
        if (a.gachaType > b.gachaType) return -1;
        return 1;
    });
    gachasInfo.value = info;
});
</script>
<template>
    <div class="你好" style="height: 100%; overflow: hidden">
        <h2>- 祈愿 -</h2>
        <el-scrollbar style="width: 100%">
            <div v-for="info in gachasInfo" :key="info.gachaType">
                <div
                    class="info"
                    v-if="isShowGacha(info.gachaType)"
                    @click="
                        $router.push({
                            path: '/gacha/logs',
                            query: { gachaType: info.gachaType },
                        })
                    "
                >
                    <h3 class="gacha-name">
                        {{ gachaTypeToName(info.gachaType) }} - {{ info.allCount }}
                    </h3>
                    <div v-if="info.avatar5 != 0">
                        <span class="item-type">五星角色 {{ info.avatar5 }}</span>
                        <el-progress
                            :color="'var(--gacha-item-prog-in-5)'"
                            :percentage="getPercentage(info.allCount, info.avatar5)"
                        />
                    </div>
                    <div v-if="info.weapon5 != 0">
                        <span class="item-type">五星武器 {{ info.weapon5 }}</span>
                        <el-progress
                            :color="'var(--gacha-item-prog-in-5)'"
                            :percentage="getPercentage(info.allCount, info.weapon5)"
                        />
                    </div>
                    <div v-if="info.avatar4 != 0">
                        <span class="item-type">四星角色 {{ info.avatar4 }}</span>
                        <el-progress
                            :color="'var(--gacha-item-prog-in-4)'"
                            :percentage="getPercentage(info.allCount, info.avatar4)"
                        />
                    </div>
                    <div v-if="info.weapon4 != 0">
                        <span class="item-type">四星武器 {{ info.weapon4 }}</span>
                        <el-progress
                            :color="'var(--gacha-item-prog-in-4)'"
                            :percentage="getPercentage(info.allCount, info.weapon4)"
                        />
                    </div>
                    <div v-if="info.weapon3 != 0 && isShowRank3Item">
                        <span class="item-type">三星武器 {{ info.weapon3 }}</span>
                        <el-progress
                            :color="'var(--gacha-item-prog-in-3)'"
                            :percentage="getPercentage(info.allCount, info.weapon3)"
                        />
                    </div>
                </div>
            </div>
        </el-scrollbar>
        <h2 class="empty" v-if="isEmpty">空空如也</h2>
    </div>
</template>
<style scoped>
.empty {
    letter-spacing: 4px;
    position: absolute;
    left: calc(50% - 56px);
    top: 30%;
    color: var(--el-text-color-disabled);
}
.gacha-name {
    margin-top: 12px;
    margin-bottom: 4px;
}
.item-type {
    display: block;
    margin-top: 12px;
}
.info {
    cursor: pointer;
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
</style>
