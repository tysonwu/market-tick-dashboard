/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // webpack: (config, { isServer }) => {
  //     if (!isServer) {
  //           config.resolve.fallback.fs = false
  //           config.resolve.fallback.dns = false
  //           config.resolve.fallback.net = false
  //     }
  //     return config;
  // }
}

module.exports = nextConfig
