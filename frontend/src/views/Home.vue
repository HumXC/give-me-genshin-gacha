<script lang="ts" setup>
import { ElMessage, ElNotification } from "element-plus";
import { onMounted, ref, watch } from "vue";
import { GetConfig, PutConfig } from "../../wailsjs/go/app/App";
import * as sync from "../../wailsjs/go/app/SyncMan";
import * as user from "../../wailsjs/go/app/UserMan";
import { models } from "../../wailsjs/go/models";

const syncTime = ref("");
const isLoading = ref(false);
const isUseProxy = ref(false);
const users = ref(new Array<models.User>());
const selectedUserIndex = ref(-1);
const selectedUid = ref("");
async function init() {
    let config = await GetConfig();
    isUseProxy.value = config.isUseProxy;
    watch(isUseProxy, async (v) => {
        let c = await GetConfig();
        c.isUseProxy = v;
        PutConfig(c);
    });
    users.value = await user.Get();

    // 选中配置里的 selectedUid
    if (config.selectedUid !== 0) {
        // 如果用户列表是空的，或者找不到这个用户，就清空 config.selectedUid
        // 否则就选中这个用户
        if (users.value.length === 0) {
            config.selectedUid === 0;
            PutConfig(config);
            // 自动同步选项
            if (config.isAutoSync) {
                // 使用代理同步会破坏体验
                startSync(false);
            }
            return;
        }
        selectUser(config.selectedUid);
        // 自动同步选项
        if (config.isAutoSync) {
            // 使用代理同步会破坏体验
            startSync(false);
        }
        return;
    }
}
function CreatProxyNotify(onClose: () => void): () => void {
    return ElNotification({
        title: "已经开启代理服务器",
        message: "重新在游戏里打开祈愿记录，关闭此通知会关闭代理服务器并取消同步",
        onClose: onClose,
        duration: 0,
        offset: 0,
    }).close;
}

function maskUid(uid: number): string {
    let reg = /(\d{3})\d{3}(\d{3})/;
    return uid.toString().replace(reg, "$1****$2");
}

async function startSync(isUseProxy: boolean) {
    isLoading.value = true;
    if (users.value.length !== 0) {
        // 没有已经选中的用户
        if (selectedUserIndex.value === -1) {
            isLoading.value = false;
            return;
        }
        let index = selectedUserIndex.value;
        let u = users.value[index];
        // 如果有已经之前同步的链接，则先使用之前的链接同步
        if (u.raw_url !== "") {
            let result = await sync.Sync(u.raw_url);
            if (result !== 0) {
                let ok = await user.Sync(result, u.raw_url);
                if (!ok) {
                    isLoading.value = false;
                    return;
                }
                isLoading.value = false;
                ElMessage({
                    type: "success",
                    message: "同步完成！",
                });
                selectUser(result);
                return;
            }
            // result===0 说明链接不可用了
            user.Sync(u.id, "");
        }
    }

    // 获取新的链接
    let closeNotifytion: (() => void) | null = null;
    if (isUseProxy) {
        closeNotifytion = CreatProxyNotify(() => {
            sync.StopProxyServer();
        });
    }
    // 获取链接
    let url = await sync.GetRawURL(isUseProxy);
    if (closeNotifytion !== null) closeNotifytion();
    if (url === "") {
        isLoading.value = false;
        return;
    }
    // 拉取祈愿数据
    let result = await sync.Sync(url);
    if (result === 0) {
        isLoading.value = false;
        return;
    }
    // 同步用户信息
    let ok = await user.Sync(result, url);
    if (!ok) {
        isLoading.value = false;
        return;
    }
    // 更新选择的 uid
    selectUser(result);
    ElMessage({
        type: "success",
        message: "同步完成！",
    });
    isLoading.value = false;
}

// 选择用户并更新左侧信息，相当于在选择框里选择了一个 user
async function selectUser(uid: number) {
    users.value = await user.Get();
    let index = -1;
    for (let i = 0; i < users.value.length; i++) {
        const user = users.value[i];
        if (user.id === uid) {
            index = i;
            break;
        }
    }
    if (index === -1) {
        return;
    }
    let u = users.value[index];
    selectedUserIndex.value = index;
    syncTime.value = formatTime(u.sync_time);
    selectedUid.value = maskUid(u.id);
    let config = await GetConfig();
    config.selectedUid = u.id;
    PutConfig(config);
}
function formatTime(time: any): string {
    let t = time as string;
    return t.substring(0, t.indexOf(".")).replace("T", " ");
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
            <div class="last-sync-time" v-if="selectedUserIndex !== -1">
                <snap>上一次同步时间: </snap>
                <snap>{{ syncTime }}</snap>
            </div>
        </div>
        <div class="right">
            <div style="height: 50px"></div>
            <el-select
                placeholder="选择你的 uid"
                value-key="id"
                v-model="selectedUid"
                @change="selectUser"
            >
                <el-option v-for="user in users" :key="user.id" :value="user"> </el-option>
            </el-select>
            <div style="height: 50px"></div>
            <div class="sync">
                <el-button
                    circle
                    :loading="isLoading"
                    class="sync-button"
                    type="primary"
                    @click="startSync(isUseProxy)"
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
.last-sync-time {
    color: var(--el-text-color-regular);
    font-size: 10px;
    position: absolute;
    bottom: 8px;
}
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
    position: relative;
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
