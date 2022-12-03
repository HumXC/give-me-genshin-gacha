<!-- 
    详细祈愿数据的展示页面
 -->

<script lang="ts" setup>
import { Close } from "@element-plus/icons-vue";
import { ref } from "vue";
import GachaItem from "./GachaItem.vue";
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
const count = ref(0);
const load = () => {
    console.log("load");
    count.value += 10;
};
</script>
<template>
    <el-drawer
        style="padding: 0px"
        size="100%"
        :show-close="false"
        v-model="data.isShow"
        direction="ttb"
    >
        <template #title>
            <div class="title">
                <span>角色活动祈愿</span>
                <el-icon class="close-btn" @click="data.isShow = false">
                    <Close />
                </el-icon>
            </div>
        </template>
        <el-scrollbar height="100%"
            ><div
                class="item"
                v-infinite-scroll="load"
                infinite-scroll-delay="20"
                infinite-scroll-distance="100"
            >
                <div v-for="i in count"><GachaItem :gacha-log="gachaLogs[0]" /></div></div
        ></el-scrollbar>
    </el-drawer>
</template>
<style scoped>
.item-list {
    list-style: none;
    padding: 0px;
    padding-inline-start: 0px;
    height: 100%;
    margin: 0;
    list-style: none;
}
.item {
    height: 70px;
    margin-right: 16px;
}
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
    color: var(--el-color-primary);
    cursor: pointer;
}
</style>
