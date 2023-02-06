<script lang="ts" setup>
import { gachaTypeToName } from "../util";
</script>
<template>
    <div style="height: 100%; overflow: hidden">
        <div>
            <h2>- {{ gachaTypeToName($route.query.gachaType as string) }} -</h2>
        </div>
        <h1 @click="$router.back()">Back</h1>
    </div>
</template>
<style scoped></style>
