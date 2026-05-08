import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  images: {
    remotePatterns: [
      {
        hostname: "www.google.com",
        protocol: "https",
      },
    ],
  },
  outputFileTracingRoot: process.cwd(),
};

export default nextConfig;
