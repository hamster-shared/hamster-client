<template>
  <PageWrapper>
    <Header :showBack="true" />
    <div class="text-center text-[30px] font-bold mb-[40px]">
      {{ t('applications.new.addTitle') }}
    </div>
    <div class="bg-white rounded-[20px] mx-[40px] mb-[100px] py-[40px] px-[90px]">
      <Form
        :model="formData"
        :rules="formRules"
        ref="formRef"
        :labelCol="{ style: { width: '200px' } }"
      >
        <FormItem :label="t('applications.new.name')" name="name">
          <Input
            :allowClear="true"
            class="input-width"
            v-model:value="formData.name"
            :placeholder="t('applications.new.namePlaceholder')"
          />
        </FormItem>
        <FormItem :label="t('applications.new.leaseTerm')" name="leaseTerm">
          <Input
            :allowClear="true"
            class="input-width"
            v-model:value="formData.leaseTerm"
            :placeholder="t('applications.new.leaseTermPlaceholder')"
          />
        </FormItem>
        <FormItem :label="t('applications.new.selectNodeType')" name="selectNodeType">
          <Select
            :allowClear="true"
            class="input-width"
            v-model:value="formData.selectNodeType"
            :placeholder="t('applications.new.selectNodeTypePlaceholder')"
            :options="pluginOptions"
          />
        </FormItem>
        <FormItem :label="t('applications.new.thegraphIndexer')" name="thegraphIndexer">
          <Input
            s
            :allowClear="true"
            class="input-width"
            v-model:value="formData.thegraphIndexer"
            :placeholder="t('applications.new.thegraphIndexerPlaceholder')"
          />
        </FormItem>
        <FormItem :label="t('applications.new.stakingAmount')" name="stakingAmount">
          <Input
            :allowClear="true"
            class="input-width"
            v-model:value="formData.stakingAmount"
            :placeholder="t('applications.new.stakingAmountPlaceholder')"
          />
        </FormItem>
        <FormItem class="text-center">
          <Button size="large" class="w-32 mt-6 ml-4" type="primary" @click="handleSubmit">
            {{ t('common.createText') }}
          </Button>
        </FormItem>
      </Form>
    </div>
  </PageWrapper>
</template>

<script lang="ts" setup>
  import { ref, reactive, computed, toRaw } from 'vue';
  import { useRouter } from 'vue-router';
  import { PageWrapper } from '/@/components/Page';
  import { useI18n } from '/@/hooks/web/useI18n';
  import { useMessage } from '/@/hooks/web/useMessage';
  import { createRule } from '/@/utils/formUtil';
  import Header from '../index/components/Header.vue';
  import { pluginConfigs } from '/@/utils/thegraphUtil';
  import { AddApplication } from '/@wails/go/app/Application';
  import { Form, FormItem, Input, Select, Button } from 'ant-design-vue';

  const { t } = useI18n();
  const router = useRouter();
  const { createErrorModal } = useMessage();

  const pluginOptions = pluginConfigs.map((item) => ({
    ...item,
    label: item.plugin,
    value: item.plugin,
  }));

  const formRef = ref();
  const formData = reactive({
    name: '',
  });

  const formRules = computed(() => ({
    name: [createRule(t('applications.new.namePlaceholder'))],
  }));

  const handleSubmit = async () => {
    await formRef.value?.validate();

    try {
      const params = toRaw(formData);
      params['leaseTerm'] = parseInt(params['leaseTerm']);
      params['stakingAmount'] = parseInt(params['leaseTerm']);
      const { id } = await AddApplication(params);
      router.push(`/applications/${id}`);
    } catch (e: any) {
      createErrorModal({
        title: t('common.errorTip'),
        content: t('applications.new.createFailed'),
      });
    }
  };
</script>

<style lang="less" scoped>
  :deep(.ant-form-item-label > label) {
    @apply !h-[42px];
  }

  :deep(.ant-input-affix-wrapper) {
    @apply !rounded-[8px] !h-[42px];
  }

  :deep(.ant-select-selector) {
    @apply !rounded-[8px] !h-[42px];

    .ant-select-selection-placeholder {
      line-height: 42px;
    }
  }

  :deep(.ant-select-selection-item) {
    @apply !leading-[40px];
  }
</style>
