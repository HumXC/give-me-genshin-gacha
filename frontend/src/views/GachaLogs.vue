<script lang="ts" setup>
import { gachaTypeToName } from "../util";
</script>
<template>
    <div>
        <div class="header">
            <h2>- {{ gachaTypeToName($route.query.gachaType as string) }} -</h2>
        </div>
        <H1 @click="$router.back()">Back</H1>
        <h2>Hello</h2>
    </div>
</template>
<style scoped></style>
