<script lang="ts" setup>
import { stringLiteral } from "@babel/types";
import { reactive } from "vue";
import { Uids, Sync, GetGachaInfo } from "../../wailsjs/go/main/App";
const data = reactive({
    name: "",
    uids: new Array<string>(),
    info: "",
    gachaInfo: {},
});

type GachaInfo = {
    // 池子名称
    name: string;
    // 获得星数物品的数量
    s3: number;
    s4: number;
    s5: number;
    s5Items: Array<{
        name: string;
        usec: number;
        time: number;
    }>;
};
function greet() {
    console.log("hi");
}
function sync() {
    Sync();
}
function init() {
    Uids().then((u) => (data.uids = u));
}
init();
function getGachaInfo() {
    var uid = data.uids[0];
    var g: GachaInfo = {
        name: "",
        s3: 0,
        s4: 0,
        s5: 0,
        s5Items: [],
    };
    GetGachaInfo(uid, "301").then((d) => (data.gachaInfo = d));
}
</script>

<template>
    <main>
        <div id="result" class="result">{{ data.uids }}</div>
        <div>{{ data.gachaInfo }}</div>
        <button class="btn" v-on:click="getGachaInfo">获取信息</button>
        <button class="btn" v-on:click="sync">同步</button>
        <div id="input" class="input-box">
            <input id="name" v-model="data.name" autocomplete="off" class="input" type="text" />
            <button class="btn" @click="greet">Greet</button>
        </div>
    </main>
</template>

<style scoped>
.result {
    height: 20px;
    line-height: 20px;
    margin: 1.5rem auto;
}

.input-box .btn {
    width: 60px;
    height: 30px;
    line-height: 30px;
    border-radius: 3px;
    border: none;
    margin: 0 0 0 20px;
    padding: 0 8px;
    cursor: pointer;
}

.input-box .btn:hover {
    background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
    color: #333333;
}

.input-box .input {
    border: none;
    border-radius: 3px;
    outline: none;
    height: 30px;
    line-height: 30px;
    padding: 0 10px;
    background-color: rgba(240, 240, 240, 1);
    -webkit-font-smoothing: antialiased;
}

.input-box .input:hover {
    border: none;
    background-color: rgba(255, 255, 255, 1);
}

.input-box .input:focus {
    border: none;
    background-color: rgba(255, 255, 255, 1);
}
</style>
