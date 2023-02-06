<script lang="ts" setup>
import { ElMessage, ElNotification } from "element-plus";
import { onMounted, ref, watch } from "vue";
import { GetConfig, PutConfig } from "../../wailsjs/go/app/App";
import * as sync from "../../wailsjs/go/app/SyncMan";
import * as user from "../../wailsjs/go/app/UserMan";
import { models } from "../../wailsjs/go/models";
import download from "../components/icons/download.vue";
import downloading from "../components/icons/downloading.vue";

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
    users.value = await user.Get();
    isLoading.value = true;
    // 数据库的 users 表里有记录但是没有被选中
    if (users.value.length !== 0 && selectedUserIndex.value === -1) {
        isLoading.value = false;
        return;
    }
    let index = selectedUserIndex.value;
    let u = users.value[index];
    // 如果有已经之前同步的链接，则先使用之前的链接同步
    if (index !== -1 && u.raw_url !== "") {
        let result = await sync.Sync(u.raw_url);
        if (result === 0) {
            user.Sync(u.id, "");
            isLoading.value = false;
            return;
        }
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
        // result===0 说明链接不可用了
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
    <div class="box">
        <h2 class="title">在此，同步你的祈愿</h2>
        <div class="view">
            <div class="left">
                <p>一些信息</p>
                <p>一些信息</p>
                <p>一些信息</p>
                <p>一些信息</p>
                <p>一些信息</p>
                <span class="last-sync-time" v-if="selectedUserIndex !== -1">{{ syncTime }}</span>
            </div>
            <div class="right">
                <el-select
                    placeholder="选择你的 uid"
                    value-key="id"
                    v-model="selectedUid"
                    @change="(v: any)=>{selectUser(v.id as number) }"
                >
                    <el-option
                        v-for="user in users"
                        :key="user.id"
                        :value="user"
                        :label="maskUid(user.id)"
                    >
                    </el-option>
                </el-select>
                <div class="sync">
                    <el-button
                        circle
                        :loading="isLoading"
                        style="font-size: 3rem"
                        class="sync-button"
                        type="primary"
                        @click="startSync(isUseProxy)"
                    >
                        <template #loading>
                            <downloading />
                        </template>
                        <template #icon> <download /> </template>
                    </el-button>
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
    </div>
</template>
<style scoped>
.box {
    height: 100%;
    display: flex;
    flex-flow: column;
    position: relative;
}
.title {
    margin-top: 6%;
    margin-bottom: 8%;
}
.last-sync-time {
    color: var(--el-text-color-regular);
    font-size: 10px;
    position: absolute;
    bottom: 8px;
}
.sync-button {
    width: calc(23vh);
    height: calc(23vh);
    margin-bottom: 10%;
}

.sync-type {
    position: absolute;
    bottom: 10px;
    color: var(--el-text-color-regular);
}
.view {
    width: 100%;
    display: flex;
    flex-grow: 1;
    align-items: center;
    margin-bottom: 10px;
}
.left {
    margin-left: 16px;
    margin-right: 6%;
    padding-left: 12px;
    position: relative;
    text-align: left;
    width: 46%;
    border-radius: 8px;
    height: 100%;
    background-color: var(--fill);
}
.right {
    margin-right: 16px;
    flex: 1;
    display: flex;
    flex-flow: column;
    height: 100%;
}
.sync {
    margin-top: 10%;
    position: relative;
    height: 100%;
    padding-left: 20px;
    padding-right: 20px;
    text-align: left;
    border-radius: 8px;
    display: flex;
    flex-flow: column;
    align-items: center;
    justify-content: center;
    background-color: var(--fill);
}
</style>
