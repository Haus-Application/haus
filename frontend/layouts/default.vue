<template>
  <div class="app-layout">
    <slot />
  </div>
</template>

<style scoped>
.app-layout {
  min-height: 100vh;
  background: var(--color-bg);
  color: var(--color-text);
}
</style>
