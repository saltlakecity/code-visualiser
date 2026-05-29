import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  base: './',
  build: {
    modulePreload: false,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('reactflow') || id.includes('@reactflow')) {
            return 'reactflow';
          }
        },
      },
    },
  },
});
