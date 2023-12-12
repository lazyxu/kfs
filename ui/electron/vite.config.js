import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

export default defineConfig({
  build: {
    commonjsOptions: {
      transformMixedEsModules: true,
    },
  },
  plugins: [react()], optimizeDeps: {
    include: ['@emotion/styled', '@mui/icons-material'],
  },
});
