<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { GetConfig, PutConfig } from "../../wailsjs/go/app/App";
import * as sync from "../../wailsjs/go/app/SyncMan";
import * as user from "../../wailsjs/go/app/UserMan";
import { models } from "../../wailsjs/go/models";

const isUseProxy = ref(false);
const users = ref(new Array<models.User>());
const selectedUid = ref("");
let SelectedUser: models.User;
async function init() {
    let config = await GetConfig();
    users.value = await user.Get();
    if (config.selectedUid !== 0) {
        // 如果用户列表是空的，或者找不到这个用户，就清空 config.selectedUid
        // 否则就选中这个用户
        if (users.value.length !== 0) {
            let u: undefined | models.User = undefined;
            for (let i = 0; i < users.value.length; i++) {
                const user = users.value[i];
                if (user.id === config.selectedUid) {
                    u = user;
                    break;
                }
            }
            if (u !== undefined) {
                SelectedUser = u;
                selectedUid.value = maskUid(u.id.toString());
            } else {
                config.selectedUid === 0;
                PutConfig(config);
            }
            return;
        }
        config.selectedUid === 0;
        PutConfig(config);
    }
}

async function changeUser(user: models.User) {
    SelectedUser = user;
    let config = await GetConfig();
    config.selectedUid = SelectedUser.id;
    PutConfig(config);
}

function maskUid(uid: string): string {
    let reg = /(\d{3})\d{3}(\d{3})/;
    return uid.replace(reg, "$1****$2");
}

async function startSync() {
    if (users.value.length !== 0) {
        if (SelectedUser === undefined) {
            return;
        }
        // 如果有已经之前同步的链接，则先使用之前的链接同步
        if (SelectedUser.raw_url !== "") {
            let result = await sync.Sync(SelectedUser.raw_url);
            if (result !== 0) {
                user.Sync(result, SelectedUser.raw_url);
                return;
            }
            // result===0 说明链接不可用了
            user.Sync(SelectedUser.id, "");
        }
    }

    // 获取新的链接
    let url = await sync.GetRawURL(isUseProxy.value);
    if (url === "") {
        return;
    }
    let result = await sync.Sync(url);
    if (result === 0) {
        return;
    }
    let ok = await user.Sync(result, url);
    if (!ok) {
        return;
    }
    // 更新选择的 uid
    users.value = await user.Get();
    for (let i = 0; i < users.value.length; i++) {
        const user = users.value[i];
        if (user.id === result) {
            changeUser(user);
            selectedUid.value = maskUid(result.toString());
            return;
        }
    }
}
onMounted(() => {
    init();
});
</script>
<template>
    <h1 style="height: 100px">Hi !</h1>
    <div class="view">
        <div class="left">
            <p>一些信息</p>
            <p>一些信息</p>
            <p>一些信息</p>
            <p>一些信息</p>
            <p>一些信息</p>
        </div>
        <div class="right">
            <div style="height: 50px"></div>
            <el-select
                placeholder="选择你的 uid"
                value-key="id"
                v-model="selectedUid"
                @change="changeUser"
            >
                <el-option
                    v-for="user in users"
                    :key="user.id"
                    :label="maskUid(user.id.toString())"
                    :value="user"
                >
                </el-option>
            </el-select>
            <div style="height: 50px"></div>
            <div class="sync">
                <el-button circle class="sync-button" type="primary" @click="startSync"
                    >尝试同步</el-button
                >
                <el-switch
                    v-model="isUseProxy"
                    class="sync-type"
                    active-text="从网络代理"
                    inactive-text="从游戏缓存"
                    style="
                        --el-switch-on-color: var(--el-color-success);
                        --el-switch-off-color: var(--el-color-success);
                    "
                />
            </div>
        </div>
    </div>
</template>
<style scoped>
.sync-button {
    height: calc(24vh);
    width: calc(24vh);
    position: absolute;
    top: 18%;
}

.sync-type {
    position: absolute;
    bottom: 10px;
    color: var(--el-text-color-regular);
}
.view {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
}
.left {
    padding-left: 16px;
    text-align: left;
    width: 46%;
    border-radius: 8px;
    height: 100%;
    background-color: var(--el-fill-color-lighter);
}
.right {
    flex: 1;
    margin-left: 5%;
    display: flex;
    flex-flow: column;
    height: 100%;
}
.sync {
    position: relative;
    height: 100%;
    padding-left: 20px;
    padding-right: 20px;
    text-align: left;
    border-radius: 8px;
    display: flex;
    flex-flow: column;
    align-items: center;
    background-color: var(--el-fill-color-lighter);
}
</style>
