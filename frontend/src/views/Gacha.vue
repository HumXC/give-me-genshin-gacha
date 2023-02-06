<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { GetConfig } from "../../wailsjs/go/app/App";
import { GetGachaInfo } from "../../wailsjs/go/app/GachaMan";
import { models } from "../../wailsjs/go/models";
import { gachaTypeToName } from "../util";
const gachasInfo = ref(new Array<models.GachaInfo>());
let fakeData = [
    {
        gachaType: "301",
        allCount: 100,
        avatar5: 3,
        avatar4: 54,
        weapon5: 13,
        weapon4: 0,
        weapon3: 0,
    },
    {
        gachaType: "302",
        allCount: 100,
        avatar5: 15,
        avatar4: 12,
        weapon5: 23,
        weapon4: 24,
        weapon3: 31,
    },
    {
        gachaType: "200",
        allCount: 200,
        avatar5: 12,
        avatar4: 32,
        weapon5: 10,
        weapon4: 13,
        weapon3: 45,
    },
    {
        gachaType: "100",
        allCount: 201,
        avatar5: 31,
        avatar4: 12,
        weapon5: 21,
        weapon4: 21,
        weapon3: 32,
    },
];

const showGacha = ref({
    g301: false,
    g302: false,
    g200: false,
    g100: false,
});
const isShowRank3Item = ref(false);
function getPercentage(n1: number, n2: number): string {
    return ((n2 / n1) * 100).toFixed(2);
}
function isShowGacha(gachaType: string): boolean {
    let show = showGacha.value;
    switch (gachaType) {
        case "301":
            return show.g301;
        case "302":
            return show.g302;
        case "200":
            return show.g200;
        case "100":
            return show.g100;
        default:
            return true;
    }
}
onMounted(async () => {
    let c = await GetConfig();
    showGacha.value = c.showGacha;
    isShowRank3Item.value = c.isShowRank3Item;
    gachasInfo.value = await GetGachaInfo();
    gachasInfo.value.sort((a, b) => {
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
});
</script>
<template>
    <div style="height: 100%; overflow: hidden">
        <div class="header">
            <h2>- 祈愿 -</h2>
        </div>
        <el-scrollbar style="width: 100%; height: calc(100% - 75px)">
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
    </div>
</template>
<style scoped>
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
