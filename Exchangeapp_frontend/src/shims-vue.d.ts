// client/src/shims-vue.d.ts
declare module '*.vue' {
    import { DefineComponent } from 'vue';
    const component: DefineComponent<{}, {}, any>;
    export default component;
}

declare module 'element-plus/dist/locale/zh-cn.mjs' {
  const zhCn: any;
  export default zhCn;
}

declare module 'element-plus/dist/locale/zh-cn' {
  const zhCn: any;
  export default zhCn;
}
  