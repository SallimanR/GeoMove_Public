<script setup lang="ts">
import { ref, provide } from "vue";

import Tabs from "primevue/tabs";
import TabList from "primevue/tablist";
import Tab from "primevue/tab";
import TabPanels from "primevue/tabpanels";
import TabPanel from "primevue/tabpanel";

import MapsTab from "./Tabs/MapsTab.vue";
import Profile from "./Tabs/Profile/Profile.vue";
import OrderSearch from "./Tabs/OrderSearch.vue";
import StreamGps from "./Tabs/StreamGps.vue";
import { ACTIVE_TAB_KEY } from "../injectionKeys";

const activeTab = ref("mainTab");
provide(ACTIVE_TAB_KEY, activeTab);
</script>

<template>
  <Tabs v-model:value="activeTab" class="flex h-full w-full flex-col">
    <TabList class="flex w-full">
      <Tab value="mainTab" class="flex-1 text-center">Главная</Tab>
      <Tab value="orderSearchTab" class="flex-1 text-center">Ищу заказ</Tab>
      <Tab value="statusTab" class="flex-1 text-center">Статус</Tab>
      <Tab value="profileTab" class="flex-1 text-center">Профиль</Tab>
    </TabList>
    <TabPanels class="flex flex-1 flex-col overflow-hidden">
      <TabPanel value="mainTab" class="m-0 flex-1 p-0">
        <MapsTab />
      </TabPanel>
      <TabPanel value="orderSearchTab" class="m-0 flex-1 p-0">
        <OrderSearch />
      </TabPanel>
      <TabPanel value="statusTab"> <StreamGps /> </TabPanel>

      <TabPanel value="profileTab" class="m-0 flex-1 p-0">
        <Profile />
      </TabPanel>
    </TabPanels>
  </Tabs>
</template>

<style scoped>
:deep(.p-tabs) {
  padding: 0 !important;
  height: 100%;
}
:deep(.p-tabpanels) {
  padding: 0 !important;
  flex: 1 !important;
  overflow: hidden !important;
  display: flex !important;
  flex-direction: column !important;
}
:deep(.p-tabpanel) {
  padding: 0 !important;
  flex: 1 !important;
  height: 100% !important;
}
:deep(.p-tablist) {
  padding: 0 !important;
}
</style>
