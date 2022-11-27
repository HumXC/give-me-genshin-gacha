<!-- 
    详细祈愿数据的展示页面
 -->

<script lang="ts" setup>
import GachaItem from "./GachaItem.vue";
import { Close } from "@element-plus/icons-vue";
import { ref } from "vue";
const props = defineProps<{ data: { isShow: boolean } }>();

const gachaLogs = ref([
    {
        uid: "111111111",
        gacha_type: "100",
        item_id: "",
        count: "1",
        time: 1669293032,
        name: "飞天御剑",
        lang: "zh-cn",
        item_type: "武器",
        rank_type: "5",
        id: "1665741960038126786",
    },
    {
        uid: "111111111",
        gacha_type: "100",
        item_id: "",
        count: "1",
        time: 1669293032,
        name: "飞天御剑",
        lang: "zh-cn",
        item_type: "武器",
        rank_type: "5",
        id: "1665741960038126786",
    },
    {
        uid: "111111111",
        gacha_type: "100",
        item_id: "",
        count: "1",
        time: 1669293032,
        name: "胡桃",
        lang: "zh-cn",
        item_type: "武器",
        rank_type: "4",
        id: "1665741960038126786",
    },
    {
        uid: "111111111",
        gacha_type: "100",
        item_id: "",
        count: "1",
        time: 1669293032,
        name: "三星",
        lang: "zh-cn",
        item_type: "武器",
        rank_type: "3",
        id: "1665741960038126786",
    },
    {
        uid: "111111111",
        gacha_type: "100",
        item_id: "",
        count: "1",
        time: 1669293032,
        name: "胡桃",
        lang: "zh-cn",
        item_type: "武器",
        rank_type: "4",
        id: "1665741960038126786",
    },
    {
        uid: "111111111",
        gacha_type: "100",
        item_id: "",
        count: "1",
        time: 1669293032,
        name: "三星",
        lang: "zh-cn",
        item_type: "武器",
        rank_type: "3",
        id: "1665741960038126786",
    },
]);
</script>
<template>
    <el-drawer size="100%" :show-close="false" v-model="data.isShow" direction="ttb">
        <template #title>
            <div class="title">
                <span>角色活动祈愿</span>
                <el-icon class="close-btn" @click="data.isShow = false">
                    <Close />
                </el-icon>
            </div>
        </template>
        <div class="items" v-for="gachaLog in gachaLogs">
            <GachaItem :gacha-log="gachaLog" />
        </div>
    </el-drawer>
</template>
<style scoped>
.title {
    position: relative;
    width: 100%;
}
.close-btn {
    position: absolute;
    width: 20px;
    height: 40px;
    top: -10px;
    right: 10px;
}
.close-btn :hover {
    color: aqua;
    cursor: pointer;
}
</style>
