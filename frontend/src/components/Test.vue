<script lang="ts" setup>
import { GetMsg, PutMsg } from "../../wailsjs/go/main/Bridge";

import { ref } from "vue";
import { type } from "os";

function putmsg(type: string) {
    var m = {
        type: type,
        msg: "",
    };
    switch (type) {
        case "success":
            m.msg = "成功";
            break;
        case "error":
            m.msg = "错误";
            break;
        case "warning":
            m.msg = "警告";
            break;
        default:
            m.msg = "普通消息";
    }
    console.log("dadqa");
    PutMsg(m);
}
function putM() {
    putmsg("message");
}
function putE() {
    putmsg("error");
}
function putW() {
    putmsg("warning");
}
function putS() {
    putmsg("success");
}
</script>
<template>
    <el-button @click="putM"> 发送消息 </el-button>
    <el-button @click="putE"> 发送错误 </el-button>
    <el-button @click="putW"> 发送警告 </el-button>
    <el-button @click="putS"> 发送成功 </el-button>
</template>
