import { type InjectionKey, type Ref } from "vue";

export const ACTIVE_TAB_KEY: InjectionKey<Ref<string>> = Symbol("activeTab");
