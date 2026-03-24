<template>
  <div class="formula-panel">
    <el-button class="close-btn" :icon="Close" @click="emit('close')" text circle />
    <div class="panel-section function-keys">
      <el-button @click="emit('insert-event')" size="small">插入事件指标</el-button>
      <el-button @click="emit('clear')" size="small" type="danger" plain>清空</el-button>
    </div>
    <div class="panel-section numbers">
      <div class="button-grid">
        <button v-for="num in numberKeys" :key="num" @click="emit('add-token', { type: 'number', value: num })">
          {{ num }}
        </button>
      </div>
    </div>
    <div class="panel-section operators">
      <div class="button-grid">
        <button v-for="op in operatorKeys" :key="op" @click="emit('add-token', { type: 'operator', value: op })">
          {{ op }}
        </button>
        <button @click="emit('delete-token')">del</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { CustomElement } from '@/types/doris/customFormula';
import { Close } from '@element-plus/icons-vue';

interface Emits {
  (e: 'add-token', token: Pick<CustomElement, 'type' | 'value'>): void;
  (e: 'delete-token'): void;
  (e: 'clear'): void;
  (e: 'insert-event'): void;
  (e: 'close'): void;
}
const emit = defineEmits<Emits>();

const numberKeys = ['1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '.'];
const operatorKeys = ['+', '-', '*', '/', '(', ')'];
</script>

<style scoped>
.formula-panel {
  position: relative;
  padding: 12px;
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,.1);
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 280px;
}
.close-btn {
  position: absolute;
  top: 5px;
  right: 5px;
}
.panel-section .button-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}
.operators .button-grid {
  grid-template-columns: repeat(4, 1fr);
}
.panel-section button {
  height: 32px;
  border-radius: 4px;
  border: 1px solid #dcdfe6;
  background-color: #fff;
  cursor: pointer;
  font-size: 14px;
}
.function-keys {
  display: flex;
  gap: 8px;
}
.function-keys .el-button {
  flex: 1;
}
</style>
