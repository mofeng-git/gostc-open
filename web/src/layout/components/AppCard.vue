<script setup>

const props = defineProps({
  loading: {
    default: false,
    type: Boolean
  },
  showBorder: {
    type: Boolean,
    default: true,
  },
  empty: {
    type: Boolean,
    default: false
  }
})

</script>


<template>
  <n-spin class="content-container" :show="props.loading" :delay="200">
    <n-el tag="div" :style="{
      borderColor:'var(--border-color)',
      borderRadius:'var(--border-radius)',
      borderWidth:props.showBorder?'1px':'0',
      borderStyle:'solid',
    }">
      <n-empty v-if="props.empty" description="没有数据" style="margin: 30px 0"></n-empty>
      <slot v-else></slot>
    </n-el>
  </n-spin>
</template>

<style lang="scss" scoped>
.content-container {
  margin: 20px;
  width: calc(100% - 40px);
  padding: 0 !important;

  :deep(.n-card) {
    .n-card__content {
      padding: 0 !important;
    }
  }
}

</style>