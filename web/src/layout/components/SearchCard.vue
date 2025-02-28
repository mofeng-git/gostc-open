<script setup>
const props = defineProps({
  showBorder: {
    type: Boolean,
    default: true,
  },
  space: {
    type: Boolean,
    default: false,
  },
  vertical:{
    type:Boolean,
    default:false
  }
})

</script>


<template>
  <n-el class="content-container" tag="div" :style="{
      borderColor:'var(--border-color)',
      borderRadius:'var(--border-radius)',
      borderWidth:props.showBorder?'1px':'0',
      borderStyle:'solid',
    }">
    <n-space v-if="props.space" align="center" :vertical="props.vertical">
      <slot></slot>
    </n-space>
    <slot v-else></slot>
  </n-el>
</template>

<style lang="scss" scoped>
.content-container {
  margin: 20px;
  width: calc(100% - 40px);
  padding: 0 !important;
}
</style>