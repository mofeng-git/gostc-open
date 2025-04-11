<script setup>

const props = defineProps({
  border: {
    required: false,
    default: false,
    type: Boolean,
  },
  description: {
    required: false,
    default: '无内容',
    type: String,
  }
})

</script>

<template>
  <n-el tag="div" :style="{
      borderColor:'var(--border-color)',
      borderWidth:props.border?'1px':'0',
      borderStyle:'solid',
      padding:'48px',
      borderRadius:'var(--border-radius)',
    }">
    <n-empty size="huge" :description="props.description"></n-empty>
  </n-el>
</template>

<style scoped lang="scss">

</style>