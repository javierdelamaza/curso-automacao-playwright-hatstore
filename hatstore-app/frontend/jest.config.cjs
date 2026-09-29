module.exports = {
  testEnvironment: 'jsdom',
  roots: ['<rootDir>/js'],
  moduleFileExtensions: ['js'],
  testMatch: ['**/*.test.js'],
  transform: {
    '^.+\\.js$': ['babel-jest', { presets: ['@babel/preset-env'] }],
  },
};
