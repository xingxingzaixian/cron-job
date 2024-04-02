import type { UserConfigExport, ConfigEnv } from 'vite';
import { loadEnv } from 'vite';
import path from 'node:path';
import vue from '@vitejs/plugin-vue';
import vueJsx from '@vitejs/plugin-vue-jsx';
import { resolve } from 'path';
import UnoCSS from 'unocss/vite';
import Icons from 'unplugin-icons/vite';
import Components from 'unplugin-vue-components/vite';
import IconsResolver from 'unplugin-icons/resolver';
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons';
import { FileSystemIconLoader } from 'unplugin-icons/loaders';
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers';

// https://vitejs.dev/config/
export default ({ mode }: ConfigEnv): UserConfigExport => {
  const { VITE_APP_PORT, VITE_ICON_PREFIX, VITE_ICON_LOCAL_PREFIX, VITE_APP_PUBLISH_PATH } = loadEnv(
    mode,
    process.cwd()
  );
  const localIconPath = path.join(process.cwd(), 'src/assets/svg-icon');
  const collectionName = VITE_ICON_LOCAL_PREFIX.replace(`${VITE_ICON_PREFIX}-`, '');

  const INVALID_CHAR_REGEX = /[\x00-\x1F\x7F<>*#"{}|^[\]`;?:&=+$,]/g;
  const DRIVE_LETTER_REGEX = /^[a-z]:/i;

  return {
    base: VITE_APP_PUBLISH_PATH,
    plugins: [
      vue({
        script: {
          defineModel: true
        }
      }),
      vueJsx(),
      UnoCSS(),
      Components({
        resolvers: [NaiveUiResolver(), IconsResolver({ prefix: VITE_ICON_PREFIX })]
      }),
      Icons({
        compiler: 'vue3',
        autoInstall: true,
        customCollections: {
          [collectionName]: FileSystemIconLoader(localIconPath, (svg) =>
            svg.replace(/^<svg\s/, '<svg width="1em" height="1em" ')
          )
        },
        scale: 1,
        defaultClass: 'inline-block'
      }),
      createSvgIconsPlugin({
        iconDirs: [localIconPath],
        symbolId: 'icon-[dir]-[name]'
      })
    ],
    resolve: {
      alias: [
        {
          find: /@\//,
          replacement: resolve(__dirname, 'src') + '/'
        }
      ],
      extensions: ['.ts', '.js', '.jsx', '.tsx']
    },
    server: {
      https: false,
      host: true,
      port: Number(VITE_APP_PORT)
    },
    build: {
      target: 'es2015',
      chunkSizeWarningLimit: 1500,
      rollupOptions: {
        output: {
          sanitizeFileName(name) {
            const match = DRIVE_LETTER_REGEX.exec(name);
            const driveLetter = match ? match[0] : '';
            // 确保名称中不包含非法字符，主要是不能以_开头
            return driveLetter + name.slice(driveLetter.length).replace(INVALID_CHAR_REGEX, '');
          }
        }
      }
    }
  };
};
