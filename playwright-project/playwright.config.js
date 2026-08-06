// @ts-check
import { defineConfig, devices } from '@playwright/test';

/**
 * Aponta para o Hat Store rodando localmente via `podman compose up --build`.
 * Pode ser sobrescrito com a variável de ambiente BASE_URL, por exemplo para
 * apontar para o ambiente publicado (https://hatstore-prd.fly.dev).
 */
const baseURL = process.env.BASE_URL || 'http://localhost:8080';

export default defineConfig({
  testDir: './tests',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 2 : undefined,
  reporter: 'html',

  use: {
    baseURL,
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },

  projects: [
    { name: 'chromium', use: { ...devices['Desktop Chrome'] } },
    { name: 'firefox', use: { ...devices['Desktop Firefox'] } },
    { name: 'webkit', use: { ...devices['Desktop Safari'] } },
  ],
});
