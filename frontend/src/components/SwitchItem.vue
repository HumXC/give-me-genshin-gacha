<script lang="ts" setup>
defineProps(["modelValue"]);
defineEmits(["update:modelValue", "change"]);
</script>

<template>
    <div class="switch-item">
        <span class="name">
            <slot> </slot>
        </span>
        <el-switch
            class="switch"
            :model-value="modelValue"
            @change="(val:boolean)=>{$emit('update:modelValue',val);$emit('change',val)}"
        />
    </div>
</template>
<style scoped>
.switch-item {
    width: 100%;
    height: 50px;
    color: var(--el-text-primary);
    background-color: var(--el-fill-color-lighter);
    border-width: 1px;
    border-style: solid;
    border-radius: 5px;
    border-color: var(--el-border-color);
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.name {
    margin-left: 10px;
}
.switch {
    margin-right: 10px;
}
</style>
