<script setup>
import {onBeforeMount, ref} from "vue";
import AppCard from "../../../layout/components/AppCard.vue";
import {appStore} from "../../../store/app.js";

const state = ref({})

onBeforeMount(() => {

})

</script>

<template>
  <div>
    <AppCard :show-border="false">
      <n-h3><b>关于</b></n-h3>
      <div>
        <b>软件版本：</b>{{ appStore().siteConfig.system.version }}
      </div>
    </AppCard>
  </div>
</template>

<style scoped lang="scss">

</style>