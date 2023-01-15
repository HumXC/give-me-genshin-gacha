<script lang="ts" setup>
import { Sync } from "../../wailsjs/go/main/App";

let sendErr = async () => {
    let data = await Sync(true);
    console.log(JSON.stringify(data));
};
</script>
<template>
    <div id="right">
        <button @click="sendErr">发送错误</button>
    </div>
</template>
