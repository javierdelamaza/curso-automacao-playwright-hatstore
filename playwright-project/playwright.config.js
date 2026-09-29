// @ts-check
import { defineConfig, devices } from '@playwright/test';

// Hat Store local (hatstore-app/, via `podman compose up --build`).
const baseURL = process.env.BASE_URL || 'http://localhost:8080';

export default defineConfig({
  testDir: './tests',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 2 : undefined,
  reporter: [['html', { open: 'never' }], ['list']],

  expect: {
    toHaveScreenshot: { maxDiffPixelRatio: 0.01, animations: 'disabled' },
  },

  use: {
    baseURL,
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
      // Testes visuais dependem de baseline gerada na sua máquina: rodam só no projeto "visual".
      grepInvert: /@visual/,
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
      grepInvert: /@visual/,
    },
    {
      name: 'visual',
      use: { ...devices['Desktop Chrome'] },
      grep: /@visual/,
    },
    {
      name: 'solucoes',
      testDir: './solucoes',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
});
