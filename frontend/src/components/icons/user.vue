<script lang="ts" setup>
import { onMounted, ref, watch } from "vue";
let props = defineProps(["modelValue"]);
defineEmits(["update:modelValue", "change"]);
const fillDefault = "var(--main-svg-fill)";
const strokeDefault = "var(--main-svg-stroke)";
const fillSelected = "var(--main-svg-fill-hover)";
const strokeSelected = "var(--main-svg-stroke-hover)";
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
    <svg
        @click="$emit('update:modelValue', true)"
        width="28"
        height="28"
        viewBox="0 0 48 48"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
    >
        <circle
            class="up"
            cx="24"
            cy="12"
            r="8"
            :fill="fillNow"
            :stroke="strokeNow"
            stroke-width="4"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            d="M42 44C42 34.0589 33.9411 26 24 26C14.0589 26 6 34.0589 6 44"
            :stroke="strokeNow"
            stroke-width="4"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
    </svg>
</template>
<style scoped>
svg:hover {
    cursor: pointer;
}
svg:hover circle {
    fill: var(--main-svg-fill-hover);
    stroke: var(--main-svg-stroke-hover);
}
svg:hover path {
    /* fill: var(--main-svg-fill-hover); */
    stroke: var(--main-svg-stroke-hover);
}
</style>
