<script setup lang="ts">
import { provide, ref } from "vue";

import Tabs from "primevue/tabs";
import TabList from "primevue/tablist";
import Tab from "primevue/tab";
import TabPanels from "primevue/tabpanels";
import TabPanel from "primevue/tabpanel";

import MapsTabPanel from "./MapsTab/MapsTabPanel.vue";
import DriverSearchTab from "./DriverSearchTab/DriverSearchTab.vue";
import OrderTabPanel from "./OrderTab/OrderPanel.vue";
import ProfileTabPanel from "./Profile.vue";
import DriverCardPopup from "./DriverCardPopup.vue";
import { ACTIVE_TAB_KEY } from "src/injectionKeys";

const activeTab = ref("mapsTab");
provide(ACTIVE_TAB_KEY, activeTab);
</script>

<template>
  <Tabs v-model:value="activeTab" class="h-full w-full flex flex-col">
    <TabList class="flex w-full">
      <Tab value="mapsTab" class="flex-1 text-center">Карты</Tab>
      <Tab value="driverSearchTab" class="flex-1 text-center"
        >Поиск эвакуатора</Tab
      >
      <Tab value="orderTab" class="flex-1 text-center">Заказ</Tab>
      <Tab value="profileTab" class="flex-1 text-center">Профиль</Tab>
    </TabList>
    <TabPanels class="flex-1 overflow-hidden flex flex-col">
      <TabPanel value="mapsTab" class="flex-1 p-0 m-0">
        <MapsTabPanel />
      </TabPanel>
      <TabPanel value="driverSearchTab" class="flex-1 p-0 m-0">
        <DriverSearchTab />
      </TabPanel>
      <TabPanel value="orderTab" class="flex-1 p-0 m-0">
        <OrderTabPanel />
      </TabPanel>
      <TabPanel value="profileTab" class="flex-1 p-0 m-0">
        <ProfileTabPanel />
      </TabPanel>
    </TabPanels>
  </Tabs>
  <DriverCardPopup />
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
