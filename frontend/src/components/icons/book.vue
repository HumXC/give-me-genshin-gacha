<script lang="ts" setup>
import { onMounted, ref, watch } from "vue";
let props = defineProps(["modelValue"]);
defineEmits(["update:modelValue", "change"]);
const fillDefault = "var(--svg-fill)";
const strokeDefault = "var(--svg-stroke)";
const fillSelected = "var(--svg-fill-hover)";
const strokeSelected = "var(--svg-stroke-hover)";
const fillNow = ref(fillDefault);
const strokeNow = ref(strokeDefault);
function light(isLighting: boolean) {
    if (isLighting) {
        fillNow.value = fillSelected;
        strokeNow.value = strokeSelected;
    } else {
        fillNow.value = fillDefault;
        strokeNow.value = strokeDefault;
    }
}
watch(props, (v) => {
    light(v.modelValue);
});
onMounted(() => {
    light(props.modelValue);
});
</script>
<template>
    <div @click="$emit('update:modelValue', true)">
        <svg
            width="28"
            height="28"
            viewBox="0 0 48 48"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
        >
            <path
                class="in"
                d="M7 37C7 29.2967 7 11 7 11C7 7.68629 9.68629 5 13 5H35V31C35 31 18.2326 31 13 31C9.7 31 7 33.6842 7 37Z"
                :fill="fillNow"
                :stroke="strokeNow"
                stroke-width="4"
                stroke-linejoin="round"
            />
            <path
                d="M35 31C35 31 14.1537 31 13 31C9.68629 31 7 33.6863 7 37C7 40.3137 9.68629 43 13 43C15.2091 43 25.8758 43 41 43V7"
                stroke-width="4"
                :stroke="strokeNow"
                stroke-linecap="round"
                stroke-linejoin="round"
            />
            <path
                d="M14 37H34"
                :stroke="strokeNow"
                stroke-width="4"
                stroke-linecap="round"
                stroke-linejoin="round"
            />
        </svg>
    </div>
</template>
<style scoped>
div {
    display: flex;
    align-items: center;
    justify-content: center;
}
svg:hover path {
    stroke: var(--svg-stroke-hover);
}
svg:hover > path.in {
    fill: var(--svg-fill-hover);
}
</style>
