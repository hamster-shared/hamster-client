<template>
  <PageWrapper>
    <Header :showBack="true" />
    <div class="text-center text-[30px] font-bold mb-[40px]">
      {{ t('applications.new.addTitle', { title: titleName }) }}
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
        <!--    thegraph form start     -->
        <FormItem
          :label="t('applications.new.selectNodeType')"
          name="selectNodeType"
          v-if="blockchain === 'thegraph'"
        >
          <Select
            :allowClear="true"
            class="input-width"
            v-model:value="formData.selectNodeType"
            :placeholder="t('applications.new.selectNodeTypePlaceholder')"
            :options="pluginConfigs"
          />
        </FormItem>
        <FormItem
          :label="t('applications.new.thegraphIndexer')"
          name="thegraphIndexer"
          v-if="blockchain === 'thegraph'"
        >
          <Input
            :allowClear="true"
            class="input-width"
            v-model:value="formData.thegraphIndexer"
            :placeholder="t('applications.new.thegraphIndexerPlaceholder')"
          />
        </FormItem>
        <FormItem
          :label="t('applications.new.stakingAmount')"
          name="stakingAmount"
          v-if="blockchain === 'thegraph'"
        >
          <Input
            :allowClear="true"
            class="input-width"
            v-model:value="formData.stakingAmount"
            :placeholder="t('applications.new.stakingAmountPlaceholder')"
          />
        </FormItem>
        <!--    thegraph form end     -->

        <!--    ethereum form start     -->
        <FormItem
          :label="t('applications.new.selectNodeType')"
          name="selectNodeType"
          v-if="blockchain === 'ethereum' || blockchain === 'starkware'"
        >
          <Select
            :allowClear="true"
            class="input-width"
            v-model:value="formData.selectNodeType"
            :placeholder="t('applications.new.selectNodeTypePlaceholder')"
            :options="ethereum_network_option"
          />
        </FormItem>

        <FormItem
          :label="t('applications.new.selectNodeType')"
          name="selectNodeType"
          v-if="blockchain === 'near'"
        >
          <Select
            :allowClear="true"
            class="input-width"
            v-model:value="formData.selectNodeType"
            :placeholder="t('applications.new.selectNodeTypePlaceholder')"
            :options="near_network_option"
          />
        </FormItem>

        <FormItem class="text-center">
          <Button
            size="large"
            class="w-32 mt-6 ml-4"
            type="primary"
            shape="round"
            @click="handleSubmit"
            :loading="createLoading"
          >
            {{ t('common.createText') }}
          </Button>
        </FormItem>
      </Form>
    </div>
  </PageWrapper>
</template>

<script lang="ts" setup>
  import { ref, reactive, computed, toRaw, onMounted } from 'vue';
  import { useRouter } from 'vue-router';
  import { PageWrapper } from '/@/components/Page';
  import { useI18n } from '/@/hooks/web/useI18n';
  import { useMessage } from '/@/hooks/web/useMessage';
  import { createRule } from '/@/utils/formUtil';
  import Header from '../index/components/Header.vue';
  import { pluginConfigs } from '/@/utils/thegraphUtil';
  import { ethereum_network_option } from '/@/utils/ethereum';
  import { near_network_option } from '/@/utils/near';
  import { AddApplication } from '/@wails/go/app/Application';
  import { Form, FormItem, Input, Select, Button } from 'ant-design-vue';

  const { t } = useI18n();
  const router = useRouter();
  const { createErrorModal } = useMessage();

  const { blockchain } = router.currentRoute.value.query;
  const titleName = ref('');

  // const formCoponents = {
  //   thegraph: ThegraphForm,
  //   ethereum: EthereumForm,
  //   bsc: BscForm,
  //   polygon: PolygonForm,
  //   avalanche: AvalancheForm,
  //   optimism: OptimismForm,
  //   zksync: ZksyncForm,
  //   starkware: StarkwareForm,
  //   near: NearForm,
  //   aptos: AptosForm,
  //   sui: SuiForm,
  // };

  // const BlockChainForm = formCoponents[blockchain] || ThegraphForm;

  const formRef = ref();
  const createLoading = ref(false);
  const formData = reactive({
    name: '',
  });

  const formRules = computed(() => ({
    name: [createRule(t('applications.new.namePlaceholder'))],
  }));

  const handleSubmit = async () => {
    await formRef.value?.validate();
    createLoading.value = true;

    try {
      const params = toRaw(formData);
      params['serviceType'] = blockchain;
      params['leaseTerm'] = parseInt(params['leaseTerm']);
      params['stakingAmount'] = parseInt(params['stakingAmount']);
      const { id } = await AddApplication(params);
      router.push(`/applications/${id}`);
    } catch (error: any) {
      const errorMessage = typeof error === 'string' ? error : error?.message;
      createErrorModal({
        title: t('common.errorTip'),
        content: errorMessage || t('applications.new.createFailed'),
      });
    } finally {
      createLoading.value = false;
    }
  };
  onMounted(() => {
    console.log('blockchain', blockchain);
    formData.serviceType = blockchain;
    console.log('formData', formData);
    titleName.value = blockchain?.slice(0, 1).toUpperCase() + blockchain?.slice(1);
  });
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
