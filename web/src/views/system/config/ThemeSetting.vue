<template>
  <div>
    <n-spin :show="show" description="请稍候...">
      <n-form :label-width="80" :model="formValue" :rules="rules" ref="formRef">
        <n-form-item label="默认主题" path="themeDarkTheme">
          <n-input v-model:value="formValue.themeDarkTheme" placeholder="" />
          <template #feedback> 可选：'dark' 、 'light' </template>
        </n-form-item>

        <n-form-item label="默认系统主题" path="themeAppTheme">
          <n-input v-model:value="formValue.themeAppTheme" placeholder="" />
          <template #feedback> 默认：#2d8cf0 </template>
        </n-form-item>
        <n-form-item label="默认侧边栏风格" path="themeNavTheme">
          <n-input v-model:value="formValue.themeNavTheme" placeholder="" />
          <template #feedback>可选：'light'、 'header-dark'</template>
        </n-form-item>

        <div>
          <n-space>
            <n-button type="primary" @click="formSubmit">保存更新</n-button>
          </n-space>
        </div>
      </n-form>
    </n-spin>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue';
  import { useMessage } from 'naive-ui';
  import { getConfig, updateConfig } from '@/api/sys/config';

  const group = ref('theme');
  const show = ref(false);

  const rules = {
    themeDarkTheme: {
      required: true,
      message: '请输入默认主题',
      trigger: 'blur',
    },
  };

  const formRef: any = ref(null);
  const message = useMessage();

  const formValue = ref({
    themeDarkTheme: 'dark',
    themeAppTheme: '#2d8cf0',
    themeNavTheme: 'dark',
  });

  function formSubmit() {
    formRef.value.validate((errors) => {
      if (!errors) {
        updateConfig({ group: group.value, list: formValue.value }).then((_res) => {
          message.success('更新成功');
          load();
        });
      } else {
        message.error('验证失败，请填写完整信息');
      }
    });
  }

  onMounted(() => {
    load();
  });

  function load() {
    show.value = true;
    new Promise((_resolve, _reject) => {
      getConfig({ group: group.value })
        .then((res) => {
          formValue.value = res.list;
        })
        .finally(() => {
          show.value = false;
        });
    });
  }
</script>
