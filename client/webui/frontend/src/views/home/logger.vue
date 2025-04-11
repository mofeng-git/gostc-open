<script setup>
import {computed, onBeforeMount, ref} from "vue";
import {apiLoggerList} from "../../api/logger.js";
import {appStore} from "../../store/app.js";

const state = ref({
  data: [],
  logs:'',
})


const listFunc = async () => {
  try {
    let res = await apiLoggerList()
    state.value.data = res.data || []
    let data = []
    state.value.data.forEach(item=>{
      data.push(`[${item.timestamp}] ${item.message}`)
    })
    data.reverse()
    state.value.logs = data.join('\n')
  } finally {

  }
}

onBeforeMount(() => {
  listFunc()
})
</script>

<template>
  <div>
    <n-space>
      <n-button size="small" type="info" @click="listFunc">刷新</n-button>
    </n-space>

    <p/>

    <n-log
        :rows="20"
        :log="state.logs"/>
  </div>
</template>
